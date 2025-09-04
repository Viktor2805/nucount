// Package services provides functionality for analyzing DNA sequences.
package dna

import (
	"io"
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

	const bufSize = 4 << 20 
	buf := make([]byte, bufSize)

	var lutA, lutC, lutG, lutT [256]uint8
	lutA['A'], lutA['a'] = 1, 1
	lutC['C'], lutC['c'] = 1, 1
	lutG['G'], lutG['g'] = 1, 1
	lutT['T'], lutT['t'] = 1, 1

	var total BasesCount

	inHeader := false

	for {
		n, err := file.Read(buf)
		if n > 0 {
			p := buf[:n]

			i := 0
			for i < len(p) {
				if inHeader {
					j := i
					for j < len(p) && p[j] != '\n' {
						j++
					}
					if j < len(p) && p[j] == '\n' {
						inHeader = false
						i = j + 1
						continue
					}
					break
				}

				if p[i] == '>' && (i == 0 || p[i-1] == '\n') {
					inHeader = true
					continue
				}

				j := i
				for j < len(p) && p[j] != '\n' {
					j++
				}
				seg := p[i:j]

				for _, b := range seg {
					total.A += int(lutA[b])
					total.C += int(lutC[b])
					total.G += int(lutG[b])
					total.T += int(lutT[b])
				}

				if j < len(p) && p[j] == '\n' {
					j++
				}
				i = j
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

