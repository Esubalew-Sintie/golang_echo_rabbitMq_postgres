package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"payment-gateway/internal/constant/model/dto"
	"payment-gateway/internal/pkg/logger"
	"payment-gateway/internal/storage"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Queue names
	paymentProcessingQueue = "payment.processing"
	// Exchange names
	paymentExchange = "payment.exchange"
	// Routing keys
	paymentProcessingKey = "payment.processing"
)

// RabbitMQBroker implements the Messaging interface using RabbitMQ with publisher confirms
type RabbitMQBroker struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	queue         amqp.Queue
	confirms      <-chan amqp.Confirmation
	logger        logger.Logger
	connMutex     sync.RWMutex
	consumers     []chan dto.ProcessPaymentMessage
	consumerMutex sync.RWMutex
	running       bool
}

func NewRabbitMQBroker(host, port, user, password, vhost string, logger logger.Logger) (storage.Messaging, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Enable publisher confirms for guaranteed message delivery
	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	// Create confirmation channel
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	err = ch.ExchangeDeclare(
		paymentExchange, // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare(
		paymentProcessingQueue, // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		q.Name,               // queue name
		paymentProcessingKey, // routing key
		paymentExchange,      // exchange
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	broker := &RabbitMQBroker{
		connection: conn,
		channel:    ch,
		queue:      q,
		confirms:   confirms,
		logger:     logger,
		consumers:  make([]chan dto.ProcessPaymentMessage, 0),
		running:    false,
	}

	logger.Info(context.Background(), "RabbitMQ broker initialized successfully")

	return broker, nil
}

// PublishPaymentProcessing publishes a payment processing message to RabbitMQ with guaranteed delivery
func (r *RabbitMQBroker) PublishPaymentProcessing(ctx context.Context, paymentID uuid.UUID) error {
	r.connMutex.RLock()
	defer r.connMutex.RUnlock()

	message := dto.ProcessPaymentMessage{
		PaymentID: paymentID.String(),
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish message with publisher confirms enabled
	err = r.channel.PublishWithContext(
		ctx,
		paymentExchange,      // exchange
		paymentProcessingKey, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Wait for publisher confirmation to ensure message was delivered
	select {
	case confirmed := <-r.confirms:
		if !confirmed.Ack {
			return fmt.Errorf("message delivery not confirmed by RabbitMQ")
		}
		r.logger.Info(ctx, "Message delivery confirmed for payment: %s", paymentID.String())
	case <-ctx.Done():
		return fmt.Errorf("timeout waiting for message delivery confirmation: %w", ctx.Err())
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for message delivery confirmation")
	}

	r.logger.Info(ctx, "Published payment processing message for payment: %s", paymentID.String())
	return nil
}

// ConsumePaymentProcessing starts consuming messages from the payment processing queue
func (r *RabbitMQBroker) ConsumePaymentProcessing() (<-chan dto.ProcessPaymentMessage, error) {
	r.consumerMutex.Lock()
	defer r.consumerMutex.Unlock()

	if !r.running {
		r.running = true
		// Start the message dispatcher in a goroutine
		go r.startDispatching()
	}

	// Create a new consumer channel
	consumerChan := make(chan dto.ProcessPaymentMessage, 100) // Buffered channel
	r.consumers = append(r.consumers, consumerChan)

	r.logger.Info(context.Background(), "Registered new RabbitMQ message consumer")
	return consumerChan, nil
}

// startDispatching consumes messages from RabbitMQ and distributes them to registered consumers
func (r *RabbitMQBroker) startDispatching() {
	r.logger.Info(context.Background(), "Starting RabbitMQ message dispatcher")

	// Set up QoS to ensure fair distribution
	err := r.channel.Qos(1, 0, false) // Prefetch count = 1
	if err != nil {
		r.logger.Error(context.Background(), "Failed to set QoS: %v", err)
		return
	}

	// Start consuming messages
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack (we'll ack manually)
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		r.logger.Error(context.Background(), "Failed to register consumer: %v", err)
		return
	}

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				r.logger.Info(context.Background(), "Message channel closed")
				return
			}

			// Parse the message
			var paymentMsg dto.ProcessPaymentMessage
			if err := json.Unmarshal(msg.Body, &paymentMsg); err != nil {
				r.logger.Error(context.Background(), "Failed to unmarshal message: %v", err)
				msg.Nack(false, false) // Don't requeue malformed messages
				continue
			}

			// Distribute to all registered consumers
			r.consumerMutex.RLock()
			distributed := false
			for _, consumerChan := range r.consumers {
				select {
				case consumerChan <- paymentMsg:
					distributed = true
				default:
					r.logger.Warn(context.Background(), "Consumer channel full, message dropped for one consumer")
				}
			}
			r.consumerMutex.RUnlock()

			if distributed {
				// Acknowledge the message only if it was successfully distributed
				if err := msg.Ack(false); err != nil {
					r.logger.Error(context.Background(), "Failed to acknowledge message: %v", err)
				}
			} else {
				// If no consumers were available, requeue the message
				if err := msg.Nack(false, true); err != nil {
					r.logger.Error(context.Background(), "Failed to nack message: %v", err)
				}
			}
		}
	}
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQBroker) Close() error {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	r.logger.Info(context.Background(), "Closing RabbitMQ broker connection")

	// Close consumer channels
	r.consumerMutex.Lock()
	for _, consumerChan := range r.consumers {
		close(consumerChan)
	}
	r.consumers = nil
	r.consumerMutex.Unlock()

	// Close channel and connection
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}

	if r.connection != nil {
		if err := r.connection.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}

	return nil
}
