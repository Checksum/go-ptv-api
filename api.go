package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "http://timetableapi.ptv.vic.gov.au"
)

// Client manages the communication with the API
// Modeled after https://github.com/digitalocean/godo/blob/master/godo.go
type Client struct {
	// HTTP client
	client *http.Client
	// Base URL for requests
	BaseURL *url.URL
	// Developer ID
	developerID string
	// Security key
	securityKey string
	// Optional callback that is called after request completion
	onRequestCompleted RequestCompletedCallback

	// Services
	HealthCheck HealthCheckService
}

// RequestCompletedCallback defines the type of request callback function
type RequestCompletedCallback func(*http.Request, *http.Response)

// NewClient returns a new PTV API client
func NewClient(developerID, securityKey string) *Client {
	httpClient := http.DefaultClient
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL,
		developerID: developerID, securityKey: securityKey,
	}

	c.HealthCheck = &HealthCheckServiceOp{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if urlStr[0] != '/' {
		urlStr = "/" + urlStr
	}

	rel, err := url.Parse(urlStr)
	rel.Path = strings.TrimSuffix(rel.Path, "/")
	if err != nil {
		return nil, err
	}

	// The devid param has to be added before the signature is calculated
	q := rel.Query()
	q.Set("devid", c.developerID)
	rel.RawQuery = q.Encode()

	// Signature
	sig := c.GenerateSignature(rel.String())
	q.Set("signature", sig)
	rel.RawQuery = q.Encode()

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// req.Header.Add("Content-Type", mediaType)
	// req.Header.Add("Accept", mediaType)
	// req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do actually sends the API request. It calculates the signature using the request url,
// appends it to the request and returns the API response. The API response is JSON decoded
// and stored as "v" or returns an error. If "v" implements the io.Writer interface, the raw
// response is written to it, instead of decoding
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	// Fire callback if required
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err := json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return resp, err
}

// GenerateSignature returns the HMAC-SHA1 of the complete request string
func (c *Client) GenerateSignature(urlStr string) string {
	key := []byte(c.securityKey)
	sig := hmac.New(sha1.New, key)
	// Only the request path should be used for generating the signature
	sig.Write([]byte(urlStr))

	signature := hex.EncodeToString(sig.Sum(nil))
	return signature
}
