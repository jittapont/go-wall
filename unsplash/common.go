package unsplash

import (
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func newClient(minTimeout, maxTimeout time.Duration, retry int) (*retryablehttp.Client, error) {
	c := retryablehttp.NewClient()
	c.RetryWaitMin = minTimeout
	c.RetryWaitMax = maxTimeout
	c.RetryMax = retry
	return c, nil
}

func newRequest(u *url.URL, method string, header http.Header) (*retryablehttp.Request, error) {
	req := http.Request{
		Method: method,
		URL:    u,
		Header: header,
	}
	r, err := retryablehttp.FromRequest(&req)
	if err != nil {
		return nil, err
	}
	return r, err
}
