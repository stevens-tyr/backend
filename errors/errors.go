package errors

import (
	"errors"
	"net/http"
)

type APIError interface {
	Error() string
	GetError() error
	StatusCode() int
}

type Error struct {
	Err error
	SC  int
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) GetError() error {
	return e.Err
}

func (e *Error) StatusCode() int {
	return e.SC
}

var (
	// UserNotFoundError an error to throw for when a User is not found.
	ErrorResourceNotFound = &Error{errors.New("RESOURCE DOES NOT EXIST"), http.StatusNotFound}
	// IncorrectPasswordCredentials an error to throw for when login credentials are incorrect.
	ErrorIncorrectCredentials = &Error{errors.New("INCORRECT CREDENTIALS"), http.StatusUnauthorized}
	// ErrorNonMatchingPassword an error to throw when a password cofirmation does not match the password.
	ErrorNonMatchingPassword = &Error{errors.New("CONFIRMATION MUST MATCH"), http.StatusBadRequest}
	// ErrorFailedToCreateUser an error for when you fail to create a user.
	ErrorDatabaseFailedCreate      = &Error{errors.New("DATABASE CREATE OPERATION FAILURE"), http.StatusInternalServerError}
	ErrorDatabaseFailedUpdate      = &Error{errors.New("DATABASE UPDATE OPERATION FAILURE"), http.StatusInternalServerError}
	ErrorDatabaseFailedQuery       = &Error{errors.New("DATABASE QUERY OPERATION FAILURE"), http.StatusInternalServerError}
	ErrorDatabaseFailedExtract     = &Error{errors.New("DATABASE EXTRACT DATA OPERATION FAILURE"), http.StatusInternalServerError}
	ErrorCannotCreateDuplicateData = &Error{errors.New("CANNOT CREATE DUPLICATE DATABASE ENTRY"), http.StatusConflict}
	ErrorHashFailure               = &Error{errors.New("FAILED TO HASH CONTENT"), http.StatusInternalServerError}
	ErrorUserAlreadyEnrolled       = &Error{errors.New("USER ALREADY ENROLLED IN COURSE"), http.StatusConflict}
	ErrorInvalidObjectID           = &Error{errors.New("INVALID OBJECT ID"), http.StatusBadRequest}
	ErrorInvlaidJSON               = &Error{errors.New("INVALID JSON"), http.StatusBadRequest}
	ErrorInvlaidBSON               = &Error{errors.New("INVALID BSON OBJECT DECODED"), http.StatusInternalServerError}
	ErrorGenerateTokenFailure      = &Error{errors.New("GENERATE TOKEN FAILURE"), http.StatusInternalServerError}
)
