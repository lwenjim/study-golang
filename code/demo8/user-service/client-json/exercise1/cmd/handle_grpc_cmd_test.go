package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleGrpc(t *testing.T) {
	usageMessage := `
grpc: A gRPC client.

grpc: <options> server
Options: 
  -method string
    	Method to call
  -pretty-print
    	Pretty print the JSON output
  -request string
    	Request to send
  -service string
    	gRPC service to send the request to
`
	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args: []string{"-service", "Users", "localhost:8880"},
			err:  errors.New("Invalid gRPC method"),
		},
		{
			args: []string{},
			err:  InvalidInputError{ErrNoServerSpecified},
		},
		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: usageMessage,
		},
	}

	byteBuf := new(bytes.Buffer)
	for i, tc := range testConfigs {
		t.Log(i)
		err := HandleGrpc(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}
		if tc.err != nil && err == nil {
			t.Fatalf("Expected non-nil error, got nil")
		}
		if tc.err != nil && err != nil && tc.err.Error() != err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}
		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Fatalf("Expected output to be: %#v, Got: %#v", tc.output, gotOutput)
			}
		}
		byteBuf.Reset()
	}
}
