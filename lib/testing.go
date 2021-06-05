// build +testing
package lib

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

type roundTripFunc func(req *http.Request) (int, string)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	status, in := f(req)
	resp := &http.Response{
		StatusCode:    status,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: int64(len(in)),
	}
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(in))
	return resp, nil
}

func NewTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func TestOnline(t *testing.T) {
	if os.Getenv("TEST_ONLINE") == "" {
		t.Skipf("Skipping this test, please define the env variable TEST_ONLINE")
	}
}
