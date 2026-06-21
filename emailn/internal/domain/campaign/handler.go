package campaign

import (
	"emailn/cmd/api/handlers"
	"emailn/internal/contract/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewCampaignHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleCreate(c *gin.Context) (*handlers.ResponseConfig, *Campaign, error) {
	campaign := c.MustGet("data").(dto.NewCampaign)

	createdCampaign, err := h.service.Create(campaign)

	return &handlers.ResponseConfig{
		SuccessStatus: http.StatusCreated,
		ErrorStatus:   http.StatusBadRequest,
	}, createdCampaign, err
}

func (h *Handler) HandleGet(c *gin.Context) (*handlers.ResponseConfig, *[]Campaign, error) {
	campaigns, err := h.service.Get()

	return &handlers.ResponseConfig{
		SuccessStatus: http.StatusOK,
		ErrorStatus:   http.StatusBadRequest,
	}, campaigns, err
}
