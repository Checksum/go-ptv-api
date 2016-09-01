package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client *Client

// Initialize the client with the auth info from env
func init() {
	devid := os.Getenv("PTV_DEVID")
	secret := os.Getenv("PTV_SECRET")
	client = NewClient(devid, secret)
}

func TestClient(t *testing.T) {
	assert.NotNil(t, client.developerID)
	assert.NotNil(t, client.securityKey)
}

func TestGenerateSignature(t *testing.T) {
	client := NewClient("2", "7car2d2b-7527-14e1-8975-06cf1059afe0")
	signature := client.GenerateSignature("/v2/healthcheck?devid=2")
	assert.Equal(t, signature, "7a98b58785754b6af5fa51899666e767085b8ef4")
}

func TestNewRequest(t *testing.T) {
	client := NewClient("2", "7car2d2b-7527-14e1-8975-06cf1059afe0")
	req, err := client.NewRequest("GET", "/v2/healthcheck", nil)
	assert.Nil(t, err)
	assert.Equal(t, req.URL.RawQuery, "devid=2&signature=7a98b58785754b6af5fa51899666e767085b8ef4")
	assert.Equal(t, req.URL.String(), client.BaseURL.String()+"/v2/healthcheck?devid=2&signature=7a98b58785754b6af5fa51899666e767085b8ef4")
}
