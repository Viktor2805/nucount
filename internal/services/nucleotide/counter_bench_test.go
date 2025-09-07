package nucleotide_test

import (
	nucleotide "golang/internal/services/nucleotide"
	"strings"
	"testing"
)

func BenchmarkCountNucleotides_Small(b *testing.B) {
	counter := nucleotide.NewCounter()
	input := ">seq1\nACGTACGT\n"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := counter.Count(asMultipartFile(input))
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkCountNucleotides_1MB(b *testing.B) {
	counter := nucleotide.NewCounter()
	input := strings.Repeat("ACGT", (1<<20)/4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := counter.Count(asMultipartFile(input))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCountNucleotides_10MB(b *testing.B) {
	counter := nucleotide.NewCounter()
	input := strings.Repeat("ACGT", (10<<20)/4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := counter.Count(asMultipartFile(input))
		if err != nil {
			b.Fatal(err)
		}
	}
}
