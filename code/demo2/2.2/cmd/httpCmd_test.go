package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleHttp(t *testing.T) {
	usageMessage := `
http: A HTTP client.
		
http: <options> server

  -verb string
    	HTTP method (default "GET")
`
	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		//在未指定位置参数的情况下调用http子命令时测试行为
		{
			args: []string{},
			err:  ErrNoServerSpecified,
		},
		//使用"-h" 调用http子命令时测试行为
		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: usageMessage,
		},
		//使用指定服务器URL的位置参数调用http子命令时的测试行为
		{
			args:   []string{"http://localhost"},
			err:    nil,
			output: "Executing http command\n",
		},
	}
	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleHttp(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}
		if tc.err != nil && err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("expected error %v, got %v", tc.err, err)
		}
		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Errorf("Expected output to be: %#v, Got: %#v", tc.output, gotOutput)
			}
		}
		byteBuf.Reset()
	}
}
