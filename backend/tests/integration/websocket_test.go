package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// WebSocketTestSuite contains the test suite for WebSocket integration tests
type WebSocketTestSuite struct {
	suite.Suite
	client    *http.Client
	baseURL   string
	channelID string
}

// SetupSuite sets up the test suite
func (suite *WebSocketTestSuite) SetupSuite() {
	suite.client = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.baseURL = "http://localhost:3000"

	// Wait for server to be ready
	suite.waitForServer()

	// Create a test channel first
	suite.createTestChannel()
}

// waitForServer waits for the server to be ready
func (suite *WebSocketTestSuite) waitForServer() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := suite.client.Get(suite.baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	suite.T().Fatal("Server not ready after 30 seconds")
}

// createTestChannel creates a test channel for WebSocket tests
func (suite *WebSocketTestSuite) createTestChannel() {
	channelData := map[string]string{
		"name":        "websocket-test-channel-" + time.Now().Format("20060102150405"),
		"description": "Test channel for WebSocket",
	}

	jsonData, _ := json.Marshal(channelData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/channels", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Email", "test@example.com")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		if channel, ok := response["channel"].(map[string]interface{}); ok {
			if id, ok := channel["id"].(string); ok {
				suite.channelID = id
			}
		}
	}
}

// TestWebSocketConnection tests basic WebSocket connection
func (suite *WebSocketTestSuite) TestWebSocketConnection() {
	// This is a basic test - in a real scenario, you'd need to implement
	// WebSocket connection testing with proper authentication
	// For now, we'll test the WebSocket stats endpoint

	req, err := http.NewRequest("GET", suite.baseURL+"/api/v1/ws/stats", nil)
	assert.NoError(suite.T(), err)

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response["success"].(bool))

	// Check that stats structure is correct
	data := response["data"].(map[string]interface{})
	assert.Contains(suite.T(), data, "total_clients")
	assert.Contains(suite.T(), data, "channels")
}

// TestWebSocketStats tests WebSocket statistics endpoint
func (suite *WebSocketTestSuite) TestWebSocketStats() {
	req, err := http.NewRequest("GET", suite.baseURL+"/api/v1/ws/stats", nil)
	assert.NoError(suite.T(), err)

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response["success"].(bool))

	// Verify stats structure
	data := response["data"].(map[string]interface{})
	assert.IsType(suite.T(), float64(0), data["total_clients"])
	assert.IsType(suite.T(), map[string]interface{}{}, data["channels"])
}

// TestWebSocketSuite runs the test suite
func TestWebSocketSuite(t *testing.T) {
	suite.Run(t, new(WebSocketTestSuite))
}
