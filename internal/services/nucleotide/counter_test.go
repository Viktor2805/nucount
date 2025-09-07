package nucleotide_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"strings"
	"testing"

	nucleotide "golang/internal/services/nucleotide"
)

type errAfterFile struct {
	*bytes.Reader
}

func (f errAfterFile) Read(p []byte) (int, error) { return 0, errors.New("read failure") }
func (f errAfterFile) Close() error               { return nil }

type nopCloser struct{ *bytes.Reader }

func (n nopCloser) Close() error { return nil }

func asMultipartFile(s string) multipart.File {
	return nopCloser{bytes.NewReader([]byte(s))}
}

type limitedFile struct {
	*bytes.Reader
	limit int
}

func (lf *limitedFile) Read(p []byte) (int, error) {
	if len(p) > lf.limit {
		p = p[:lf.limit]
	}
	return lf.Reader.Read(p)
}
func (lf *limitedFile) Close() error { return nil }

func TestCountBases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    nucleotide.NucleotideCount
		wantErr bool
	}{
		{
			name:  "simple FASTA with headers and newlines",
			input: ">seq1\nACGTACGTNN\n>seq2\nGGCC\n",
			want:  nucleotide.NucleotideCount{A: 2, C: 4, G: 4, T: 2},
		},
		{
			name:    "empty file",
			input:   "",
			want:    nucleotide.NucleotideCount{A: 0, C: 0, G: 0, T: 0},
			wantErr: false,
		},
		{
			name:  "headers only are ignored",
			input: ">seq1 something\n>seq2\n",
			want:  nucleotide.NucleotideCount{},
		},
		{
			name:  "lowercase bases are ignored by current implementation",
			input: "acgtACGT\n",
			want:  nucleotide.NucleotideCount{A: 2, C: 2, G: 2, T: 2},
		},
	}

	counter := nucleotide.NewCounter()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file := asMultipartFile(tt.input)
			got, err := counter.Count(file)
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
	counter := nucleotide.NewCounter()

	input := ">seq1\nACGTAC\nGT\n>seq2 description\nCC\nTTAA\n"
	want := nucleotide.NucleotideCount{A: 4, C: 4, G: 2, T: 4}

	got, err := counter.Count(asMultipartFile(input))
	if err != nil {
		t.Fatalf("CountBases() error: %v", err)
	}
	if got != want {
		t.Fatalf("CountBases() = %+v, want %+v", got, want)
	}
}

func TestCountBases_ErrorFile(t *testing.T) {
	counter := nucleotide.NewCounter()

	_, err := counter.Count(errAfterFile{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCountBases_HeaderBreak(t *testing.T) {
	counter := nucleotide.NewCounter()

	input := ">" + strings.Repeat("HEADER", 50) + "\nACGT\n"

	f := &limitedFile{
		Reader: bytes.NewReader([]byte(input)),
		limit:  1,
	}

	got, err := counter.Count(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := nucleotide.NucleotideCount{A: 1, C: 1, G: 1, T: 1}
	if got != want {
		t.Fatalf("CountBases() = %+v, want %+v", got, want)
	}
}
