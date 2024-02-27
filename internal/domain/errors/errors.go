package errors

import (
	"errors"
)

var (
	InvalidGender = errors.New("must be a valid value")
	InvalidData   = errors.New("must be in a valid format")
	EmptyField    = errors.New("cannot be blank")
)
