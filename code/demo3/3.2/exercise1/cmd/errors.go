package cmd

import "errors"

var ErrNoServerSpecified = errors.New("ErrNoServerSpecified")
var ErrInvalidHTTPMethod = errors.New("Invalid HTTP method")
