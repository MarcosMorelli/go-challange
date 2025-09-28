package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"jobsity-backend/pkg/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ChannelTestSuite contains the test suite for channel integration tests
type ChannelTestSuite struct {
	suite.Suite
	client  *http.Client
	baseURL string
}

// SetupSuite sets up the test suite
func (suite *ChannelTestSuite) SetupSuite() {
	suite.client = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.baseURL = "http://localhost:3000"

	// Wait for server to be ready
	suite.waitForServer()
}

// waitForServer waits for the server to be ready
func (suite *ChannelTestSuite) waitForServer() {
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

// TestCreateChannel tests channel creation
func (suite *ChannelTestSuite) TestCreateChannel() {
	channelData := domain.CreateChannelRequest{
		Name:        "test-channel-" + time.Now().Format("20060102150405"),
		Description: "Test channel description",
	}

	jsonData, _ := json.Marshal(channelData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/channels", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Email", "test@example.com") // Mock authentication

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 201, resp.StatusCode)

	var response domain.ChannelResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Channel created successfully", response.Message)
	assert.NotNil(suite.T(), response.Channel)
}

// TestCreateChannelInvalidData tests channel creation with invalid data
func (suite *ChannelTestSuite) TestCreateChannelInvalidData() {
	channelData := domain.CreateChannelRequest{
		Name:        "", // Invalid: empty name
		Description: "Test channel description",
	}

	jsonData, _ := json.Marshal(channelData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/channels", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Email", "test@example.com")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 400, resp.StatusCode)
}

// TestGetAllChannels tests getting all channels
func (suite *ChannelTestSuite) TestGetAllChannels() {
	req, err := http.NewRequest("GET", suite.baseURL+"/api/v1/channels", nil)
	assert.NoError(suite.T(), err)

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response domain.ChannelsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Channels retrieved successfully", response.Message)
}

// TestChannelSuite runs the test suite
func TestChannelSuite(t *testing.T) {
	suite.Run(t, new(ChannelTestSuite))
}
