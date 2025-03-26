package controllers

import (
	sequence "golang/internal/services/sequenceService"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SequenceController struct {
	service sequence.SequenceServiceI
}

func NewSequenceController(service sequence.SequenceServiceI) *SequenceController {
	return &SequenceController{service: service}
}

// UploadCSVFile uploads a CSV file.
// @Summary Upload CSV file
// @Description Uploads a CSV file containing transactions
// @Tags transactions
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "CSV file"
// @Success 200 {object} apierror.SuccessResponse
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 500 {object} apierror.ErrorResponse
// @Router /transactions/upload [post]
func (c *SequenceController) AnalyzeDNASeq(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/octet-stream" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Please upload a TXT file"})
		return
	}

	skewValues, err := c.service.ProcessFASTAFromMultipart(file, 1000, 500)

	// skewValues, err: = c.service P(file, 1000, 500)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"skewValues": skewValues})
}
