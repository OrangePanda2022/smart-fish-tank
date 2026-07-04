package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"tank/internal/config"
	"tank/internal/controller"
	"tank/internal/domain"
	"tank/internal/infra/db"
	"tank/internal/infra/mq"
	"tank/internal/repo"
	"tank/internal/service"
	"tank/internal/stream"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err)
	}

	// 数据库
	db, _ := db.NewSQLite(cfg.Database.DSN)
	tankRepo := repo.NewTankRepo(db)
	ensureTestTank(tankRepo)

	// 帧缓冲区
	frameBuffer := stream.NewFrameBuffer(cfg.Stream.MaxRingSize, cfg.Stream.ViewerBufSize)

	// 服务层
	tankSvc := service.NewTankService(tankRepo)

	// NATS连接（帧推送和消费共用）
	var natsConn *mq.NATSConsumer
	var nc *nats.Conn // 用于帧元数据发布的原生连接
	natsConn, err = mq.NewNATSConsumer(cfg.NATSURL, tankRepo, frameBuffer)
	if err != nil {
		log.Printf("Failed to create NATS consumer: %s\n", err)
	} else {
		nc = natsConn.NATSClient
		go func() {
			if err := natsConn.Start(context.Background()); err != nil {
				log.Printf("NATS consumer error: %s\n", err)
			}
		}()
	}

	// Handler和路由
	handler := controller.NewHandler(tankSvc, frameBuffer, nc)
	router := controller.NewRouter(handler)

	// HTTP服务器（WriteTimeout=0，支持MJPEG长连接）
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // 必须为0，否则MJPEG流会被超时杀死
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器
	log.Printf("Server started at %s\n", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen error: %s\n", err)
	}
}

func ensureTestTank(tankRepo *repo.TankRepo) {
	const (
		testTankID = "019f1807-dbca-753d-afb8-6b2db160fa8d"
		testUserID = "019ceb4d-95ef-75cd-8364-0144f0b984a7"
	)

	ctx := context.Background()
	existing, err := tankRepo.GetTankByTankID(ctx, testTankID)
	if err != nil {
		log.Printf("Failed to check test tank: %s\n", err)
		return
	}
	if existing != nil {
		return
	}

	if err := tankRepo.Create(ctx, &domain.Tank{
		TankID:   testTankID,
		UserID:   testUserID,
		TankName: "测试鱼缸",
		TankSize: 120,
	}); err != nil {
		log.Printf("Failed to seed test tank: %s\n", err)
		return
	}
	log.Printf("Seeded test tank %s\n", testTankID)
}
