package response

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type ErrorResponse struct {
	Message string `json:"message"`
}
type StatusResponse struct {
	Status string `json:"status"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	slog.Debug(message)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
