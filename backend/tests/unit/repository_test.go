package unit

import (
	"context"
	"testing"
	"time"

	"jobsity-backend/internal/repository"
	"jobsity-backend/pkg/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// RepositoryTestSuite contains the test suite for repository tests
type RepositoryTestSuite struct {
	suite.Suite
	mt *mtest.T
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.mt = mtest.New(suite.T(), mtest.NewOptions().ClientType(mtest.Mock))
}

func (suite *RepositoryTestSuite) TearDownTest() {
	// mtest.T doesn't have a Close method, cleanup is handled automatically
}

// TestMongoUserRepository_FindByEmail tests finding user by email
func (suite *RepositoryTestSuite) TestMongoUserRepository_FindByEmail() {
	suite.mt.Run("success", func(mt *mtest.T) {
		// Arrange
		email := "test@example.com"
		expectedUser := domain.User{
			ID:        "507f1f77bcf86cd799439011",
			Email:     email,
			Password:  "testpass",
			CreatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedUser.ID},
			{Key: "email", Value: expectedUser.Email},
			{Key: "password", Value: expectedUser.Password},
			{Key: "created_at", Value: expectedUser.CreatedAt},
		}))

		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		user, err := repo.FindByEmail(context.Background(), email)

		// Assert
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), expectedUser.Email, user.Email)
	})

	suite.mt.Run("user not found", func(mt *mtest.T) {
		// Arrange
		email := "nonexistent@example.com"
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		user, err := repo.FindByEmail(context.Background(), email)

		// Assert
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
		assert.Equal(suite.T(), mongo.ErrNoDocuments, err)
	})
}

// TestMongoUserRepository_Create tests user creation
func (suite *RepositoryTestSuite) TestMongoUserRepository_Create() {
	suite.mt.Run("success", func(mt *mtest.T) {
		// Arrange
		user := &domain.User{
			Email:    "new@example.com",
			Password: "newpass",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		err := repo.Create(context.Background(), user)

		// Assert
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), user.ID)
		assert.False(suite.T(), user.CreatedAt.IsZero())
	})
}

// TestMongoUserRepository_FindByID tests finding user by ID
func (suite *RepositoryTestSuite) TestMongoUserRepository_FindByID() {
	suite.mt.Run("success", func(mt *mtest.T) {
		// Arrange
		userID := "507f1f77bcf86cd799439011"
		expectedUser := domain.User{
			ID:        userID,
			Email:     "test@example.com",
			Password:  "testpass",
			CreatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userID},
			{Key: "email", Value: expectedUser.Email},
			{Key: "password", Value: expectedUser.Password},
			{Key: "created_at", Value: expectedUser.CreatedAt},
		}))

		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		user, err := repo.FindByID(context.Background(), userID)

		// Assert
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), expectedUser.ID, user.ID)
		assert.Equal(suite.T(), expectedUser.Email, user.Email)
	})

	suite.mt.Run("invalid ID", func(mt *mtest.T) {
		// Arrange
		invalidID := "invalid"
		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		user, err := repo.FindByID(context.Background(), invalidID)

		// Assert
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
	})
}

// TestMongoUserRepository_Update tests user update
func (suite *RepositoryTestSuite) TestMongoUserRepository_Update() {
	suite.mt.Run("success", func(mt *mtest.T) {
		// Arrange
		user := &domain.User{
			ID:       "507f1f77bcf86cd799439011",
			Email:    "updated@example.com",
			Password: "updatedpass",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		err := repo.Update(context.Background(), user)

		// Assert
		assert.NoError(suite.T(), err)
	})

	suite.mt.Run("invalid ID", func(mt *mtest.T) {
		// Arrange
		user := &domain.User{
			ID:       "invalid",
			Email:    "updated@example.com",
			Password: "updatedpass",
		}

		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		err := repo.Update(context.Background(), user)

		// Assert
		assert.Error(suite.T(), err)
	})
}

// TestMongoUserRepository_Delete tests user deletion
func (suite *RepositoryTestSuite) TestMongoUserRepository_Delete() {
	suite.mt.Run("success", func(mt *mtest.T) {
		// Arrange
		userID := "507f1f77bcf86cd799439011"
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		err := repo.Delete(context.Background(), userID)

		// Assert
		assert.NoError(suite.T(), err)
	})

	suite.mt.Run("invalid ID", func(mt *mtest.T) {
		// Arrange
		invalidID := "invalid"
		repo := repository.NewMongoUserRepository(mt.Coll)

		// Act
		err := repo.Delete(context.Background(), invalidID)

		// Assert
		assert.Error(suite.T(), err)
	})
}

// TestSuite runs the test suite
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
