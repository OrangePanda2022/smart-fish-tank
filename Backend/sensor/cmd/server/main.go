package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	natsgo "github.com/nats-io/nats.go"

	"sensor/internal/config"
	"sensor/internal/database"
	"sensor/internal/handler"
	"sensor/internal/mqtt"
	"sensor/internal/nats"
	"sensor/internal/repo"
	"sensor/internal/service"
)

func main() {

	// 初始化Gin路由
	router := gin.Default()

	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化SQLite数据库
	sqliteDB, err := database.NewSQLiteDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("连接SQLite失败: %v", err)
	}
	defer sqliteDB.Close()

	// 初始化InfluxDB
	influxDB, err := database.NewInfluxDB(cfg.InfluxDB)
	if err != nil {
		log.Printf("连接InfluxDB失败: %v", err)
		// 继续运行，部分功能可能不可用
	} else {
		defer influxDB.Close()
	}

	// 初始化NATS连接
	var natsConn *natsgo.Conn
	nc, natsErr := natsgo.Connect(cfg.NATSURL)
	if natsErr != nil {
		log.Printf("连接NATS失败: %v", natsErr)
	} else {
		natsConn = nc
		log.Printf("已连接NATS: %s", cfg.NATSURL)
		defer natsConn.Drain()
	}

	// 初始化仓储
	// deviceRepo := repo.NewDeviceRepo(sqliteDB.GetDB())
	sensorRepo := repo.NewSensorRepo(influxDB)

	// 初始化服务
	sensorService := service.NewSensorService(
		sensorRepo,
		natsConn,
		// deviceRepo,
	)

	// 初始化MQTT消费者
	mqttConsumer := mqtt.NewConsumer(cfg.MQTT, sensorService)
	if err := mqttConsumer.Start(); err != nil {
		log.Printf("启动MQTT消费者失败: %v", err)
	} else {
		defer mqttConsumer.Stop()
	}

	// 初始化NATS消费者
	var natsConsumer *nats.NATSConsumer
	if natsConn != nil {
		natsConsumer = nats.NewNATSConsumerWithConn(natsConn, sensorRepo)
		go func() {
			if err := natsConsumer.Start(context.Background()); err != nil {
				log.Printf("NATS消费者异常退出: %v", err)
			}
		}()
		defer natsConsumer.Stop()
	} else {
		// NATS不可用时尝试独立连接
		var natsErr error
		natsConsumer, natsErr = nats.NewNATSConsumer(cfg.NATSURL, sensorRepo)
		if natsErr != nil {
			log.Printf("初始化NATS消费者失败: %v", natsErr)
		} else {
			go func() {
				if err := natsConsumer.Start(context.Background()); err != nil {
					log.Printf("NATS消费者异常退出: %v", err)
				}
			}()
			defer natsConsumer.Stop()
		}
	}

	// 初始化处理器
	sensorHandler := handler.NewSensorHandler(sensorService)
	sensorHandler.RegisterRoutes(router)

	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	fmt.Println(addr)

	// 在goroutine中启动服务器
	go func() {
		log.Printf("启动HTTP服务器于 %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 优雅关闭，设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("服务器强制关闭: %v", err)
	}

	log.Println("服务器已正常退出")
}
