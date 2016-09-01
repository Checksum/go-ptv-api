package ptv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStopsNearMe(t *testing.T) {
	lat := float64(-37.8572239)
	lon := float64(144.9995594)
	stops, resp, err := client.StopsNearMe.Get(lat, lon)

	if err != nil {
		t.Errorf("%+v", err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, stops)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, 200)

	if len(stops) == 0 {
		t.Errorf("Number of stops nearby should be > 0")
	}

	assert.EqualValues(t, stops[0].Result.RouteType, NightBus)
}
