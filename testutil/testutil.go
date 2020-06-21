package testutil

import (
	"net/http"
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
