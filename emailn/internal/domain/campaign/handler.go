package campaign

import (
	"emailn/cmd/api/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ServiceInterface
}

func NewCampaignHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleCreate(c *gin.Context) *handlers.ResponseConfig[*CampaignResponse] {
	campaign := c.MustGet("data").(CreateCampaignRequest)

	createdCampaign, err := h.service.Create(&campaign)

	return &handlers.ResponseConfig[*CampaignResponse]{
		SuccessStatus: http.StatusCreated,
		ErrorStatus:   http.StatusBadRequest,
		Data:          createdCampaign,
		Error:         err,
	}
}

func (h *Handler) HandleGet(c *gin.Context) *handlers.ResponseConfig[*[]CampaignResponse] {
	campaigns, err := h.service.Get()

	return &handlers.ResponseConfig[*[]CampaignResponse]{
		SuccessStatus: http.StatusOK,
		ErrorStatus:   http.StatusBadRequest,
		Data:          campaigns,
		Error:         err,
	}
}

func (h *Handler) HandleShow(c *gin.Context) *handlers.ResponseConfig[*CampaignResponse] {
	params := c.MustGet("path").(ShowCampaignParams)

	campaign, err := h.service.Show(params.Id)

	return &handlers.ResponseConfig[*CampaignResponse]{
		SuccessStatus: http.StatusOK,
		ErrorStatus:   http.StatusBadRequest,
		Data:          campaign,
		Error:         err,
	}
}
