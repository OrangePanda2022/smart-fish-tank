package mq

import (
	"context"
	"encoding/json"
	"tank/internal/dto/request"
	"tank/internal/repo"

	"github.com/nats-io/nats.go"
)

type NATSConsumer struct {
	NATSClient *nats.Conn
	tankRepo   *repo.TankRepo
}

func NewNATSConsumer(URL string, tankRepo *repo.TankRepo) (*NATSConsumer, error) {
	nc, err := nats.Connect(URL)
	if err != nil {
		return nil, err
	}
	return &NATSConsumer{
		NATSClient: nc,
		tankRepo:   tankRepo,
	}, nil
}

func (c *NATSConsumer) Start(ctx context.Context) error {
	_, err := c.NATSClient.QueueSubscribe("tank.get", "tank-group", func(msg *nats.Msg) {
		var req request.GetTankRequest

		// 解析请求
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			msg.Respond([]byte(`{"error":"bad request"}`))
			return
		}

		// subject 已路由到 tank 服务，无需 type 字段；tank 查询要求 limit == 1
		if req.Limit != 1 {
			msg.Respond([]byte(`{"error":"bad request"}`))
			return
		}

		tankData, err := c.tankRepo.GetTankByTankID(
			ctx,
			req.TankID,
		)
		if err != nil {
			data, _ := json.Marshal(map[string]string{"error": err.Error()})
			msg.Respond(data)
			return
		}

		data, err := json.Marshal(tankData)
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
