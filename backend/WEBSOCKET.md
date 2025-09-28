# WebSocket Real-Time Chat Implementation

This document describes the WebSocket implementation for real-time chat functionality in the Jobsity backend.

## 🚀 Features

- **Real-time message broadcasting** to all clients in a channel
- **Channel-based messaging** with automatic client management
- **Message lifecycle events** (create, update, delete) with real-time notifications
- **Connection management** with automatic cleanup
- **WebSocket statistics** for monitoring
- **Authentication integration** with existing user system

## 🏗️ Architecture

### Components

1. **WebSocket Hub** (`internal/websocket/hub.go`)

   - Manages all WebSocket connections
   - Handles client registration/unregistration
   - Broadcasts messages to specific channels or all clients

2. **WebSocket Client** (`internal/websocket/client.go`)

   - Represents individual WebSocket connections
   - Handles message reading/writing
   - Manages channel joining/leaving

3. **WebSocket Handler** (`internal/websocket/handler.go`)

   - HTTP-to-WebSocket upgrade handling
   - Message broadcasting logic
   - Statistics endpoint

4. **WebSocket Message Service** (`internal/service/websocket_message_service.go`)
   - Integrates WebSocket broadcasting with message operations
   - Automatically broadcasts message events

## 📡 WebSocket API

### Connection

**Endpoint:** `GET /api/v1/ws`

**Headers:**

- `User-Email`: User email for authentication

**Query Parameters:**

- `channel_id` (optional): Channel to join immediately

**Example:**

```javascript
const ws = new WebSocket(
  "ws://localhost:3000/api/v1/ws?channel_id=channel123",
  {
    headers: {
      "User-Email": "user@example.com",
    },
  }
);
```

### Message Types

#### Client → Server Messages

**Join Channel:**

```json
{
  "type": "join_channel",
  "channel_id": "channel123"
}
```

**Leave Channel:**

```json
{
  "type": "leave_channel",
  "channel_id": "channel123"
}
```

**Ping:**

```json
{
  "type": "ping"
}
```

#### Server → Client Messages

**New Message:**

```json
{
  "type": "new_message",
  "channel_id": "channel123",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "id": "msg123",
    "channel_id": "channel123",
    "user_email": "user@example.com",
    "content": "Hello world!",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

**Message Updated:**

```json
{
  "type": "message_updated",
  "channel_id": "channel123",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "id": "msg123",
    "channel_id": "channel123",
    "user_email": "user@example.com",
    "content": "Updated message",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

**Message Deleted:**

```json
{
  "type": "message_deleted",
  "channel_id": "channel123",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "id": "msg123",
    "channel_id": "channel123",
    "user_email": "user@example.com"
  }
}
```

**Channel Joined:**

```json
{
  "type": "channel_joined",
  "channel_id": "channel123",
  "user_email": "user@example.com",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Channel Left:**

```json
{
  "type": "channel_left",
  "channel_id": "channel123",
  "user_email": "user@example.com",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Pong:**

```json
{
  "type": "pong",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Statistics Endpoint

**Endpoint:** `GET /api/v1/ws/stats`

**Response:**

```json
{
  "success": true,
  "data": {
    "total_clients": 5,
    "channels": {
      "channel123": 3,
      "channel456": 2
    }
  }
}
```

## 🧪 Testing

### HTML Test Client

A complete HTML test client is available at `http://localhost:3000/` when the server is running.

**Features:**

- WebSocket connection management
- Channel joining/leaving
- Real-time message sending/receiving
- Connection statistics
- Ping/pong testing

### Integration Tests

Run WebSocket integration tests:

```bash
# Run all integration tests including WebSocket
make test-integration-docker

# Or run WebSocket tests specifically
go test -v ./tests/integration/... -run "TestWebSocketSuite"
```

## 🔧 Usage Examples

### JavaScript Client

```javascript
// Connect to WebSocket
const ws = new WebSocket("ws://localhost:3000/api/v1/ws?channel_id=general", {
  headers: {
    "User-Email": "user@example.com",
  },
});

// Handle messages
ws.onmessage = function (event) {
  const message = JSON.parse(event.data);

  switch (message.type) {
    case "new_message":
      console.log("New message:", message.data.content);
      break;
    case "message_updated":
      console.log("Message updated:", message.data.content);
      break;
    case "message_deleted":
      console.log("Message deleted:", message.data.id);
      break;
  }
};

// Join a channel
function joinChannel(channelId) {
  ws.send(
    JSON.stringify({
      type: "join_channel",
      channel_id: channelId,
    })
  );
}

// Leave current channel
function leaveChannel() {
  ws.send(
    JSON.stringify({
      type: "leave_channel",
      channel_id: currentChannelId,
    })
  );
}

// Send a message via REST API (triggers WebSocket broadcast)
function sendMessage(channelId, content) {
  fetch("/api/v1/messages", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "User-Email": "user@example.com",
    },
    body: JSON.stringify({
      channel_id: channelId,
      content: content,
    }),
  });
}
```

### Go Client

```go
package main

import (
    "net/url"
    "github.com/gorilla/websocket"
)

func main() {
    // Connect to WebSocket
    u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/api/v1/ws"}
    headers := http.Header{}
    headers.Set("User-Email", "user@example.com")

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
    if err != nil {
        log.Fatal("dial:", err)
    }
    defer conn.Close()

    // Join channel
    joinMessage := map[string]string{
        "type": "join_channel",
        "channel_id": "general",
    }
    conn.WriteJSON(joinMessage)

    // Read messages
    for {
        var message map[string]interface{}
        err := conn.ReadJSON(&message)
        if err != nil {
            log.Println("read:", err)
            return
        }
        log.Printf("Received: %v", message)
    }
}
```

## 🚀 Real-Time Features

### Automatic Broadcasting

When messages are created, updated, or deleted via the REST API, they are automatically broadcast to all connected clients in the relevant channel.

### Channel Management

- Clients can join/leave channels dynamically
- Messages are only broadcast to clients in the relevant channel
- Connection cleanup when clients disconnect

### Connection Health

- Automatic ping/pong for connection health monitoring
- Graceful connection cleanup
- Statistics for monitoring connection counts

## 🔒 Security

- **Authentication Required**: All WebSocket connections require valid user email
- **Channel Isolation**: Messages are only broadcast to clients in the relevant channel
- **Input Validation**: All incoming messages are validated
- **Rate Limiting**: Built-in connection limits and message size limits

## 📊 Monitoring

### Statistics

Monitor WebSocket connections with the stats endpoint:

```bash
curl http://localhost:3000/api/v1/ws/stats
```

### Logs

WebSocket events are logged for debugging:

- Client connections/disconnections
- Channel joins/leaves
- Message broadcasts
- Connection errors

## 🎯 Next Steps

To enhance the WebSocket implementation:

1. **JWT Authentication**: Replace header-based auth with JWT tokens
2. **Message History**: Implement message persistence and history
3. **User Presence**: Add online/offline status
4. **File Sharing**: Support file uploads with WebSocket notifications
5. **Message Reactions**: Add emoji reactions to messages
6. **Private Messages**: Support direct messaging between users
7. **Message Encryption**: Add end-to-end encryption for sensitive messages

The WebSocket implementation provides a solid foundation for real-time chat functionality! 🎉
