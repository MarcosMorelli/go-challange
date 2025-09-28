#!/bin/bash

# Test runner script for Jobsity Backend

set -e

echo "🧪 Starting Jobsity Backend Tests"
echo "================================="

# Check if MongoDB is running
echo "📊 Checking MongoDB connection..."
if ! nc -z localhost 27017; then
    echo "❌ MongoDB is not running. Please start MongoDB first."
    echo "   On macOS: brew services start mongodb-community"
    echo "   Or run: mongod"
    exit 1
fi
echo "✅ MongoDB is running"

# Set test environment variables
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="jobsity_test"
export SERVER_PORT="3001"

echo ""
echo "🔧 Running Unit Tests..."
echo "========================"
go test -v ./tests/unit/...

echo ""
echo "🔧 Running Integration Tests..."
echo "=============================="
go test -v ./tests/integration/...

echo ""
echo "📊 Generating Coverage Report..."
echo "==============================="
go test -v -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "✅ All tests completed successfully!"
echo "📈 Coverage report: coverage.html"
echo "🎉 Ready for deployment!"
