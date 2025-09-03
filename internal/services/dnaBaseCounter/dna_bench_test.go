package dna_test

import (
	dna "golang/internal/services/dnaBaseCounter"
	"strings"
	"testing"
)

func BenchmarkCountBases_Small(b *testing.B)  {
	counter := 	dna.NewBasesCounter()
	input := ">seq1\nACGTACGT\n"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := counter.CountBases(asMultipartFile(input))
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkCountBases_1MB(b *testing.B) {
	counter := dna.NewBasesCounter()
	input := strings.Repeat("ACGT", (1<<20)/4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := counter.CountBases(asMultipartFile(input))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCountBases_10MB(b *testing.B) {
	counter := dna.NewBasesCounter()
	input := strings.Repeat("ACGT", (10<<20)/4) 

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := counter.CountBases(asMultipartFile(input))
		if err != nil {
			b.Fatal(err)
		}
	}
}
