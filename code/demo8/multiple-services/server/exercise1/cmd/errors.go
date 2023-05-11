package cmd

import "errors"

var ErrNoServerSpecified = errors.New("you have to specify the remote server")
var ErrInvalidHTTPMethod = errors.New("invalid HTTP method")
var ErrInvalidHTTPPostCommand = errors.New("cannot specify both body and body-file")
var ErrInvalidHTTPPostRequest = errors.New("HTTP POST request must specify a non-empty JSON body")
var ErrInvalidHTTPCommand = errors.New("invalid HTTP command")
var ErrInvalidGrpcMethod = errors.New("invalid gRPC method")

type FlagParsingError struct {
	err error
}

func (e FlagParsingError) Error() string {
	return e.err.Error()
}

type InvalidInputError struct {
	Err error
}

func (e InvalidInputError) Error() string {
	return e.Err.Error()
}
