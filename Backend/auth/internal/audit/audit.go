package audit

import (
	"auth/internal/domain/event"
	"auth/internal/infra/mq"
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

type AuditLogger interface {
	Log(ctx context.Context, event event.AuditEvent)
	Close() error
}

type KafkaAuditLogger struct {
	sender mq.MessageSender
	logger *slog.Logger
}

func NewKafkaAuditLogger(brokers []string, topic string, logger *slog.Logger) AuditLogger {
	// 未配置 Kafka 时直接降级为 stdout
	if len(brokers) == 0 || topic == "" {
		logger.Warn("kafka disabled, fallback to stdout audit logger")
		return &noopAuditLogger{logger: logger}
	}

	// 调用通用的初始化逻辑
	client := mq.NewKafkaWriter(brokers, topic, "audit-breaker", logger)

	return &KafkaAuditLogger{
		sender: client,
		logger: logger,
	}
}

func (l *KafkaAuditLogger) Log(ctx context.Context, event event.AuditEvent) {
	event.At = time.Now()
	payload, err := json.Marshal(event)
	if err != nil {
		l.logger.Error("marshal audit event failed", slog.String("error", err.Error()))
		return
	}

	// 使用通用的 Send 方法
	err = l.sender.Send(ctx, nil, payload)
	if err != nil {
		l.logger.Warn("kafka write failed, fallback local audit",
			slog.String("error", err.Error()),
			slog.String("event", string(payload)),
		)
	}
}

func (l *KafkaAuditLogger) Close() error {
	return l.sender.Close()
}

// 无 MQ 场景
type noopAuditLogger struct {
	logger *slog.Logger
}

func (n *noopAuditLogger) Log(_ context.Context, event event.AuditEvent) {
	// 无 MQ 场景下仍保留结构化审计记录，便于本地开发排查
	event.At = time.Now()
	data, _ := json.Marshal(event)
	n.logger.Info("audit event", slog.String("event", string(data)))
}

// Close 为 no-op，实现接口契约
func (n *noopAuditLogger) Close() error { return nil }
