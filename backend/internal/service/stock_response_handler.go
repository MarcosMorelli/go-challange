package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

type StockResponseHandler struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	broadcastFunc func(channelID string, message []byte)
}

type StockResponseMessage struct {
	ChannelID string `json:"channel_id"`
	UserEmail string `json:"user_email"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func NewStockResponseHandler(conn *amqp.Connection, broadcastFunc func(channelID string, message []byte)) (*StockResponseHandler, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &StockResponseHandler{
		conn:          conn,
		ch:            ch,
		broadcastFunc: broadcastFunc,
	}, nil
}

func (h *StockResponseHandler) Start() error {
	// Start consuming stock responses
	go h.consumeStockResponses()
	return nil
}

func (h *StockResponseHandler) consumeStockResponses() {
	msgs, err := h.ch.Consume(
		"stock_responses", // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		fmt.Printf("Failed to register a consumer for stock responses: %s\n", err)
		return
	}

	for msg := range msgs {
		h.processStockResponse(msg.Body)
	}
}

func (h *StockResponseHandler) processStockResponse(body []byte) {
	// Parse the stock response
	response := string(body)
	parts := strings.Split(response, "|")
	if len(parts) != 4 {
		fmt.Printf("Invalid stock response format: %s\n", response)
		return
	}

	channelID := parts[0]
	_ = parts[1] // userEmail - not used in this context
	botEmail := parts[2]
	message := parts[3]

	// Create the bot message
	botMessage := map[string]interface{}{
		"type": "new_message",
		"data": map[string]interface{}{
			"id":          fmt.Sprintf("bot_%d", time.Now().UnixNano()),
			"channel_id":  channelID,
			"user_email":  botEmail,
			"content":     message,
			"created_at":  time.Now().Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	messageBytes, err := json.Marshal(botMessage)
	if err != nil {
		fmt.Printf("Error marshaling stock response: %v\n", err)
		return
	}

	// Broadcast to WebSocket clients
	h.broadcastFunc(channelID, messageBytes)
	fmt.Printf("Broadcasted stock response to channel %s: %s\n", channelID, message)
}

func (h *StockResponseHandler) Close() error {
	if h.ch != nil {
		h.ch.Close()
	}
	return nil
}
