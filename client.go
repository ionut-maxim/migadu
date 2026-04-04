package migadu

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.migadu.com/v1/"

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func New(user, apiKey string) *Client {
	base, _ := url.Parse(baseURL)
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &roundTripper{
				base:   base,
				next:   &http.Transport{},
				user:   user,
				apiKey: apiKey,
			},
		},
	}
}

type roundTripper struct {
	base   *url.URL
	next   http.RoundTripper
	user   string
	apiKey string
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme = r.base.Scheme
	req.URL.Host = r.base.Host
	req.URL.Path = r.base.Path + req.URL.Path
	req.SetBasicAuth(r.user, r.apiKey)
	req.Header.Add("Content-Type", "application/json")
	return r.next.RoundTrip(req)
}

func doAndDecode[T any](c *Client, req *http.Request) (T, error) {
	var zero T
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return zero, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zero, newError(resp)
	}
	return decode[T](resp)
}

func decode[T any](resp *http.Response) (T, error) {
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}

func encode(v any) (io.Reader, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return &buf, nil
}
