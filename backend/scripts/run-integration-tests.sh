#!/bin/bash

# Script to run integration tests against Docker container
# This script starts the Docker services and runs the integration tests

set -e

echo "🚀 Starting integration tests with Docker container..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Navigate to project root
cd "$(dirname "$0")/../.."

echo "📦 Starting Docker services..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if the backend service is healthy
echo "🔍 Checking if backend service is ready..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if curl -f http://localhost:3000/health > /dev/null 2>&1; then
        echo "✅ Backend service is ready!"
        break
    fi
    
    attempt=$((attempt + 1))
    echo "⏳ Waiting for backend service... (attempt $attempt/$max_attempts)"
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "❌ Backend service failed to start within expected time"
    echo "📋 Docker logs:"
    docker-compose logs backend
    exit 1
fi

echo "🧪 Running integration tests..."
cd backend
go test -v ./tests/integration/...

echo "🧹 Cleaning up Docker services..."
cd ..
docker-compose down

echo "✅ Integration tests completed!"
