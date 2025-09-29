import React from "react";
import "./Message.css";

const Message = ({ message, isOwn, showAvatar }) => {
  const formatTime = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  const getInitials = (email) => {
    return email.split("@")[0].substring(0, 2).toUpperCase();
  };

  return (
    <div
      className={`message ${isOwn ? "own" : "other"} ${
        showAvatar ? "with-avatar" : "no-avatar"
      }`}
    >
      {showAvatar && !isOwn && (
        <div className="message-avatar">{getInitials(message.user_email)}</div>
      )}

      <div className="message-content">
        {showAvatar && !isOwn && (
          <div className="message-sender">{message.user_email}</div>
        )}

        <div className="message-bubble">
          <div className="message-text">{message.content}</div>
          <div className="message-time">{formatTime(message.created_at)}</div>
        </div>
      </div>
    </div>
  );
};

export default Message;
