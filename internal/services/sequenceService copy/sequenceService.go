package dna

import (
	"bufio"
	"fmt"
	"golang/internal/models"
	repository "golang/internal/repository/sequence"
	"golang/internal/utils"
	"io"
	"mime/multipart"
	"runtime"
	"strings"
	"sync"
)

type SequenceServiceI interface {
	Create(sequence *models.Sequence) error
}

// gcPrecision defines rounding precision for GC content
var gcPrecision = 1

// BasesCount holds base counts (G, C, T, A)
type BasesCount struct {
	G, C, T, A int
}

// BasesCounter is used to count bases in a DNA sequence.
type BasesCounter struct {
	repo repository.SequenceRepositoryI
}

// NewBasesCounter returns a new instance of BasesCounter.
func NewBasesCounter(repo repository.SequenceRepositoryI) *BasesCounter {
	return &BasesCounter{repo: repo}
}

// CountBases reads file in chunks, processes sequences in parallel, and aggregates results
func (s *BasesCounter) CountBases(file multipart.File) (float64, error) {
	defer file.Close()

	reader := bufio.NewReaderSize(file, 1024*1024) // 1MB buffer for faster reading
	linesChan := make(chan string, 500)            // Buffered channel for sequences
	countsChan := make(chan BasesCount, 50)        // Channel for base count results
	var wg sync.WaitGroup

	// runtime.BlockProfile()

	// Memory-efficient buffer reusehttps://ua.iherb.com/pr/now-foods-probiotic-10-100-billion-30-veg-capsules/64465?gad_source=1&gclid=Cj0KCQjw7dm-BhCoARIsALFk4v9ZUI38-ATgsMl6jP14n_ap0y1yZU6U3xD9Fv_kWzUVSpaPqB1uNgkaAt8tEALw_wcB&gclsrc=aw.ds
	bufferPool := sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 1024*64) // 64KB chunk size
			return &buf
		},
	}

	// Worker Goroutines for counting bases
	numWorkers := runtime.NumCPU() * 2
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for sequence := range linesChan {
				countsChan <- s.countSeqBases(sequence)
			}
		}()
	}

	// Read file and push sequences to workers
	go func() {
		defer close(linesChan)
		var builder strings.Builder

		for {
			buf := bufferPool.Get().(*[]byte) // Get buffer from pool
			n, err := reader.Read(*buf)
			if n > 0 {
				chunk := string((*buf)[:n]) // Convert bytes to string
				for _, line := range strings.Split(chunk, "\n") {
					if !strings.HasPrefix(line, ">") { // Skip FASTA headers
						builder.WriteString(line) // Append to sequence
					} else if builder.Len() > 0 { // If full sequence ready
						linesChan <- builder.String()
						builder.Reset()
					}
				}
			}
			bufferPool.Put(buf) // Return buffer to pool

			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error reading file:", err)
				break
			}
		}
	}()

	// Close countsChan once workers finish processing
	go func() {
		wg.Wait()
		close(countsChan)
	}()

	// Aggregate results
	totalCounts := BasesCount{}
	for counts := range countsChan {
		totalCounts.G += counts.G
		totalCounts.C += counts.C
		totalCounts.T += counts.T
		totalCounts.A += counts.A
	}

	// Compute GC Content
	totalGCCount := totalCounts.G + totalCounts.C
	totalSeqLength := totalGCCount + totalCounts.A + totalCounts.T

	if totalSeqLength == 0 {
		return 0, nil
	}

	gcPercentage := (float64(totalGCCount) / float64(totalSeqLength)) * 100
	return utils.Round(gcPercentage, gcPrecision), nil
}

// countSeqBases counts nucleotides in a sequence
func (s *BasesCounter) countSeqBases(sequence string) BasesCount {
	counts := BasesCount{}
	for i := 0; i < len(sequence); i++ {
		switch sequence[i] {
		case 'G':
			counts.G++
		case 'C':
			counts.C++
		case 'T':
			counts.T++
		case 'A':
			counts.A++
		}
	}
	return counts
}
