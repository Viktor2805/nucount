package dna_test

import (
	"bytes"
	"mime/multipart"
	"strings"
	"testing"

	dna "golang/internal/services/dnaBaseCounter"
)

// nopCloser wraps bytes.Reader to satisfy multipart.File (adds Close()).
type nopCloser struct{ *bytes.Reader }

func (n nopCloser) Close() error { return nil }

func asMultipartFile(s string) multipart.File {
	return nopCloser{bytes.NewReader([]byte(s))}
}

func TestCountBases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		want     dna.BasesCount
		wantErr  bool
	}{
		{
			name: "simple FASTA with headers and newlines",
			input: ">seq1\nACGTACGTNN\n>seq2\nGGCC\n",
			want: dna.BasesCount{A: 2, C: 4, G: 4, T: 2},
		},
		{
			name:    "empty file",
			input:   "",
			want:    dna.BasesCount{A: 0, C: 0, G: 0, T: 0},
			wantErr: false,
		},
		{
			name: "headers only are ignored",
			input: ">seq1 something\n>seq2\n",
			want: dna.BasesCount{},
		},
		{
			name: "lowercase bases are ignored by current implementation",
			input: "acgtACGT\n",
			want: dna.BasesCount{A: 1, C: 1, G: 1, T: 1},
		},
		{
			name: "long line under 1MB scanner buffer",
			input: func() string {
				n := (1 << 20) - 16
				return strings.Repeat("G", n) + "\n"
			}(),
			want: dna.BasesCount{G: (1 << 20) - 16},
		},
	}

	counter := dna.NewBasesCounter()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file := asMultipartFile(tt.input)
			got, err := counter.CountBases(file)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CountBases() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got.A != tt.want.A || got.C != tt.want.C || got.G != tt.want.G || got.T != tt.want.T {
				t.Fatalf("CountBases() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestCountBases_MultipleLinesPerSeq(t *testing.T) {
	counter := dna.NewBasesCounter()

	input := ">seq1\nACGTAC\nGT\n>seq2 description\nCC\nTTAA\n"
	want := dna.BasesCount{A: 4, C: 4, G: 2, T: 4}

	got, err := counter.CountBases(asMultipartFile(input))
	if err != nil {
		t.Fatalf("CountBases() error: %v", err)
	}
	if got != want {
		t.Fatalf("CountBases() = %+v, want %+v", got, want)
	}
}
