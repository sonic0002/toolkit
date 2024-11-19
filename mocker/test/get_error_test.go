package test

import (
	"errors"
	"mocker"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetError(t *testing.T) {
	defer mocker.RestoreMock(&getError)()

	err := errors.New("error here")
	getError = func(str string) error {
		return err
	}

	assert.Equal(t, err, getError("dummy error"))
}

func TestGetErrorAgain(t *testing.T) {
	errDummy := errors.New("dummy error")

	assert.Equal(t, errDummy, getError("dummy error"))
}
