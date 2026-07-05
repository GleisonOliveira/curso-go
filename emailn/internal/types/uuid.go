package types

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type UUID uuid.UUID

func (u *UUID) UnmarshalParam(param string) error {
	parsed, err := uuid.Parse(param)

	if err != nil {
		return err
	}

	*u = UUID(parsed)

	return nil
}

var _ binding.BindUnmarshaler = (*UUID)(nil)
