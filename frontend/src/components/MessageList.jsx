import React from "react";
import Message from "./Message";
import "./MessageList.css";

const MessageList = ({ messages, currentUser }) => {
  if (messages.length === 0) {
    return (
      <div className="message-list empty">
        <div className="empty-messages">
          <h3>No messages yet</h3>
          <p>Be the first to send a message in this channel!</p>
        </div>
      </div>
    );
  }

  console.log(messages);

  return (
    <div className="message-list">
      {messages.map((message, index) => {
        const prevMessage = index > 0 ? messages[index - 1] : null;
        const showAvatar =
          !prevMessage || prevMessage.user_email !== message.user_email;

        return (
          <Message
            key={message.id}
            message={message}
            isOwn={message.user_email === currentUser}
            showAvatar={showAvatar}
          />
        );
      })}
    </div>
  );
};

export default MessageList;
