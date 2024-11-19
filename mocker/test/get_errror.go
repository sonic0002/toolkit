package test

import (
	"errors"
)

var getError = func(str string) error {
	return errors.New(str)
}
