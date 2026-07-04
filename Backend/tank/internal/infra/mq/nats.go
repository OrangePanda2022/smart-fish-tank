package mq

import (
	"context"
	"encoding/json"
	"tank/internal/dto/request"
	"tank/internal/repo"
	"tank/internal/stream"

	"github.com/nats-io/nats.go"
)

type NATSConsumer struct {
	NATSClient  *nats.Conn
	tankRepo    *repo.TankRepo
	frameBuffer *stream.FrameBuffer
}

func NewNATSConsumer(URL string, tankRepo *repo.TankRepo, frameBuffer *stream.FrameBuffer) (*NATSConsumer, error) {
	nc, err := nats.Connect(URL)
	if err != nil {
		return nil, err
	}
	return &NATSConsumer{
		NATSClient:  nc,
		tankRepo:    tankRepo,
		frameBuffer: frameBuffer,
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

	// tank.frame: 返回最新一帧 JPEG（供 mind 服务做多模态分析）
	// 契约: 成功回原始 JPEG 字节；无帧/坏请求回 {"error":"..."}。调用方按 SOI(0xFF 0xD8) 判定有效帧。
	if _, err := c.NATSClient.QueueSubscribe("tank.frame", "tank-group", func(msg *nats.Msg) {
		var req struct {
			TankID string `json:"tank_id"`
		}
		if err := json.Unmarshal(msg.Data, &req); err != nil || req.TankID == "" {
			msg.Respond([]byte(`{"error":"bad request"}`))
			return
		}
		jpeg := c.frameBuffer.GetLatestFrame(req.TankID)
		if len(jpeg) == 0 {
			msg.Respond([]byte(`{"error":"no frame"}`))
			return
		}
		msg.Respond(jpeg) // 原始 JPEG 字节，不是 JSON
	}); err != nil {
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
