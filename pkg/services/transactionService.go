package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"golang/pkg/models"
	"golang/pkg/repository"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionService struct {
	repo       repository.TransactionRepository
	csvService CSVServiceI
}

type TransactionServiceI interface {
	GetTransactions(filters map[string]interface{}, limit, offset int) ([]models.Transaction, int64, error)
	CreateTransactionInBatches(transactions []*models.Transaction, batchSize int) error
}

func NewTransactionService(repo repository.TransactionRepository, csvService CSVServiceI) *TransactionService {
	return &TransactionService{repo: repo, csvService: csvService}
}

func (s *TransactionService) GetTransactions(filters map[string]interface{}, limit, offset int) ([]models.Transaction, error) {
	return s.repo.GetTransactions(filters, limit, offset)
}

func (s *TransactionService) CreateTransactionInBatches(transactions []*models.Transaction, batchSize int) error {
	return s.repo.CreateTransactionInBatches(transactions, batchSize)
}

func (s *TransactionService) ParseTransaction(headers []string, record []string) (models.Transaction, error) {
	var transaction models.Transaction

	transactionType := reflect.TypeOf(transaction)
	transactionValue := reflect.ValueOf(&transaction).Elem()

	for i, header := range headers {
		var field reflect.StructField
		for j := 0; j < transactionType.NumField(); j++ {
			f := transactionType.Field(j)
			if strings.EqualFold(f.Name, header) {
				field = f
				break
			}
		}

		if field.Name == "" {
			continue
		}

		fieldValue := record[i]

		switch field.Type.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(fieldValue, 10, 64)
			if err != nil {
				return transaction, err
			}
			transactionValue.FieldByName(field.Name).SetUint(val)
		case reflect.Float64:
			val, err := strconv.ParseFloat(fieldValue, 64)
			if err != nil {
				return transaction, err
			}
			transactionValue.FieldByName(field.Name).SetFloat(val)
		case reflect.String:
			transactionValue.FieldByName(field.Name).SetString(fieldValue)
		case reflect.Struct:
			if field.Type == reflect.TypeOf(time.Time{}) {
				val, err := time.Parse("2006-01-02 15:04:05", fieldValue)
				if err != nil {
					return transaction, err
				}
				transactionValue.FieldByName(field.Name).Set(reflect.ValueOf(val))
			}
		}
	}

	return transaction, nil
}

func (s *TransactionService) ProcessCSVFile(reader *csv.Reader) error {
	const chunkSize = 1000

	headers, err := s.csvService.ReadCSVHeaders(reader)
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transactionsCh := make(chan [][]string, 10)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	go s.readCSVChunks(ctx, reader, chunkSize, transactionsCh, errCh)
	s.startWorkers(ctx, headers, transactionsCh, errCh, &wg)

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			cancel()
			return err
		}
	case <-ctx.Done():
	}

	return nil
}

func (s *TransactionService) readCSVChunks(ctx context.Context, reader *csv.Reader, chunkSize int, transactionsCh chan<- [][]string, errCh chan<- error) {
	defer close(transactionsCh)
	for {
		var transactions [][]string
		for i := 0; i < chunkSize; i++ {
			record, err := s.csvService.ReadCSVRecord(reader)
			if err == io.EOF {
				break
			}
			if err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
				}
				return
			}
			transactions = append(transactions, record)
		}

		if len(transactions) == 0 {
			break
		}

		select {
		case transactionsCh <- transactions:
		case <-ctx.Done():
			return
		}
	}
}

func (s *TransactionService) startWorkers(ctx context.Context, headers []string, transactionsCh <-chan [][]string, errCh chan<- error, wg *sync.WaitGroup) {
	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for transactions := range transactionsCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if _, err := s.processTransactions(headers, transactions); err != nil {
					select {
					case errCh <- err:
					case <-ctx.Done():
					}
					return
				}
			}
		}()
	}
}

func (s *TransactionService) processTransactions(headers []string, records [][]string) ([]*models.Transaction, error) {
	parsedTransactions := make([]*models.Transaction, len(records))
	for j, record := range records {
		transaction, err := s.ParseTransaction(headers, record)
		if err != nil {
			return nil, fmt.Errorf("error parsing record: %v", err)
		}
		parsedTransactions[j] = &transaction
	}

	if err := s.repo.CreateTransactionInBatches(parsedTransactions, len(records)); err != nil {
		return nil, fmt.Errorf("error inserting transactions: %v", err)
	}

	return parsedTransactions, nil
}

func (c *TransactionService) ParseFilters(ctx *gin.Context) (map[string]interface{}, error) {
	filters := make(map[string]interface{})

	if transactionID := ctx.Query("transaction_id"); transactionID != "" {
		filters["transaction_id"] = transactionID
	}
	if terminalID := ctx.Query("terminal_id"); terminalID != "" {
		filters["terminal_id"] = terminalID
	}
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}
	if paymentType := ctx.Query("payment_type"); paymentType != "" {
		filters["payment_type"] = paymentType
	}
	if dateFrom := ctx.Query("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo := ctx.Query("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}
	if paymentNarrative := ctx.Query("payment_narrative"); paymentNarrative != "" {
		filters["payment_narrative"] = paymentNarrative
	}

	return filters, nil
}

// csv export
func (s *TransactionService) ExportTransactionsCSV(ctx *gin.Context, filters map[string]interface{}) error {
	reader, writer := io.Pipe()
	csvWriter := csv.NewWriter(writer)

	headers := s.csvService.GetCSVHeaders(models.Transaction{})

	if err := csvWriter.Write(headers); err != nil {
		writer.CloseWithError(fmt.Errorf("error writing CSV headers: %w", err))
		return err
	}

	transactionsCh := make(chan []models.Transaction, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		s.fetchTransactions(filters, transactionsCh, writer)
	}()

	numWorkers := 1
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.writeTransactionsToCSV(ctx, transactionsCh, csvWriter, writer)
		}()
	}

	go func() {
		io.Copy(ctx.Writer, reader)
	}()

	wg.Wait()
	writer.Close()
	return nil
}

func (s *TransactionService) fetchTransactions(filters map[string]interface{}, transactionsCh chan<- []models.Transaction, writer *io.PipeWriter) {
	defer close(transactionsCh)
	limit := 500

	for offset := 0; ; offset += limit {
		transactions, err := s.repo.GetTransactions(filters, limit, offset)
		if err != nil {
			writer.CloseWithError(fmt.Errorf("error fetching transactions: %w", err))
			return
		}

		if len(transactions) == 0 {
			return
		}

		transactionsCh <- transactions
	}
}

func (s *TransactionService) writeTransactionsToCSV(ctx *gin.Context, transactionsCh <-chan []models.Transaction, csvWriter *csv.Writer, writer *io.PipeWriter) {
	for transactionChunk := range transactionsCh {
		select {
		case <-ctx.Done():
			writer.CloseWithError(fmt.Errorf("context canceled"))
			return
		default:
			for _, t := range transactionChunk {
				record := transactionToCSVRecord(t)
				if err := csvWriter.Write(record); err != nil {
					writer.CloseWithError(fmt.Errorf("error writing CSV record: %w", err))
					return
				}
			}
			csvWriter.Flush()
			if err := csvWriter.Error(); err != nil {
				writer.CloseWithError(fmt.Errorf("error flushing CSV writer: %w", err))
				return
			}
		}
	}
}

func transactionToCSVRecord(t models.Transaction) []string {
	return []string{
		strconv.Itoa(int(t.TransactionID)),
		strconv.Itoa(int(t.RequestID)),
		strconv.Itoa(int(t.TerminalID)),
		strconv.Itoa(int(t.PartnerObjectID)),
		strconv.FormatFloat(t.AmountTotal, 'f', 2, 64),
		strconv.FormatFloat(t.AmountOriginal, 'f', 2, 64),
		strconv.FormatFloat(t.CommissionPS, 'f', 2, 64),
		strconv.FormatFloat(t.CommissionClient, 'f', 2, 64),
		strconv.FormatFloat(t.CommissionProvider, 'f', 2, 64),
		t.DateInput.Format(time.RFC3339),
		t.DatePost.Format(time.RFC3339),
		t.Status,
		t.PaymentType,
		t.PaymentNumber,
		strconv.Itoa(int(t.ServiceID)),
		t.Service,
		strconv.Itoa(int(t.PayeeID)),
		t.PayeeName,
		strconv.Itoa(int(t.PayeeBankMfo)),
		t.PayeeBankAccount,
		t.PaymentNarrative,
	}
}

// json export
func (s *TransactionService) ExportTransactionsJSON(ctx *gin.Context, filters map[string]interface{}) error {
	reader, writer := io.Pipe()

	transactionsCh := make(chan []models.Transaction)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.fetchTransactions(filters, transactionsCh, writer)
	}()

	go func() {
		defer wg.Done()
		s.writeTransactionsToJSON(ctx, transactionsCh, writer)
	}()

	ctx.Writer.Header().Set("Content-Type", "application/json")
	_, err := io.Copy(ctx.Writer, reader)
	wg.Wait()

	return err
}

func (s *TransactionService) writeTransactionsToJSON(ctx *gin.Context, transactionsCh <-chan []models.Transaction, writer *io.PipeWriter) {
	defer writer.Close()

	writer.Write([]byte("["))

	first := true
	encoder := json.NewEncoder(writer)

	for transactionChunk := range transactionsCh {
		select {
		case <-ctx.Done():
			writer.CloseWithError(fmt.Errorf("context canceled"))
			return
		default:
			for _, transaction := range transactionChunk {
				if !first {
					writer.Write([]byte(","))
				}
				first = false

				if err := encoder.Encode(transaction); err != nil {
					writer.CloseWithError(fmt.Errorf("error encoding JSON: %w", err))
					return
				}
			}
		}
	}

	writer.Write([]byte("]"))
}
