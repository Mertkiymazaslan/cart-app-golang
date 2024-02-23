package validator

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
)

func msgForValidationError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"

	case "min":
		return "This fields minimum value is " + fe.Param()

	case "max":
		return "This fields maximum value is " + fe.Param()
	}

	return fe.Error()
}

func GetValidatorMessages(err error) error {
	var ve validator.ValidationErrors

	isSuccess := errors.As(err, &ve)
	if !isSuccess {
		return err
	}

	out := make(map[string]string)
	for _, fe := range ve {
		out[fe.Field()] = msgForValidationError(fe)
	}

	jsonData, _ := json.Marshal(out)
	errString := string(jsonData)
	return errors.New(errString)
}
