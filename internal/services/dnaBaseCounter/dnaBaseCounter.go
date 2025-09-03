// Package services provides functionality for analyzing DNA sequences.
package dna

import (
	"bufio"
	"mime/multipart"
	_ "net/http/pprof"
)

type BasesCounterServiceI interface {
	CountBases(file multipart.File) (BasesCount, error)
}

// BasesCount holds the counts of bases.
type BasesCount struct {
	G, C, T, A int
}

// BasesCounter is used to count the bases in a DNA sequence.
type BasesCounter struct{}

func NewBasesCounter() *BasesCounter {
	return &BasesCounter{}
}

// CountBases counts the GC content in the uploaded file.
func (s *BasesCounter) CountBases(file multipart.File) (BasesCount, error) {
	defer file.Close()
	scanner := bufio.NewScanner(file)

	buf := make([]byte, 1024*1024)     
	scanner.Buffer(buf, 1024*1024)     
	total := BasesCount{}

	for scanner.Scan() {
		sequence := scanner.Bytes()

		if len(sequence) > 0 && sequence[0] == '>' {
			continue
		}

		for _, b := range sequence {
			switch b {
			case 'G', 'g':
					total.G++
			case 'C', 'c':
					total.C++
			case 'T', 't':
					total.T++
			case 'A', 'a':
					total.A++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return BasesCount{}, err
	}
	return total, nil
}


