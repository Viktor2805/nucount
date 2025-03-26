package models

// Sequence represents a genomic sequence stored in the database.
type Sequence struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Accession   string  `gorm:"not null;uniqueIndex" json:"accession"`     // Unique GenBank ID
	Species     *string `gorm:"default:null" json:"species,omitempty"`     // Nullable if unknown
	Chromosome  *string `gorm:"default:null" json:"chromosome,omitempty"`  // Nullable scaffold/chromosome
	Description *string `gorm:"default:null" json:"description,omitempty"` // Metadata, optional
	Sequence    string  `gorm:"not null" json:"sequence"`                  // Raw DNA sequence
	StartPos    int     `gorm:"default:1" json:"start_position"`           // Always starts at 1
	EndPos      int     `gorm:"-" json:"end_position"`                     // Will be computed dynamically
	CreatedAt   int64   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64   `gorm:"autoUpdateTime" json:"updated_at"`
}
