package slack

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	testutil "github.com/murakamiii/go-check-app-release/testutil"
)

func TestPostMessage(t *testing.T) {
	cases := []struct {
		client    *http.Client
		expectErr error
	}{
		{
			testutil.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				}
			}),
			nil,
		},
		{
			testutil.NewTestClient(func(req *http.Request) *http.Response {
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
