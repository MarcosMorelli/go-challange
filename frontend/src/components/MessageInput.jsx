import React, { useState } from "react";
import "./MessageInput.css";

const MessageInput = ({ onSendMessage, disabled }) => {
  const [message, setMessage] = useState("");
  const [isTyping, setIsTyping] = useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (message.trim() && !disabled) {
      onSendMessage(message);
      setMessage("");
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const handleChange = (e) => {
    setMessage(e.target.value);

    // Simple typing indicator
    if (e.target.value.trim() && !isTyping) {
      setIsTyping(true);
    } else if (!e.target.value.trim() && isTyping) {
      setIsTyping(false);
    }
  };

  return (
    <div className="message-input">
      <form onSubmit={handleSubmit} className="message-form">
        <div className="input-container">
          <textarea
            value={message}
            onChange={handleChange}
            onKeyPress={handleKeyPress}
            placeholder={disabled ? "Connecting..." : "Type a message..."}
            className="message-textarea"
            disabled={disabled}
            rows="1"
          />
          <button
            type="submit"
            disabled={!message.trim() || disabled}
            className="send-button"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path
                d="M22 2L11 13M22 2L15 22L11 13M22 2L2 9L11 13"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
        </div>
        {isTyping && (
          <div className="typing-indicator">
            <span>Typing...</span>
          </div>
        )}
      </form>
    </div>
  );
};

export default MessageInput;
