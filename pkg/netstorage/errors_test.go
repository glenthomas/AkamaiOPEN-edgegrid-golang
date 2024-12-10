package netstorage

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestNewError(t *testing.T) {
	sess, err := session.New()
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodHead,
		"/",
		nil)
	require.NoError(t, err)

	tests := map[string]struct {
		response *http.Response
		expected *Error
	}{
		"valid response, status code 500": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body: ioutil.NopCloser(strings.NewReader(
					`{"type":"a","title":"b","detail":"c"}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:   "a",
				Title:  "b",
				Detail: "c",
				Status: http.StatusInternalServerError,
			},
		},
		"invalid response body, assign status code": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body: ioutil.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:  "Failed to unmarshal error body. NetStorage API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*netstorage).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestJsonErrorUnmarshalling(t *testing.T) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodHead,
		"/",
		nil)
	require.NoError(t, err)
	tests := map[string]struct {
		input    *http.Response
		expected *Error
	}{
		"API failure with HTML response": {
			input: &http.Response{
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       ioutil.NopCloser(strings.NewReader(`<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>`))},
			expected: &Error{
				Type:   "",
				Title:  "Failed to unmarshal error body. NetStorage API failed. Check details for more information.",
				Detail: "<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>",
				Status: http.StatusServiceUnavailable,
			},
		},
		"API failure with plain text response": {
			input: &http.Response{
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       ioutil.NopCloser(strings.NewReader("Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z"))},
			expected: &Error{
				Type:   "",
				Status: http.StatusServiceUnavailable,
				Title:  "Failed to unmarshal error body. NetStorage API failed. Check details for more information.",
				Detail: "Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z",
			},
		},
		"API failure with XML response": {
			input: &http.Response{
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       ioutil.NopCloser(strings.NewReader(`<Root><Item id="1" name="Example" /></Root>`))},
			expected: &Error{
				Type:   "",
				Status: http.StatusServiceUnavailable,
				Title:  "Failed to unmarshal error body. NetStorage API failed. Check details for more information.",
				Detail: "<Root><Item id=\"1\" name=\"Example\" /></Root>",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			p := netstorage{
				Session: sess,
			}
			assert.Equal(t, test.expected, p.Error(test.input))
		})
	}
}
