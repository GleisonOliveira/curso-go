package internalerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessError_NotFound(t *testing.T) {
	err := &ErrNotFound{Entity: "campaign"}
	result := ProcessError(err)

	var notFoundErr *ErrNotFound
	assert.True(t, errors.As(result, &notFoundErr))
	assert.Equal(t, "campaign", notFoundErr.Entity)
}

func TestProcessError_Internal(t *testing.T) {
	result := ProcessError(ErrInternal)
	assert.ErrorIs(t, result, ErrInternal)
}

func TestProcessError_Unknown(t *testing.T) {
	result := ProcessError(errors.New("some error"))
	assert.ErrorIs(t, result, ErrInternal)
}
