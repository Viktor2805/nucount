package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxUploadSizeMiddleware restricts the maximum size of the request body.
func MaxUploadSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)

		if err := ctx.Request.ParseMultipartForm(maxSize); err != nil {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File size exceeds the maximum allowed limit."})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
