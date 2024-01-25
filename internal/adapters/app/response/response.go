package response

import (
	"github.com/Kosodaka/enricher-service/internal/domain/model"
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

type IdResponse struct {
	Id int `json:"id"`
}

type PersonResponse struct {
	Person model.Person `json:"person"`
}

type PersonsPesponse struct {
	Persons []model.Person `json:"persons"`
}
