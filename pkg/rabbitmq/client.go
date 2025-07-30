package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Client struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	url           string
	connected     bool
	mu            sync.RWMutex
	reconnectChan chan struct{}
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

type Config struct {
	URL             string
	ReconnectDelay  time.Duration
	MaxReconnectTry int
}

type Message struct {
	Body        []byte
	ContentType string
	Headers     map[string]interface{}
	Exchange    string
	RoutingKey  string
}

type ConsumerConfig struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

func NewClient(config Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	
	client := &Client{
		url:           config.URL,
		reconnectChan: make(chan struct{}, 1),
		ctx:           ctx,
		cancel:        cancel,
	}
	
	return client
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}
	
	c.conn = conn
	c.channel = channel
	c.connected = true
	
	go c.handleReconnect()
	
	return nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.cancel()
	c.wg.Wait()
	
	if c.channel != nil {
		c.channel.Close()
	}
	
	if c.conn != nil {
		c.conn.Close()
	}
	
	c.connected = false
	return nil
}

func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected && c.conn != nil && !c.conn.IsClosed()
}

func (c *Client) handleReconnect() {
	c.wg.Add(1)
	defer c.wg.Done()
	
	notifyClose := make(chan *amqp.Error)
	c.conn.NotifyClose(notifyClose)
	
	for {
		select {
		case err := <-notifyClose:
			if err != nil {
				log.Printf("RabbitMQ connection closed: %v", err)
				c.mu.Lock()
				c.connected = false
				c.mu.Unlock()
				
				for {
					select {
					case <-c.ctx.Done():
						return
					default:
						log.Println("Attempting to reconnect to RabbitMQ...")
						if err := c.reconnect(); err != nil {
							log.Printf("Failed to reconnect: %v", err)
							time.Sleep(5 * time.Second)
							continue
						}
						log.Println("Successfully reconnected to RabbitMQ")
						notifyClose = make(chan *amqp.Error)
						c.conn.NotifyClose(notifyClose)
						break
					}
					break
				}
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) reconnect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	
	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
	}
	if c.channel != nil {
		c.channel.Close()
	}
	
	c.conn = conn
	c.channel = channel
	c.connected = true
	c.mu.Unlock()
	
	return nil
}

func (c *Client) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if !c.connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}
	
	return c.channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (c *Client) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if !c.connected {
		return amqp.Queue{}, fmt.Errorf("not connected to RabbitMQ")
	}
	
	return c.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (c *Client) BindQueue(queueName, key, exchange string, noWait bool, args amqp.Table) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if !c.connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}
	
	return c.channel.QueueBind(queueName, key, exchange, noWait, args)
}

func (c *Client) Publish(exchange, routingKey string, mandatory, immediate bool, msg amqp.Publishing) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if !c.connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}
	
	return c.channel.Publish(exchange, routingKey, mandatory, immediate, msg)
}

func (c *Client) PublishMessage(msg Message) error {
	publishing := amqp.Publishing{
		ContentType: msg.ContentType,
		Body:        msg.Body,
		Headers:     msg.Headers,
		Timestamp:   time.Now(),
	}
	
	return c.Publish(msg.Exchange, msg.RoutingKey, false, false, publishing)
}

func (c *Client) PublishJSON(exchange, routingKey string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	msg := Message{
		Body:        body,
		ContentType: "application/json",
		Exchange:    exchange,
		RoutingKey:  routingKey,
	}
	
	return c.PublishMessage(msg)
}

func (c *Client) Consume(config ConsumerConfig) (<-chan amqp.Delivery, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if !c.connected {
		return nil, fmt.Errorf("not connected to RabbitMQ")
	}
	
	return c.channel.Consume(
		config.Queue,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)
}

func (c *Client) ConsumeWithHandler(config ConsumerConfig, handler func(amqp.Delivery) error) error {
	deliveries, err := c.Consume(config)
	if err != nil {
		return err
	}
	
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case delivery, ok := <-deliveries:
				if !ok {
					log.Println("Consumer channel closed")
					return
				}
				
				if err := handler(delivery); err != nil {
					log.Printf("Error handling message: %v", err)
					delivery.Nack(false, true)
				} else {
					delivery.Ack(false)
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
	
	return nil
}

func (c *Client) SetQoS(prefetchCount, prefetchSize int, global bool) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if !c.connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}
	
	return c.channel.Qos(prefetchCount, prefetchSize, global)
}