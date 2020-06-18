package app

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

func TestGetiOSVersion(t *testing.T) {
	cases := []struct {
		client    *http.Client
		expectStr string
		expectErr error
	}{
		{
			NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{ "results": [{ "version": "1.23.4" }] }`)),
				}
			}),
			"1.23.4",
			nil,
		},
		{
			NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 500,
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				}
			}),
			"",
			errors.New("unexpected end of JSON input"),
		},
	}
	for _, c := range cases {
		app := App{c.client}
		str, err := app.GetiOSVersion("123456789")
		if err != nil && err.Error() != c.expectErr.Error() {
			t.Errorf("invalid error expect: %s, actual: %s", c.expectErr, err)
		}
		if str != c.expectStr {
			t.Errorf("invalid error str: %s, actual: %s", c.expectStr, str)
		}
	}
}

func TestApp_GetAndroidVersion(t *testing.T) {
	tests := []struct {
		client    *http.Client
		want    string
		wantErr bool
	}{
		{
			NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(htmlMock)),
				}
			}),
			"4.14.2",
			false,
		},
		{
			NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 500,
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				}
			}),
			"",
			true,
		},
	}
	for _, tt := range tests {
		app := App{ Client: tt.client }
		got, err := app.GetAndroidVersion("com.test")
		if (err != nil) != tt.wantErr {
			t.Errorf("App.GetAndroidVersion() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			t.Errorf("App.GetAndroidVersion() = %v, want %v", got, tt.want)
		}
	}
}

const htmlMock = `<!DOCTYPE html>
<html lang="en">
   <head>
      <meta charset="utf-8">
      <title>title</title>
      <link rel="stylesheet" href="style.css">
      <script src="script.js"></script>
   </head>
   <body>
      <div class="hAyfc"><div class="BgcNfc">現在のバージョン</div><span class="htlgb"><div class="IQ1z0d"><span class="htlgb">4.14.2</span></div></span></div>
   </body>
</html>`