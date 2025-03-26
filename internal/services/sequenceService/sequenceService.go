package sequence

import (
	"bufio"
	repository "golang/internal/repository/sequence"
	"mime/multipart"
	"runtime"
	"strings"
	"sync"
)

// SequenceServiceI defines the interface for sequence processing.
type SequenceServiceI interface {
	ProcessFASTAFromMultipart(file multipart.File, windowSize, stepSize int) (map[string][]float64, error)
}

// SequenceService handles FASTA file processing.
type SequenceService struct {
	repo repository.SequenceRepositoryI
}

// NewSequenceService creates a new SequenceService.
func NewSequenceService(repo repository.SequenceRepositoryI) *SequenceService {
	return &SequenceService{repo: repo}
}

// ProcessFASTAFromMultipart streams a FASTA file and returns GC skew results **per window**.
func (s *SequenceService) ProcessFASTAFromMultipart(file multipart.File, windowSize, stepSize int) (map[string][]float64, error) {
	reader := bufio.NewReader(file)
	gcSkewResults := make(map[string][]float64)
	var accession string
	seqBuffer := make([]byte, 0, windowSize*2) // Preallocate rolling buffer

	// Detect number of CPU cores
	numThreads := runtime.NumCPU()

	// Create worker pool for parallel processing
	taskQueue := make(chan struct {
		Accession string
		Sequence  []byte
	}, numThreads)

	resultsQueue := make(chan struct {
		Accession string
		SkewData  []float64
	}, numThreads)

	// Worker pool to process tasks
	var wg sync.WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskQueue {
				resultsQueue <- struct {
					Accession string
					SkewData  []float64
				}{Accession: task.Accession, SkewData: computeGCSkewPerWindow(task.Sequence, windowSize, stepSize)}
			}
		}()
	}

	// Separate Goroutine to read results
	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for result := range resultsQueue {
			gcSkewResults[result.Accession] = append(gcSkewResults[result.Accession], result.SkewData...)
		}
	}()

	// Read and process FASTA file
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // EOF
		}
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, ">") { // New sequence header
			if len(seqBuffer) > 0 {
				taskQueue <- struct {
					Accession string
					Sequence  []byte
				}{Accession: accession, Sequence: append([]byte{}, seqBuffer...)}
				seqBuffer = seqBuffer[:0]
			}

			// Extract Accession ID (ignore description)
			parts := strings.Fields(line[1:])
			if len(parts) > 0 {
				accession = parts[0]
			}
		} else {
			// Append new sequence data
			seqBuffer = append(seqBuffer, line...)

			// Process in rolling windows once enough data is available
			for len(seqBuffer) >= windowSize {
				taskQueue <- struct {
					Accession string
					Sequence  []byte
				}{Accession: accession, Sequence: seqBuffer[:windowSize]}
				seqBuffer = seqBuffer[stepSize:]
			}
		}
	}

	// Process remaining sequence chunk
	if len(seqBuffer) > 0 {
		taskQueue <- struct {
			Accession string
			Sequence  []byte
		}{Accession: accession, Sequence: append([]byte{}, seqBuffer...)}
	}

	// Close channels and wait for completion
	close(taskQueue)
	wg.Wait()
	close(resultsQueue)
	resultsWg.Wait()

	return gcSkewResults, nil
}
func computeGCSkewPerWindow(sequence []byte, windowSize, stepSize int) []float64 {
	gCount, cCount := 0, 0
	length := len(sequence)
	skewValues := make([]float64, 0, length/windowSize) // Preallocate

	for i := 0; i < length; i++ {
		if sequence[i] == 'G' {
			gCount++
		} else if sequence[i] == 'C' {
			cCount++
		}

		if (i+1)%stepSize == 0 { // Store **one value per window**
			if gCount+cCount > 0 {
				skewValues = append(skewValues, float64(gCount-cCount)/float64(gCount+cCount))
			} else {
				skewValues = append(skewValues, 0.0)
			}
			gCount, cCount = 0, 0
		}
	}

	if len(skewValues) == 0 { // Ensure at least one value is returned
		skewValues = append(skewValues, 0.0)
	}

	return skewValues
}
