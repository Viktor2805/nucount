// Package nucleotide provides functionality for analyzing DNA sequences.
package nucleotide

import (
	"io"
	"mime/multipart"
	_ "net/http/pprof"
)

type Service interface {
	Count(file multipart.File) (NucleotideCount, error)
}

// NucleotideCount represents raw counts of A, C, G, and T
type NucleotideCount struct {
	A int
	C int
	G int
	T int
}

// Total returns the total number of nucleotides counted.
func (n NucleotideCount) Total() int {
	return n.A + n.C + n.G + n.T
}

// NucleotideCounter is used to count nucleotides in a DNA sequence.
type Counter struct{}

func NewCounter() Service {
	return &Counter{}
}

// CountNucleotides counts the nucleotides in the uploaded FASTA file.
func (s *Counter) Count(file multipart.File) (NucleotideCount, error) {
	defer file.Close()

	const bufSize = 4 << 20
	buf := make([]byte, bufSize)

	var lutA, lutC, lutG, lutT [256]uint8
	lutA['A'], lutA['a'] = 1, 1
	lutC['C'], lutC['c'] = 1, 1
	lutG['G'], lutG['g'] = 1, 1
	lutT['T'], lutT['t'] = 1, 1

	var count NucleotideCount
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
					i++
					continue
				}

				j := i
				for j < len(p) && p[j] != '\n' {
					j++
				}
				segment := p[i:j]

				for _, b := range segment {
					count.A += int(lutA[b])
					count.C += int(lutC[b])
					count.G += int(lutG[b])
					count.T += int(lutT[b])
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
			return count, err
		}
	}

	return count, nil
}
