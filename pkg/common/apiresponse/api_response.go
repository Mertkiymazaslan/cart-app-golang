package apiresponse

import (
	errs "checkoutProject/pkg/common/errors"
	"errors"
	"net/http"
)

type Responder interface {
	Response() interface{}
}

func OK(r Responder) (int, interface{}) {
	return http.StatusOK, r.Response()
}

func Created(r Responder) (int, interface{}) {
	return http.StatusCreated, r.Response()
}

func Failed(err error) (int, interface{}) {
	responseCode := http.StatusBadRequest
	genericResponse := GenericResponse{Result: false, Message: err.Error()}

	if errors.Is(err, errs.InternalServerErr) {
		responseCode = errs.InternalServerErrCode
	}

	if errors.Is(err, errs.RecordNotFoundErr) {
		responseCode = errs.RecordNotFoundErrCode
	}

	return responseCode, genericResponse
}
