package apiresponse

type GenericResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

type GenericResponseSerializer struct {
	Result  bool
	Message string
}

func (s GenericResponseSerializer) Response() interface{} {
	return GenericResponse{
		Result:  s.Result,
		Message: s.Message,
	}
}
