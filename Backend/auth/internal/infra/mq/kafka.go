package mq

import (
	"context"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sony/gobreaker"
)

// MessageSender 定义了一个通用的发送接口
type MessageSender interface {
	Send(ctx context.Context, key, value []byte) error
	Close() error
}

// MessageReader 通用读取接口
type MessageReader interface {
	Fetch(ctx context.Context) (kafka.Message, error) // 获取消息（需手动提交）
	Read(ctx context.Context) (kafka.Message, error)  // 自动提交模式
	Commit(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// KafkaProducer 封装了 Kafka Writer 和熔断器
type KafkaProducer struct {
	writer  *kafka.Writer
	breaker *gobreaker.CircuitBreaker
	logger  *slog.Logger
}

// KafkaConsumer 封装了 Kafka Reader
type KafkaConsumer struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewKafkaWriter(brokers []string, topic string, name string, logger *slog.Logger) *KafkaProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireOne,
		Async:        true,
		BatchTimeout: 200 * time.Millisecond,
	}

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		Timeout:     5 * time.Second,
		MaxRequests: 3,
		ReadyToTrip: func(c gobreaker.Counts) bool {
			return c.ConsecutiveFailures >= 5
		},
		OnStateChange: func(n string, from, to gobreaker.State) {
			logger.Warn("[ALERT] MQ circuit state changed",
				slog.String("name", n),
				slog.String("from", from.String()),
				slog.String("to", to.String()))
		},
	})

	return &KafkaProducer{
		writer:  w,
		breaker: cb,
		logger:  logger,
	}
}

// NewKafkaReader 初始化一个消费者
func NewKafkaReader(brokers []string, topic, groupID string, logger *slog.Logger) *KafkaConsumer {
	config := kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	}

	return &KafkaConsumer{
		reader: kafka.NewReader(config),
		logger: logger,
	}
}

func (c *KafkaProducer) Send(ctx context.Context, key, value []byte) error {
	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.writer.WriteMessages(ctx, kafka.Message{
			Key:   key,
			Value: value,
		})
	})
	return err
}

func (c *KafkaProducer) Close() error {
	return c.writer.Close()
}

//

func (c *KafkaConsumer) Read(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *KafkaConsumer) Fetch(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

func (c *KafkaConsumer) Commit(ctx context.Context, msgs ...kafka.Message) error {
	return c.reader.CommitMessages(ctx, msgs...)
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
