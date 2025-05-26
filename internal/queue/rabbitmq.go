package queue

import (
    "context"
    "fmt"
    "time"

    "github.com/streadway/amqp"
)

// RabbitMQ represents a RabbitMQ connection
type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

// Config holds RabbitMQ configuration
type Config struct {
    URI          string `yaml:"uri"`
    Exchange     string `yaml:"exchange"`
    ExchangeType string `yaml:"exchange_type"`
    Queue        string `yaml:"queue"`
    RoutingKey   string `yaml:"routing_key"`
}

// NewRabbitMQ creates a new RabbitMQ connection
func NewRabbitMQ(cfg Config) (*RabbitMQ, error) {
    conn, err := amqp.Dial(cfg.URI)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to open channel: %w", err)
    }

    // Declare exchange
    err = ch.ExchangeDeclare(
        cfg.Exchange,     // name
        cfg.ExchangeType, // type
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

    // Declare queue
    _, err = ch.QueueDeclare(
        cfg.Queue, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, fmt.Errorf("failed to declare queue: %w", err)
    }

    // Bind queue to exchange
    err = ch.QueueBind(
        cfg.Queue,      // queue name
        cfg.RoutingKey, // routing key
        cfg.Exchange,   // exchange
        false,
        nil,
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, fmt.Errorf("failed to bind queue: %w", err)
    }

    return &RabbitMQ{
        conn:    conn,
        channel: ch,
    }, nil
}

// Publish publishes a message to the queue
func (r *RabbitMQ) Publish(ctx context.Context, body []byte) error {
    return r.channel.Publish(
        "",           // exchange
        "crypto_data", // routing key
        false,        // mandatory
        false,        // immediate
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType:  "application/json",
            Body:        body,
            Timestamp:   time.Now(),
        },
    )
}

// Consume starts consuming messages from the queue
func (r *RabbitMQ) Consume(ctx context.Context, handler func([]byte) error) error {
    msgs, err := r.channel.Consume(
        "crypto_data", // queue
        "",           // consumer
        false,        // auto-ack
        false,        // exclusive
        false,        // no-local
        false,        // no-wait
        nil,          // args
    )
    if err != nil {
        return fmt.Errorf("failed to register a consumer: %w", err)
    }

    for {
        select {
        case <-ctx.Done():
            return nil
        case msg := <-msgs:
            // 处理消息，包含重试机制
            for retries := 0; retries < 3; retries++ {
                err := handler(msg.Body)
                if err == nil {
                    msg.Ack(false) // 确认消息
                    break
                }
                if retries == 2 {
                    // 最后一次重试失败，拒绝消息并重新入队
                    msg.Reject(true)
                }
                time.Sleep(time.Second * time.Duration(retries+1))
            }
        }
    }
}

// Close closes the RabbitMQ connection
func (r *RabbitMQ) Close() error {
    if err := r.channel.Close(); err != nil {
        return fmt.Errorf("failed to close channel: %w", err)
    }
    if err := r.conn.Close(); err != nil {
        return fmt.Errorf("failed to close connection: %w", err)
    }
    return nil
} 