import React, { useEffect, useRef, useState } from "react";
import "./ChatInterface.css";
import MessageInput from "./MessageInput";
import MessageList from "./MessageList";

const ChatInterface = ({ channel, user, onConnectionChange }) => {
  const [messages, setMessages] = useState([]);
  const [connected, setConnected] = useState(false);
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef(null);
  const wsRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Cleanup on component unmount
  useEffect(() => {
    return () => {
      if (wsRef.current) {
        console.log("Component unmounting, closing WebSocket connection");
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, []);

  useEffect(() => {
    if (!channel) {
      // If no channel selected, close any existing connection
      if (wsRef.current) {
        console.log("No channel selected, closing WebSocket connection");
        wsRef.current.close();
        wsRef.current = null;
        setConnected(false);
      }
      return;
    }

    // Fetch existing messages
    fetchMessages();

    // Connect to WebSocket
    connectWebSocket();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
        setConnected(false);
      }
    };
  }, [channel]);

  const fetchMessages = async () => {
    if (!channel) return;

    setLoading(true);
    try {
      // Fetch messages with limit and ensure they're sorted by created_at desc
      const response = await fetch(
        `/api/v1/channels/${channel.id}/messages?limit=50`
      );
      const data = await response.json();

      if (data.success) {
        // Messages should already be sorted by created_at desc from backend
        // But let's ensure they're in the correct order for display
        const messages = data.messages || [];
        setMessages(messages);
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
    } finally {
      setLoading(false);
    }
  };

  const connectWebSocket = () => {
    if (!channel) {
      console.log("No channel selected, skipping WebSocket connection");
      return;
    }

    // Close existing connection if any
    if (wsRef.current) {
      console.log(
        "Closing existing WebSocket connection before creating new one"
      );
      wsRef.current.close();
      wsRef.current = null;
      setConnected(false);
    }

    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//localhost:3000/api/v1/ws?channel_id=${
      channel.id
    }&user_email=${encodeURIComponent(user.email)}`;

    const websocket = new WebSocket(wsUrl);

    websocket.onopen = () => {
      setConnected(true);
      wsRef.current = websocket;
      onConnectionChange?.(channel);
    };

    websocket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };

    websocket.onclose = (event) => {
      setConnected(false);
      wsRef.current = null;
      onConnectionChange?.(null);
    };

    websocket.onerror = (error) => {
      console.error("WebSocket error:", error);
      console.error("Failed to connect to:", wsUrl);
      setConnected(false);
    };
  };

  const handleWebSocketMessage = (message) => {
    switch (message.type) {
      case "new_message":
        // Add new message and sort by created_at desc to maintain order
        setMessages((prev) => {
          const updated = [...prev, message.data];
          return updated.sort(
            (a, b) => new Date(b.created_at) - new Date(a.created_at)
          );
        });
        break;
      case "message_updated":
        setMessages((prev) =>
          prev.map((msg) => (msg.id === message.data.id ? message.data : msg))
        );
        break;
      case "message_deleted":
        setMessages((prev) => prev.filter((msg) => msg.id !== message.data.id));
        break;
    }
  };

  const handleSendMessage = async (content) => {
    if (!channel || !content.trim()) return;

    try {
      const response = await fetch("/api/v1/messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "User-Email": user.email,
        },
        body: JSON.stringify({
          channel_id: channel.id,
          content: content.trim(),
        }),
      });

      const data = await response.json();

      if (!data.success) {
        console.error("Failed to send message:", data.message);
      }
    } catch (error) {
      console.error("Error sending message:", error);
    }
  };

  if (!channel) {
    return (
      <div className="chat-interface no-channel">
        <div className="no-channel-content">
          <h2>Select a channel to start chatting</h2>
          <p>Choose a channel from the sidebar to begin your conversation.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="chat-interface">
      <div className="chat-header">
        <div className="channel-info">
          <h2>{channel.name}</h2>
          <p>{channel.description}</p>
        </div>
        <div className="connection-status">
          <div
            className={`status-indicator ${
              connected ? "connected" : "disconnected"
            }`}
          >
            {connected ? `Connected to ${channel.name}` : "Disconnected"}
          </div>
          {!connected && (
            <button
              onClick={connectWebSocket}
              className="btn btn-secondary"
              style={{
                marginLeft: "10px",
                padding: "4px 8px",
                fontSize: "12px",
              }}
            >
              Retry Connection
            </button>
          )}
        </div>
      </div>

      <div className="chat-messages">
        {loading ? (
          <div className="loading-messages">
            <div className="spinner"></div>
            <p>Loading messages...</p>
          </div>
        ) : (
          <MessageList messages={messages} currentUser={user.email} />
        )}
        <div ref={messagesEndRef} />
      </div>

      <MessageInput onSendMessage={handleSendMessage} disabled={!connected} />
    </div>
  );
};

export default ChatInterface;
