package nats

import (
	"context"
	"encoding/json"
	"sensor/internal/domain"
	"sensor/internal/repo"
	"time"

	"github.com/nats-io/nats.go"
)

type NATSConsumer struct {
	NATSClient *nats.Conn
	sensorRepo *repo.SensorRepo
}

func NewNATSConsumer(URL string, sensorRepo *repo.SensorRepo) (*NATSConsumer, error) {
	nc, err := nats.Connect(URL)
	if err != nil {
		return nil, err
	}
	return &NATSConsumer{
		NATSClient: nc,
		sensorRepo: sensorRepo,
	}, nil
}

func (c *NATSConsumer) Start(ctx context.Context) error {
	_, err := c.NATSClient.QueueSubscribe("sensor.get", "sensor-group", func(msg *nats.Msg) {
		var req domain.Request

		// 解析请求
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			msg.Respond([]byte(`{"error":"bad request"}`))
			return
		}

		var resp interface{}

		switch req.Type {
		case "sensor":
			sensorData, err := c.sensorRepo.QueryHistoryByTank(
				req.TankID,
				time.Now().Add(-48*time.Hour).Format(time.RFC3339),
				time.Now().Format(time.RFC3339),
				req.Limit,
			)
			if err != nil {
				resp = map[string]string{"error": err.Error()}
			} else {
				resp = sensorData
			}
			break
		default:
			resp = map[string]string{"error": "bad request"}
		}

		// 返回响应
		data, err := json.Marshal(resp)
		if err != nil {
			msg.Respond([]byte(`{"error":"marshal error"}`))
			return
		}

		msg.Respond(data)
	})

	if err != nil {
		return err
	}

	// 确保订阅生效
	if err := c.NATSClient.Flush(); err != nil {
		return err
	}

	// 检查连接状态
	if err := c.NATSClient.LastError(); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (c *NATSConsumer) Stop() error {
	if c.NATSClient == nil {
		return nil
	}

	c.NATSClient.Drain()
	c.NATSClient.Close()
	return nil
}
