package response

import (
	"github.com/iwanlaudin/go-microservice/services/ticket/internal/models"
)

type EventResponse struct {
	Status  string        `json:"status"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Items   *models.Event `json:"items"`
}
