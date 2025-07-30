package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"

	"github.com/brianchul/airline_booking/pkg/api"
	"github.com/brianchul/airline_booking/pkg/rabbitmq"
)

type BookingQueue interface {
	ProduceNormalBookingQueue(message *api.BookFlightRequest) error
	ProduceVipBookingQueue(message *api.BookFlightRequest) error
	ConsumeNormalBookingQueue(handler BookingHandler) error
	ConsumeVipBookingQueue(handler BookingHandler) error
	Start() error
	Stop() error
	IsHealthy() bool
}

// BookingHandler is a function type for handling booking requests
type BookingHandler func(ctx context.Context, request *BookingRequest) error

// RabbitMQBookingQueue implements BookingQueue using RabbitMQ
type RabbitMQBookingQueue struct {
	client        *rabbitmq.Client
	normalHandler BookingHandler
	vipHandler    BookingHandler
	config        *BookingQueueConfig
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	running       bool
	mu            sync.RWMutex
}

// BookingQueueConfig holds configuration for the booking queue
type BookingQueueConfig struct {
	ExchangeName       string
	NormalQueueName    string
	VIPQueueName       string
	NormalRoutingKey   string
	VIPRoutingKey      string
	DeadLetterExchange string
	PrefetchCount      int
}

// NewRabbitMQBookingQueue creates a new RabbitMQ booking queue
func NewRabbitMQBookingQueue(client *rabbitmq.Client, config *BookingQueueConfig) BookingQueue {
	ctx, cancel := context.WithCancel(context.Background())

	return &RabbitMQBookingQueue{
		client: client,
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

// SetNormalHandler sets the handler for normal booking queue
func (r *RabbitMQBookingQueue) SetNormalHandler(handler BookingHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.normalHandler = handler
}

// SetVIPHandler sets the handler for VIP booking queue
func (r *RabbitMQBookingQueue) SetVIPHandler(handler BookingHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vipHandler = handler
}

// ProduceNormalBookingQueue produces a message to the normal booking queue
func (r *RabbitMQBookingQueue) ProduceNormalBookingQueue(message *api.BookFlightRequest) error {
	return r.publishMessage(message, r.config.NormalRoutingKey)
}

// ProduceVipBookingQueue produces a message to the VIP booking queue
func (r *RabbitMQBookingQueue) ProduceVipBookingQueue(message *api.BookFlightRequest) error {
	return r.publishMessage(message, r.config.VIPRoutingKey)
}

// publishMessage publishes a booking message to the specified routing key
func (r *RabbitMQBookingQueue) publishMessage(message *api.BookFlightRequest, routingKey string) error {
	if !r.client.IsConnected() {
		return fmt.Errorf("rabbitmq client is not connected")
	}


	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal booking request: %w", err)
	}

	msg := rabbitmq.Message{
		Body:        body,
		ContentType: "application/json",
		Exchange:    r.config.ExchangeName,
		RoutingKey:  routingKey,
		Headers: map[string]interface{}{
			"message_type":    "booking_request",
			"user_tier":       string(message.UserTier),
			"flight_number":   message.FlightNumber,
			"timestamp":       time.Now().Unix(),
		},
	}

	if err := r.client.PublishMessage(msg); err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %w", routingKey, err)
	}

	log.Printf("Successfully published booking request for flight %s to queue %s", message.FlightNumber, routingKey)
	return nil
}

// ConsumeNormalBookingQueue starts consuming messages from the normal booking queue
func (r *RabbitMQBookingQueue) ConsumeNormalBookingQueue(handler BookingHandler) error {
	r.mu.Lock()
	r.normalHandler = handler
	r.mu.Unlock()

	return r.startConsumer(r.config.NormalQueueName, "normal_booking_consumer", handler)
}

// ConsumeVipBookingQueue starts consuming messages from the VIP booking queue
func (r *RabbitMQBookingQueue) ConsumeVipBookingQueue(handler BookingHandler) error {
	r.mu.Lock()
	r.vipHandler = handler
	r.mu.Unlock()

	return r.startConsumer(r.config.VIPQueueName, "vip_booking_consumer", handler)
}

// startConsumer starts a consumer for the specified queue
func (r *RabbitMQBookingQueue) startConsumer(queueName, consumerTag string, handler BookingHandler) error {
	if !r.client.IsConnected() {
		return fmt.Errorf("rabbitmq client is not connected")
	}

	config := rabbitmq.ConsumerConfig{
		Queue:     queueName,
		Consumer:  consumerTag,
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}

	messageHandler := func(delivery amqp.Delivery) error {
		var request BookingRequest
		if err := json.Unmarshal(delivery.Body, &request); err != nil {
			log.Printf("Failed to unmarshal booking request: %v", err)
			return err
		}

		log.Printf("Processing booking request %s from queue %s", request.BookingUUID, queueName)

		ctx, cancel := context.WithTimeout(r.ctx, 30*time.Second)
		defer cancel()

		if err := handler(ctx, &request); err != nil {
			log.Printf("Failed to process booking request %s: %v", request.BookingUUID, err)
			return err
		}

		log.Printf("Successfully processed booking request %s", request.BookingUUID)
		return nil
	}

	return r.client.ConsumeWithHandler(config, messageHandler)
}

// Start initializes the queues and exchanges
func (r *RabbitMQBookingQueue) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return fmt.Errorf("booking queue is already running")
	}

	if !r.client.IsConnected() {
		return fmt.Errorf("rabbitmq client is not connected")
	}

	// Declare exchange
	if err := r.client.DeclareExchange(
		r.config.ExchangeName,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // args
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare dead letter exchange if configured
	if r.config.DeadLetterExchange != "" {
		if err := r.client.DeclareExchange(
			r.config.DeadLetterExchange,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to declare dead letter exchange: %w", err)
		}
	}

	// Declare queues
	queueArgs := amqp.Table{}
	if r.config.DeadLetterExchange != "" {
		queueArgs["x-dead-letter-exchange"] = r.config.DeadLetterExchange
	}

	// Normal queue
	if _, err := r.client.DeclareQueue(
		r.config.NormalQueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		queueArgs,
	); err != nil {
		return fmt.Errorf("failed to declare normal queue: %w", err)
	}

	// VIP queue
	if _, err := r.client.DeclareQueue(
		r.config.VIPQueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		queueArgs,
	); err != nil {
		return fmt.Errorf("failed to declare VIP queue: %w", err)
	}

	// Bind queues to exchange
	if err := r.client.BindQueue(
		r.config.NormalQueueName,
		r.config.NormalRoutingKey,
		r.config.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind normal queue: %w", err)
	}

	if err := r.client.BindQueue(
		r.config.VIPQueueName,
		r.config.VIPRoutingKey,
		r.config.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind VIP queue: %w", err)
	}

	// Set QoS
	if r.config.PrefetchCount > 0 {
		if err := r.client.SetQoS(r.config.PrefetchCount, 0, false); err != nil {
			return fmt.Errorf("failed to set QoS: %w", err)
		}
	}

	r.running = true
	log.Printf("Booking queue started successfully with exchange: %s", r.config.ExchangeName)
	return nil
}

// Stop stops the booking queue
func (r *RabbitMQBookingQueue) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.running {
		return nil
	}

	r.cancel()
	r.wg.Wait()
	r.running = false

	log.Println("Booking queue stopped")
	return nil
}

// IsHealthy checks if the booking queue is healthy
func (r *RabbitMQBookingQueue) IsHealthy() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.running && r.client.IsConnected()
}
