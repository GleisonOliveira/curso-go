package campaign

import "emailn/internal/types"

type Repository interface {
	Save(campaign *Campaign) error
	Get() (*[]Campaign, error)
	Show(types.UUID) (*Campaign, error)
}
