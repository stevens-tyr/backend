package models

import (
	"errors"
)

// gStorage types/structs

var (
	ErrorUnableToWrite = errors.New("UNABLE TO WRITE TO BUCKET")
	ErrorUnableToClose = errors.New("UNABLE TO CLOSE BUCKET")
	ErrorUnableToOpen  = errors.New("UNABLE TO OPEN BUCKET FILE")
	ErrorUnableToRead  = errors.New("UNABLE TO READ FILE")
)
