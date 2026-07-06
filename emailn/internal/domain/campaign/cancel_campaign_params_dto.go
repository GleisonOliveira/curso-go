package campaign

import "emailn/internal/types"

type CancelCampaignParams struct {
	Id types.UUID `uri:"id" binding:"required" json:"id"`
}
