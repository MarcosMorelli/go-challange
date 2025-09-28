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

// APITestSuite contains the test suite for API integration tests
type APITestSuite struct {
	suite.Suite
	client  *http.Client
	baseURL string
}

// SetupSuite sets up the test suite
func (suite *APITestSuite) SetupSuite() {
	suite.client = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.baseURL = "http://localhost:3000"

	// Wait for server to be ready
	suite.waitForServer()
}

// waitForServer waits for the server to be ready
func (suite *APITestSuite) waitForServer() {
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

// TestHealthCheck tests the health endpoint
func (suite *APITestSuite) TestHealthCheck() {
	resp, err := suite.client.Get(suite.baseURL + "/health")
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ok", response["status"])
}

// TestCreateUser tests user creation
func (suite *APITestSuite) TestCreateUser() {
	userData := domain.CreateUserRequest{
		Email:    "testuser_" + time.Now().Format("20060102150405") + "@example.com",
		Password: "testpass",
	}

	jsonData, _ := json.Marshal(userData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/users", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response["success"].(bool))
	assert.Equal(suite.T(), "User created successfully", response["message"])
}

// TestCreateUserInvalidData tests user creation with invalid data
func (suite *APITestSuite) TestCreateUserInvalidData() {
	userData := domain.CreateUserRequest{
		Email:    "", // Invalid: empty email
		Password: "testpass",
	}

	jsonData, _ := json.Marshal(userData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/users", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 400, resp.StatusCode)
}

// TestLogin tests user login
func (suite *APITestSuite) TestLogin() {
	// First create a user
	userData := domain.CreateUserRequest{
		Email:    "logintest_" + time.Now().Format("20060102150405") + "@example.com",
		Password: "loginpass",
	}

	jsonData, _ := json.Marshal(userData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/users", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Now test login
	loginData := domain.LoginRequest{
		Email:    userData.Email,
		Password: "loginpass",
	}

	jsonData, _ = json.Marshal(loginData)
	req, err = http.NewRequest("POST", suite.baseURL+"/api/v1/login", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response domain.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Login successful", response.Message)
}

// TestLoginInvalidCredentials tests login with invalid credentials
func (suite *APITestSuite) TestLoginInvalidCredentials() {
	loginData := domain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpass",
	}

	jsonData, _ := json.Marshal(loginData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/login", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 401, resp.StatusCode)

	var response domain.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "Invalid email or password", response.Message)
}

// TestGetUser tests getting user by email
func (suite *APITestSuite) TestGetUser() {
	// First create a user
	userData := domain.CreateUserRequest{
		Email:    "getusertest_" + time.Now().Format("20060102150405") + "@example.com",
		Password: "getuserpass",
	}

	jsonData, _ := json.Marshal(userData)
	req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/users", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Now test getting the user
	req, err = http.NewRequest("GET", suite.baseURL+"/api/v1/users/"+userData.Email, nil)
	assert.NoError(suite.T(), err)

	resp, err = suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response["success"].(bool))
}

// TestGetUserNotFound tests getting a non-existent user
func (suite *APITestSuite) TestGetUserNotFound() {
	req, err := http.NewRequest("GET", suite.baseURL+"/api/v1/users/nonexistent@example.com", nil)
	assert.NoError(suite.T(), err)

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), 404, resp.StatusCode)
}

// TestSuite runs the test suite
func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
