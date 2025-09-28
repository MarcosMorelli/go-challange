package unit

import (
	"context"
	"testing"

	"jobsity-backend/internal/service"
	"jobsity-backend/pkg/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ServiceTestSuite contains the test suite for service unit tests
type ServiceTestSuite struct {
	suite.Suite
	userService service.UserService
	mockRepo    *MockUserRepository
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockUserRepository)
	suite.userService = service.NewUserService(suite.mockRepo)
}

// TestLoginSuccess tests successful login
func (suite *ServiceTestSuite) TestLoginSuccess() {
	// Arrange
	email := "test@example.com"
	password := "testpass"
	user := &domain.User{
		ID:       "123",
		Email:    email,
		Password: password,
	}

	suite.mockRepo.On("FindByEmail", mock.Anything, email).Return(user, nil)

	// Act
	req := &domain.LoginRequest{
		Email:    email,
		Password: password,
	}
	response, err := suite.userService.Login(context.Background(), req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Login successful", response.Message)
	assert.NotEmpty(suite.T(), response.Token)
	assert.Equal(suite.T(), user, response.User)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestLoginInvalidPassword tests login with invalid password
func (suite *ServiceTestSuite) TestLoginInvalidPassword() {
	// Arrange
	email := "test@example.com"
	password := "testpass"
	wrongPassword := "wrongpass"
	user := &domain.User{
		ID:       "123",
		Email:    email,
		Password: password,
	}

	suite.mockRepo.On("FindByEmail", mock.Anything, email).Return(user, nil)

	// Act
	req := &domain.LoginRequest{
		Email:    email,
		Password: wrongPassword,
	}
	response, err := suite.userService.Login(context.Background(), req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "Invalid email or password", response.Message)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestLoginUserNotFound tests login when user is not found
func (suite *ServiceTestSuite) TestLoginUserNotFound() {
	// Arrange
	email := "nonexistent@example.com"
	password := "testpass"

	suite.mockRepo.On("FindByEmail", mock.Anything, email).Return(nil, mongo.ErrNoDocuments)

	// Act
	req := &domain.LoginRequest{
		Email:    email,
		Password: password,
	}
	response, err := suite.userService.Login(context.Background(), req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "Invalid email or password", response.Message)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestLoginEmptyCredentials tests login with empty credentials
func (suite *ServiceTestSuite) TestLoginEmptyCredentials() {
	// Act
	req := &domain.LoginRequest{
		Email:    "",
		Password: "",
	}
	response, err := suite.userService.Login(context.Background(), req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "Email and password are required", response.Message)
}

// TestCreateUserSuccess tests successful user creation
func (suite *ServiceTestSuite) TestCreateUserSuccess() {
	// Arrange
	req := &domain.CreateUserRequest{
		Email:    "new@example.com",
		Password: "newpass",
	}

	suite.mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)
	suite.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	// Act
	user, err := suite.userService.CreateUser(context.Background(), req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), req.Email, user.Email)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateUserAlreadyExists tests user creation when user already exists
func (suite *ServiceTestSuite) TestCreateUserAlreadyExists() {
	// Arrange
	req := &domain.CreateUserRequest{
		Email:    "existing@example.com",
		Password: "newpass",
	}

	existingUser := &domain.User{
		ID:       "123",
		Email:    req.Email,
		Password: "oldpass",
	}

	suite.mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	// Act
	user, err := suite.userService.CreateUser(context.Background(), req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user with this email already exists", err.Error())

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateUserInvalidData tests user creation with invalid data
func (suite *ServiceTestSuite) TestCreateUserInvalidData() {
	// Arrange
	req := &domain.CreateUserRequest{
		Email:    "", // Invalid: empty email
		Password: "newpass",
	}

	// Act
	user, err := suite.userService.CreateUser(context.Background(), req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "email and password are required", err.Error())
}

// TestGetUserByEmailSuccess tests successful user retrieval by email
func (suite *ServiceTestSuite) TestGetUserByEmailSuccess() {
	// Arrange
	email := "test@example.com"
	user := &domain.User{
		ID:       "123",
		Email:    email,
		Password: "testpass",
	}

	suite.mockRepo.On("FindByEmail", mock.Anything, email).Return(user, nil)

	// Act
	result, err := suite.userService.GetUserByEmail(context.Background(), email)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user, result)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetUserByIDSuccess tests successful user retrieval by ID
func (suite *ServiceTestSuite) TestGetUserByIDSuccess() {
	// Arrange
	userID := "123"
	user := &domain.User{
		ID:       userID,
		Email:    "test@example.com",
		Password: "testpass",
	}

	suite.mockRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

	// Act
	result, err := suite.userService.GetUserByID(context.Background(), userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user, result)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestSuite runs the test suite
func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
