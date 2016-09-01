package ptv

import (
	"net/http"
	"strconv"
)

// RouteType represents the type of route. Ex: train, tram, bus
type RouteType int

// Types of stops
const (
	Train RouteType = iota
	Tram
	Bus
	VLine
	NightBus
)

// Stop represents a single transport stop
type Stop struct {
	Distance      float64   `json:"distance"`
	Suburb        string    `json:"suburb"`
	TransportType string    `json:"transport_type"`
	RouteType     RouteType `json:"route_type"`
	StopID        int       `json:"stop_id"`
	LocationName  string    `json:"location_name"`
	Latitude      float64   `json:"lat"`
	Longitude     float64   `json:"lon"`
}

// NearmeResponse represents the result from the API
type NearmeResponse []struct {
	Result Stop   `json:"result"`
	Type   string `json:"type"`
}

// StopService interfaces with the /nearme endpoint
type StopService interface {
	Get(lat, lon float64) (NearmeResponse, *http.Response, error)
}

// StopServiceOp handles communication
type StopServiceOp struct {
	client *Client
}

var _ StopService = &StopServiceOp{}

// Get stops nearme
func (s *StopServiceOp) Get(lat, lon float64) (NearmeResponse, *http.Response, error) {
	path := "/v2/nearme"
	path += "/latitude/" + FloatToString(lat)
	path += "/longitude/" + FloatToString(lon)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(NearmeResponse)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, err
}

// FloatToString converts a float64 number to string
func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
