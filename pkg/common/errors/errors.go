package errors

import (
	"errors"
	"net/http"
)

const (
	RecordNotFoundErrCode = http.StatusNotFound
	InternalServerErrCode = http.StatusInternalServerError
)

var (
	InternalServerErr = errors.New("internal server error")
	RecordNotFoundErr = errors.New("record not found")
)
