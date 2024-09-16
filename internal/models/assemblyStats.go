package models

import (
	"time"
)

type AssemblyStats struct {
	ID                  int       `json:"id" db:"id"`
	AssemblyID          int       `json:"assembly_id" db:"assembly_id"`
	ContigL50           int       `json:"contig_l50" db:"contigl50"`
	ContigN50           int       `json:"contig_n50" db:"contign50"`
	TotalSequenceLength int64     `json:"total_sequence_length" db:"total_sequence_length"`
	TotalUngappedLength int64     `json:"total_ungapped_length" db:"total_ungapped_length"`
	GCCount             int64     `json:"gc_count" db:"gc_count"`
	GCPercent           float64   `json:"gc_percent" db:"gc_percent"`
	GenomeCoverage      string    `json:"genome_coverage" db:"genome_coverage"`
	NumberOfContigs     int       `json:"number_of_contigs" db:"number_of_contigs"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}