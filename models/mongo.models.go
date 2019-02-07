package models

import (
	"errors"
)

var (
	// ErrorEmailNotValid an error to throw when an email format is not valid
	ErrorEmailNotValid = errors.New("EMAIL NOT VALID")
	// ErrorUnresolvedEmailHost an error to throw when the email host is unresolvable
	ErrorUnresolvableEmailHost = errors.New("EMAIL HOST UNRESOLVABLE")
)
