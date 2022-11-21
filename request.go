package fourbyte

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

type request struct {
	url url.URL
}

func newRequest(u url.URL, opts ...filterOption) *request {
	req := &request{u}
	for _, opt := range opts {
		if opt != nil {
			opt.apply(req)
		}
	}
	return req
}

func (r *request) Get(ctx context.Context) (status int, data []byte, err error) {
	return doRequest(ctx, http.MethodGet, r.url.String(), nil, nil)
}

func (r *request) Post(ctx context.Context, contentType string, body []byte) (status int, data []byte, err error) {
	return doRequest(ctx, http.MethodPost, r.url.String(), body, http.Header{"Content-Type": []string{contentType}})
}

func (r *request) addQueryValue(key, val string) {
	vals := r.url.Query()
	vals.Add(key, val)
	r.url.RawQuery = vals.Encode()
}

func (r *request) setQueryValue(key, val string) {
	vals := r.url.Query()
	vals.Set(key, val)
	r.url.RawQuery = vals.Encode()
}

func doRequest(ctx context.Context, method, url string, body []byte, headers http.Header) (status int, data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return -1, nil, err
	}
	req.Header = headers
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode, b, nil
}
