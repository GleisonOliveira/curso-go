package dto

type NewCampaign struct {
	Name    string   `json:"name" binding:"required,min=2,max=255"`
	Content string   `json:"content" binding:"required,min=2,max=1000"`
	Emails  []string `json:"emails" binding:"required,min=1,dive,required,email"`
}
