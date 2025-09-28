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

// MessageTestSuite contains the test suite for message integration tests
type MessageTestSuite struct {
	suite.Suite
	client    *http.Client
	baseURL   string
	channelID string
}

// SetupSuite sets up the test suite
func (suite *MessageTestSuite) SetupSuite() {
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
func (suite *MessageTestSuite) waitForServer() {
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

// createTestChannel creates a test channel for message tests
func (suite *MessageTestSuite) createTestChannel() {
	channelData := domain.CreateChannelRequest{
		Name:        "message-test-channel-" + time.Now().Format("20060102150405"),
		Description: "Test channel for messages",
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
		var response domain.ChannelResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		suite.channelID = response.Channel.ID
	}
}

// TestCreateMessage tests message creation
func (suite *MessageTestSuite) TestCreateMessage() {
	messageData := domain.CreateMessageRequest{
		ChannelID: suite.channelID,
		Content:   "Hello, this is a test message!",
	}

	jsonData, _ := json.Marshal(messageData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/messages", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Email", "test@example.com")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 201, resp.StatusCode)

	var response domain.MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Message created successfully", response.Message)
	assert.NotNil(suite.T(), response.Data)
}

// TestCreateMessageInvalidData tests message creation with invalid data
func (suite *MessageTestSuite) TestCreateMessageInvalidData() {
	messageData := domain.CreateMessageRequest{
		ChannelID: "", // Invalid: empty channel ID
		Content:   "Hello, this is a test message!",
	}

	jsonData, _ := json.Marshal(messageData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/messages", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Email", "test@example.com")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 400, resp.StatusCode)
}

// TestGetMessagesByChannel tests getting messages for a channel
func (suite *MessageTestSuite) TestGetMessagesByChannel() {
	req, err := http.NewRequest("GET", suite.baseURL+"/api/v1/channels/"+suite.channelID+"/messages", nil)
	assert.NoError(suite.T(), err)

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response domain.MessagesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Messages retrieved successfully", response.Message)
}

// TestMessageSuite runs the test suite
func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}
