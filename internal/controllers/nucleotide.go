package controllers

import (
	"net/http"

	"golang/internal/helpers"
	nucleotide "golang/internal/services/nucleotide"

	"github.com/gin-gonic/gin"
)


type NucleotideCountDTO struct {
	A int `json:"A"`
	C int `json:"C"`
	G int `json:"G"`
	T int `json:"T"`
}

type NucleotideCountResponse struct {
	Bases NucleotideCountDTO `json:"bases"`
	Total int           `json:"total"`
}

type NucleotideController struct {
	service nucleotide.Service
}

func NewNucleotideController(service nucleotide.Service) *NucleotideController {
	return &NucleotideController{service: service}
}

// CountBases uploads a FASTA file and returns A/C/G/T counts.
// @Summary Count nucleotides from a FASTA file
// @Description Upload a .fasta/.fa (or plain text FASTA) file and get counts of A/C/G/T with total.
// @Tags nucleotides
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "FASTA file (.fasta, .fa) or text/plain"
// @Success 200 {object} controllers.NucleotideCountResponse
// @Failure 400 {object} apierror.ErrorResponse "Invalid file or parameters"
// @Failure 500 {object} apierror.ErrorResponse "Internal error while processing file"
// @Router /nucleotides/count [post]
func (c *NucleotideController) Count(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	if !helpers.IsFASTAByExtension(header.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension. Please upload a .fasta or .fa file."})
		return
	}
	
	if !helpers.IsValidFASTAContentType(header) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type. Expected 'application/octet-stream' or 'text/plain'."})
		return
	}

	bc, err := c.service.Count(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := NucleotideCountResponse{
		Bases: NucleotideCountDTO{A: bc.A, C: bc.C, G: bc.G, T: bc.T},
		Total: bc.Total(),    
	}
	ctx.JSON(http.StatusOK, resp)
}
