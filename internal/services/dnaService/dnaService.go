// Package services provides functionality for analyzing DNA sequences.
package dna

import (
	"bufio"
	"context"
	"golang/internal/utils"
	"mime/multipart"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/sync/errgroup"
)

var (
	gcCountWorkers = 5 // concurrent gcCount workers
	gcPrecision    = 1 // Precision for rounding GC content percentage
)

type GCResult struct {
	TotalSequenceLen int
	GCContentLen     int
}

type DNAService struct {
}

type DNAServiceI interface {
	AnalyzeDNASeq() string
}

func NewDNAService() *DNAService {
	return &DNAService{}
}

// gfd
func (s *DNAService) AnalyzeGCContent(file multipart.File) (float64, error) {
	var (
		gcCountJobs = make(chan string, 1)
		gcCountRes  = make(chan GCResult, 1)
	)

	var (
		totalGCCount   = 0
		totalSeqLength = 0
	)

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return s.readDNASequences(ctx, gcCountJobs, file)
	})

	var wg sync.WaitGroup

	for i := 0; i < gcCountWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.gcCount(ctx, gcCountJobs, gcCountRes)
		}()
	}

	go func() {
		wg.Wait()
		close(gcCountRes)
	}()

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case gcRes, ok := <-gcCountRes:
				if !ok {
					return nil
				}

				totalSeqLength += gcRes.TotalSequenceLen
				totalGCCount += gcRes.GCContentLen
			}
		}
	})

	if err := g.Wait(); err != nil {
		return 0, err
	}

	if totalSeqLength == 0 {
		return 0, nil
	}

	x := (float64(totalGCCount) / float64(totalSeqLength)) * 100
	gcContent := utils.Round(x, gcPrecision)

	return gcContent, nil
}

// readDNASequences reads sequences from the file and sends them to the gcCount channel.
func (s *DNAService) readDNASequences(
	ctx context.Context,
	gcCountCh chan<- string,
	file multipart.File,
) error {
	defer close(gcCountCh)
	defer file.Close()

	var (
		totalGCCount = 0
		scanner      = bufio.NewScanner(file)
	)

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}
			return nil
		}

		sequence := scanner.Text()

		if strings.HasPrefix(sequence, ">") {
			continue
		}

		totalGCCount += len(sequence)
		gcCountCh <- sequence
	}
}

// gcCount processes sequences and sends the GC count results to the gcCountRes channel.
func (s *DNAService) gcCount(
	ctx context.Context,
	gcCountCh <-chan string,
	gcCountRes chan<- GCResult,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sequence, ok := <-gcCountCh:
			if !ok {
				return nil
			}

			gc := gcContentCount(sequence)
			totalSequenceLen := 0
			for _, base := range sequence {
				switch unicode.ToUpper(base) {
				case 'G', 'C', 'T', 'A':
					totalSequenceLen++
				}
			}
			gcCountRes <- GCResult{
				TotalSequenceLen: totalSequenceLen,
				GCContentLen:     gc,
			}
		}
	}
}

// gcContentCount counts the number of G and C bases in a DNA sequence.
func gcContentCount(sequence string) int {
	var gcCount int
	for _, base := range sequence {
		switch unicode.ToUpper(base) {
		case 'G', 'C':
			gcCount++
		}
	}
	return gcCount
}
