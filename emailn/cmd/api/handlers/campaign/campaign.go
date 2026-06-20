package campaign

import (
	"emailn/internal/contract/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handle(c *gin.Context) {
	campaign := c.MustGet("data").(dto.NewCampaign)

	c.JSON(http.StatusOK, campaign)
}
