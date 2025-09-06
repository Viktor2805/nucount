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
	A     int `json:"A"`
	C     int `json:"C"`
	G     int `json:"G"`
	T     int `json:"T"`
}

func (b BasesCount) Total() int {
	return b.A + b.C + b.G + b.T
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

	var basesCount BasesCount

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
					basesCount.A += int(lutA[b])
					basesCount.C += int(lutC[b])
					basesCount.G += int(lutG[b])
					basesCount.T += int(lutT[b])
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
			return basesCount, err
		}
	}

	return basesCount, nil
}

