package errors

import (
	"errors"
)

var (
	InvalidGender = errors.New("must be a valid value")
	InvalidData   = errors.New("must be in a valid format")
	EmptyField    = errors.New("cannot be blank")
	InvalidId     = errors.New("invalid id")
	NotExistId    = errors.New("sql: no rows in result set: no such user : failed to get person in service")
)
