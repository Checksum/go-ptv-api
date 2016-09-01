package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	healthCheck, resp, err := client.HealthCheck.Get()
	assert.Nil(t, err)
	assert.NotNil(t, healthCheck)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, 200)
	// Health check should pass
	assert.Equal(t, healthCheck.SecurityTokenOK, true)
	assert.Equal(t, healthCheck.ClientClockOK, true)
}
