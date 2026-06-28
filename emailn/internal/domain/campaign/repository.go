package campaign

import "github.com/google/uuid"

type Repository interface {
	Save(campaign *Campaign) error
	Get() (*[]Campaign, error)
	Show(*uuid.UUID) (*Campaign, error)
}
