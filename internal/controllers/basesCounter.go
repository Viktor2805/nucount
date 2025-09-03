package controllers

import (
	"golang/internal/helpers"
	basesCounter "golang/internal/services/dnaBaseCounter"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GCParams binds query parameters for GC skew analysis

const (
	ContentTypeFASTA      = "application/octet-stream"
	ErrInvalidQueryParams = "Invalid query parameters. Please ensure windowSize and stepSize are within valid ranges."
	ErrInvalidFileType    = "Invalid file type. Please upload a valid FASTA/TXT file."
	ErrFileProcessing     = "Error processing the uploaded file."
)

type BasesController struct {
	service basesCounter.BasesCounterServiceI
}

func NewBasesCounterController(service basesCounter.BasesCounterServiceI) *BasesController {
	return &BasesController{service: service}
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
func (c *BasesController) CountBases(ctx *gin.Context) {
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

	basesCount, err := c.service.CountBases(file)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"basesCount": basesCount,
		"total": basesCount.A + basesCount.C + basesCount.G + basesCount.T,
	})
}
