package _errors

import "errors"

const (
	InternalErrType        = "INTERNAL"
	InvalidInputErrType    = "INVALID_INPUT"
	ResponseReadingErrType = "READING"
	ParsingErrType         = "PARSING"
	GatewayErrType         = "GATEWAY"
)

type ErrorInformation struct {
	Type        string
	Description string
}

func (e ErrorInformation) ErrorInstannce() error {
	return errors.New(e.Description)
}
