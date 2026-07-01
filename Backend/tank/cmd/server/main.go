package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"tank/internal/config"
	"tank/internal/controller"
	"tank/internal/infra/db"
	"tank/internal/infra/mq"
	"tank/internal/repo"
	"tank/internal/service"
	"tank/internal/stream"

	"github.com/nats-io/nats.go"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	database, _ := db.NewSQLite(cfg.Database.DSN)
	tankRepo := repo.NewTankRepo(database)
	tankSvc := service.NewTankService(tankRepo)

	// 直播流帧缓冲区：ESP32-CAM 推帧 -> 环形缓冲 -> 浏览器 MJPEG 拉流
	frameBuffer := stream.NewFrameBuffer(cfg.Stream.MaxRingSize, cfg.Stream.ViewerBufSize)

	// NATS 消费者；连接失败时 natsConn 留 nil，stream_handler 会跳过帧元数据发布
	var natsConn *nats.Conn
	natsConsumer, err := mq.NewNATSConsumer(cfg.NATSURL, tankRepo)
	if err != nil {
		log.Printf("Failed to create NATS consumer: %s\n", err)
	} else {
		natsConn = natsConsumer.NATSClient
		go func() {
			if err := natsConsumer.Start(context.Background()); err != nil {
				log.Printf("NATS consumer error: %s\n", err)
			}
		}()
	}

	handler := controller.NewHandler(tankSvc, frameBuffer, natsConn)
	router := controller.NewRouter(handler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 启动服务器
	log.Printf("Server started at %s\n", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen error: %s\n", err)
	}
}
