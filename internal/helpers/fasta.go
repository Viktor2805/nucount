package helpers

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

// Allowed MIME types for FASTA files
const (
	ContentTypeFASTA1 = "application/octet-stream"
	ContentTypeFASTA2 = "text/plain" // Some FASTA files may be treated as text
)

// isFASTAByExtension checks if the file extension is .fasta or .fa
func IsFASTAByExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".fasta" || ext == ".fna"
}

// isValidFASTAContentType checks if the uploaded file has a valid FASTA content type
func IsValidFASTAContentType(header *multipart.FileHeader) bool {
	contentType := header.Header.Get("Content-Type")
	return contentType == ContentTypeFASTA1 || contentType == ContentTypeFASTA2
}
