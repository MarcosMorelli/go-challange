package service

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

type StockBot struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type StockCommand struct {
	ChannelID string `json:"channel_id"`
	UserEmail string `json:"user_email"`
	StockCode string `json:"stock_code"`
}

type StockResponse struct {
	ChannelID string `json:"channel_id"`
	Message   string `json:"message"`
	UserEmail string `json:"user_email"`
}

func NewStockBot(conn *amqp.Connection) (*StockBot, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &StockBot{
		conn: conn,
		ch:   ch,
	}, nil
}

func (bot *StockBot) Start() error {
	// Setup queues
	err := bot.setupQueues()
	if err != nil {
		return err
	}

	// Start consuming stock commands
	go bot.consumeStockCommands()

	return nil
}

func (bot *StockBot) setupQueues() error {
	// Declare stock_commands queue
	_, err := bot.ch.QueueDeclare(
		"stock_commands", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	// Declare stock_responses queue
	_, err = bot.ch.QueueDeclare(
		"stock_responses", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func (bot *StockBot) consumeStockCommands() {
	msgs, err := bot.ch.Consume(
		"stock_commands", // queue
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		fmt.Printf("Failed to register a consumer: %s\n", err)
		return
	}

	for msg := range msgs {
		bot.processStockCommand(msg.Body)
	}
}

func (bot *StockBot) processStockCommand(body []byte) {
	// Parse the stock command
	command := string(body)
	parts := strings.Split(command, "|")
	if len(parts) != 3 {
		fmt.Printf("Invalid stock command format: %s\n", command)
		return
	}

	channelID := parts[0]
	userEmail := parts[1]
	stockCode := parts[2]

	// Fetch stock data
	stockData, err := bot.fetchStockData(stockCode)
	if err != nil {
		fmt.Printf("Error fetching stock data for %s: %v\n", stockCode, err)
		// Send error message back
		errorMsg := fmt.Sprintf("Error fetching stock data for %s", stockCode)
		bot.sendStockResponse(channelID, userEmail, errorMsg)
		return
	}

	// Format the response
	response := fmt.Sprintf("%s quote is $%.2f per share", strings.ToUpper(stockCode), stockData.Price)
	bot.sendStockResponse(channelID, userEmail, response)
}

func (bot *StockBot) fetchStockData(stockCode string) (*StockData, error) {
	url := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stockCode)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse CSV response
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("invalid CSV response")
	}

	// Skip header row, get first data row
	dataRow := records[1]
	if len(dataRow) < 5 {
		return nil, fmt.Errorf("invalid CSV data format")
	}

	// Parse the close price (index 4 in the CSV)
	var price float64
	_, err = fmt.Sscanf(dataRow[4], "%f", &price)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %s", dataRow[4])
	}

	return &StockData{
		Symbol: dataRow[0],
		Price:  price,
	}, nil
}

func (bot *StockBot) sendStockResponse(channelID, userEmail, message string) {
	response := fmt.Sprintf("%s|%s|%s|%s", channelID, userEmail, "stock_bot", message)

	err := bot.ch.Publish(
		"",                // exchange
		"stock_responses", // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(response),
		},
	)

	if err != nil {
		fmt.Printf("Failed to publish stock response: %s\n", err)
	}
}

func (bot *StockBot) Close() error {
	if bot.ch != nil {
		bot.ch.Close()
	}
	if bot.conn != nil {
		bot.conn.Close()
	}
	return nil
}

type StockData struct {
	Symbol string
	Price  float64
}
