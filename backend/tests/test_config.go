package tests

import (
	"os"
	"testing"
)

// SetupTestEnvironment sets up the test environment
func SetupTestEnvironment(t *testing.T) {
	// Set test environment variables
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DATABASE", "jobsity_test")
	os.Setenv("SERVER_PORT", "3001") // Use different port for tests
}

// CleanupTestEnvironment cleans up the test environment
func CleanupTestEnvironment(t *testing.T) {
	// Clean up environment variables if needed
	os.Unsetenv("MONGODB_URI")
	os.Unsetenv("MONGODB_DATABASE")
	os.Unsetenv("SERVER_PORT")
}
