import React, { useEffect, useState } from "react";
import ChannelList from "./ChannelList";
import "./ChatApp.css";
import ChatInterface from "./ChatInterface";

const ChatApp = ({ user, onLogout }) => {
  const [selectedChannel, setSelectedChannel] = useState(null);
  const [channels, setChannels] = useState([]);
  const [loading, setLoading] = useState(true);
  const [connectedChannel, setConnectedChannel] = useState(null);

  useEffect(() => {
    fetchChannels();
  }, []);

  const fetchChannels = async () => {
    try {
      const response = await fetch("/api/v1/channels");
      const data = await response.json();

      if (data.success) {
        setChannels(data.channels || []);
        if (data.channels && data.channels.length > 0 && !selectedChannel) {
          setSelectedChannel(data.channels[0]);
        }
      }
    } catch (error) {
      console.error("Error fetching channels:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleChannelSelect = (channel) => {
    setSelectedChannel(channel);
  };

  const handleChannelCreated = (newChannel) => {
    setChannels((prev) => {
      // Check if channel already exists to avoid duplicates
      const exists = prev.some((channel) => channel.id === newChannel.id);
      if (exists) {
        return prev;
      }
      return [...prev, newChannel];
    });
    setSelectedChannel(newChannel);
  };

  if (loading) {
    return (
      <div className="chat-app loading">
        <div className="spinner"></div>
        <p>Loading channels...</p>
      </div>
    );
  }

  return (
    <div className="chat-app">
      <div className="chat-header">
        <div className="chat-header-left">
          <h1>Jobsity Chat</h1>
          <span className="user-info">Welcome, {user.email}</span>
        </div>
        <div className="chat-header-right">
          <button onClick={onLogout} className="btn btn-secondary">
            Logout
          </button>
        </div>
      </div>

      <div className="chat-content">
        <ChannelList
          channels={channels}
          selectedChannel={selectedChannel}
          onChannelSelect={handleChannelSelect}
          onChannelCreated={handleChannelCreated}
          user={user}
          connectedChannel={connectedChannel}
        />

        <ChatInterface
          channel={selectedChannel}
          user={user}
          onConnectionChange={setConnectedChannel}
        />
      </div>
    </div>
  );
};

export default ChatApp;
