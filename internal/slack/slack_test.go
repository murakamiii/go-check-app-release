package slack

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

// RoundTripFunc see http://hassansin.github.io/Unit-Testing-http-client-in-Go
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip see http://hassansin.github.io/Unit-Testing-http-client-in-Go
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
// see http://hassansin.github.io/Unit-Testing-http-client-in-Go
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestPostMessage(t *testing.T) {
	cases := []struct {
		client    *http.Client
		expectErr error
	}{
		{
			NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				}
			}),
			nil,
		},
		{
			NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 400,
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				}
			}),
			errors.New("Slack POST Error: HTTP Status 400"),
		},
	}
	for _, c := range cases {
		sl := Slack{c.client}
		err := sl.PostMessage("path/string", "message")
		if err != nil && err.Error() != c.expectErr.Error() {
			t.Errorf("invalid error expect: %s, actual: %s", c.expectErr, err)
		}
	}
}
