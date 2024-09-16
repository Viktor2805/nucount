package dnaSeq

import (
	"math"
)

type AlignmentResult struct {
	AlignedSeq1 string
	AlignedSeq2 string
	Score       int // Optional, may be 0 if not used
}

type AlignmentAlgorithm interface {
	Align(seq1, seq2 string) AlignmentResult
}

// NeedlemanWunschAlgorithm implements the Needleman-Wunsch algorithm for global alignment.
type NeedlemanWunschAlgorithm struct {
	MatchScore    int
	MismatchScore int
	GapPenalty    int
}

// NewNeedlemanWunschAlgorithm creates a new instance of NeedlemanWunschAlgorithm.
func NewNeedlemanWunschAlgorithm() *NeedlemanWunschAlgorithm {
	return &NeedlemanWunschAlgorithm{
		MatchScore:    2,
		MismatchScore: -1,
		GapPenalty:    -2,
	}
}

// Align performs global alignment using the Needleman-Wunsch algorithm.
func (nwa *NeedlemanWunschAlgorithm) Align(seq1, seq2 string) (alignedSeq1, alignedSeq2 string, score int) {
	matrix := initializeMatrix(len(seq1), len(seq2))
	fillMatrix(matrix, seq1, seq2, nwa.MatchScore, nwa.MismatchScore, nwa.GapPenalty)
	alignedSeq1, alignedSeq2 = traceback(matrix, seq1, seq2, nwa.GapPenalty)
	score = matrix[len(seq1)][len(seq2)]
	return
}

// Functions for initializing the matrix, filling the matrix, and traceback
func initializeMatrix(lenX, lenY int) [][]int {
	matrix := make([][]int, lenX+1)
	for i := range matrix {
		matrix[i] = make([]int, lenY+1)
	}
	return matrix
}

func fillMatrix(matrix [][]int, seq1, seq2 string, matchScore, mismatchScore, gapPenalty int) {
	for i := 1; i <= len(seq1); i++ {
		matrix[i][0] = i * gapPenalty
	}
	for j := 1; j <= len(seq2); j++ {
		matrix[0][j] = j * gapPenalty
	}

	for i := 1; i <= len(seq1); i++ {
		for j := 1; j <= len(seq2); j++ {
			match := matrix[i-1][j-1] + score(seq1[i-1], seq2[j-1], matchScore, mismatchScore)
			delete := matrix[i-1][j] + gapPenalty
			insert := matrix[i][j-1] + gapPenalty
			matrix[i][j] = max(match, delete, insert)
		}
	}
}

func score(a, b byte, matchScore, mismatchScore int) int {
	if a == b {
		return matchScore
	}
	return mismatchScore
}

func max(a, b, c int) int {
	return int(math.Max(float64(a), math.Max(float64(b), float64(c))))
}

func traceback(matrix [][]int, seq1, seq2 string, gapPenalty int) (string, string) {
	var alignedSeq1, alignedSeq2 string
	i, j := len(seq1), len(seq2)

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && matrix[i][j] == matrix[i-1][j-1]+score(seq1[i-1], seq2[j-1], 2, -1) {
			alignedSeq1 = string(seq1[i-1]) + alignedSeq1
			alignedSeq2 = string(seq2[j-1]) + alignedSeq2
			i--
			j--
		} else if i > 0 && matrix[i][j] == matrix[i-1][j]+gapPenalty {
			alignedSeq1 = string(seq1[i-1]) + alignedSeq1
			alignedSeq2 = "-" + alignedSeq2
			i--
		} else if j > 0 && matrix[i][j] == matrix[i][j-1]+gapPenalty {
			alignedSeq1 = "-" + alignedSeq1
			alignedSeq2 = string(seq2[j-1]) + alignedSeq2
			j--
		}
	}

	return alignedSeq1, alignedSeq2
}
