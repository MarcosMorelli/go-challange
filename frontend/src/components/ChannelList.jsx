import axios from "axios";
import React, { useState } from "react";
import "./ChannelList.css";

const ChannelList = ({
  channels,
  selectedChannel,
  onChannelSelect,
  onChannelCreated,
  user,
  connectedChannel,
}) => {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [createForm, setCreateForm] = useState({
    name: "",
    description: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleCreateChannel = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await axios.post(
        "/api/v1/channels",
        {
          name: createForm.name,
          description: createForm.description,
        },
        {
          headers: {
            "User-Email": user.email,
          },
        }
      );

      if (response.data.success) {
        onChannelCreated(response.data.channel);
        setCreateForm({ name: "", description: "" });
        setShowCreateForm(false);
      }
    } catch (err) {
      setError(err.response?.data?.message || "Failed to create channel");
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e) => {
    setCreateForm({
      ...createForm,
      [e.target.name]: e.target.value,
    });
    setError("");
  };

  return (
    <div className="channel-list">
      <div className="channel-list-header">
        <h2>Channels</h2>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          className="btn btn-primary"
        >
          + New Channel
        </button>
      </div>

      {showCreateForm && (
        <div className="create-channel-form">
          <form onSubmit={handleCreateChannel}>
            <div className="form-group">
              <input
                type="text"
                name="name"
                value={createForm.name}
                onChange={handleInputChange}
                className="input"
                placeholder="Channel name"
                required
              />
            </div>
            <div className="form-group">
              <input
                type="text"
                name="description"
                value={createForm.description}
                onChange={handleInputChange}
                className="input"
                placeholder="Channel description"
              />
            </div>
            {error && <div className="error-message">{error}</div>}
            <div className="form-actions">
              <button
                type="submit"
                className="btn btn-primary"
                disabled={loading}
              >
                {loading ? "Creating..." : "Create"}
              </button>
              <button
                type="button"
                onClick={() => setShowCreateForm(false)}
                className="btn btn-secondary"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="channels">
        {channels.length === 0 ? (
          <div className="no-channels">
            <p>No channels available</p>
            <p>Create a new channel to get started!</p>
          </div>
        ) : (
          channels.map((channel) => (
            <div
              key={channel.id}
              className={`channel-item ${
                selectedChannel?.id === channel.id ? "active" : ""
              }`}
              onClick={() => onChannelSelect(channel)}
            >
              <div className="channel-name">
                {channel.name}
                {connectedChannel?.id === channel.id && (
                  <span className="connection-indicator">🟢</span>
                )}
              </div>
              <div className="channel-description">{channel.description}</div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ChannelList;
