package ptv

import (
	"net/http"
	"net/url"
	"time"
)

// HealthCheck response
type HealthCheck struct {
	SecurityTokenOK bool `json:"securityTokenOK,omitempty"`
	ClientClockOK   bool `json:"clientClockOK,omitempty"`
	MemcacheOK      bool `json:"memcacheOK,omitempty"`
	DatabaseOK      bool `json:"databaseOK,omitempty"`
}

// HealthCheckService interfaces with the /healthcheck endpoint
type HealthCheckService interface {
	Get() (*HealthCheck, *http.Response, error)
}

// HealthCheckServiceOp handles communication
type HealthCheckServiceOp struct {
	client *Client
}

var _ HealthCheckService = &HealthCheckServiceOp{}

// Get healthcheck info
func (s *HealthCheckServiceOp) Get() (*HealthCheck, *http.Response, error) {
	path, err := url.Parse("/v2/healthcheck")
	if err != nil {
		return nil, nil, err
	}

	q := url.Values{}
	q.Set("timestamp", time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	path.RawQuery = q.Encode()

	req, err := s.client.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	healthInfo := new(HealthCheck)
	resp, err := s.client.Do(req, healthInfo)
	if err != nil {
		return nil, resp, err
	}

	return healthInfo, resp, err
}
