# Jobsity Challenge

## Quick Start

To run the application:

1. Start the services with Docker Compose:

   ```bash
   docker compose up -d --build
   ```

2. Open your browser and navigate to:

   ```
   http://localhost:3001
   ```

3. Login with the following credentials:
   - **Email:** `user1@jobsity.com` or `user2@jobsity.com` or create a new user
   - **Password:** `password`

## Services

The application consists of:

- Frontend (React) - Port 3001
- Backend (Go) - Port 3000
- MongoDB - Port 27017
- RabbitMQ - Port 5672
