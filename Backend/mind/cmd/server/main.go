package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"mind/internal/agent/llm"
	"mind/internal/agent/react"
	"mind/internal/agent/tools"
	"mind/internal/config"
	"mind/internal/domain"
	"mind/internal/handler"
	"mind/internal/predict/koopman"
	"mind/internal/repo"
	"mind/internal/router"
	"mind/internal/service"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertzlogger "github.com/hertz-contrib/logger/accesslog"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	// 先加载 .env，再读配置——否则 .env 里的 NATS_URL/LLM_API_KEY 等对 config.Load() 不可见
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.Load()

	// NATS 仓储：向 tank/sensor 服务实时拉取数据；连接失败不致命，工具/预测按调用返回错误给 agent。
	// 用具体指针变量持有连接（供 Close/Conn），用接口变量喂给工具/预测服务——
	// 这样 NATS 失败时接口为真 nil，规避 Go nil-interface 陷阱，使工具与 predict 的 nil 守卫真正生效。
	// （未采用副本的内存 mock 降级，保留 SHIT2 的 nil 守卫路线。）
	var tankRepo repo.TankRepository
	var sensorRepo repo.SensorRepository
	var frameRepo repo.FrameRepository

	tankNATS, tankErr := repo.NewTankNATSRepo(cfg.NATS.URL)
	if tankErr != nil {
		log.Printf("WARN: tank NATS repo init failed: %v (tank tool will error per-call)", tankErr)
	} else {
		tankRepo = tankNATS
		frameRepo = tankNATS // *TankNATSRepo 同时满足 TankRepository 与 FrameRepository
		defer tankNATS.Close()
	}

	sensorNATS, sensorErr := repo.NewSensorNATSRepo(cfg.NATS.URL)
	if sensorErr != nil {
		log.Printf("WARN: sensor NATS repo init failed: %v (sensor tool will error per-call)", sensorErr)
	} else {
		sensorRepo = sensorNATS
		defer sensorNATS.Close()
	}

	chatmodel, err := llm.CreateModel(context.Background())
	if err != nil {
		fmt.Printf("Failed to create chat model: %v\n", err)
		os.Exit(1)
	}

	baseTools := []tool.BaseTool{
		tools.NewTankTool(context.Background(), tankRepo),
		tools.NewSensorTool(context.Background(), sensorRepo),
		// tools.NewRAGTool(context.Background(), retriever),
		// tools.NewWeatherTool(context.Background()),
		// tools.NewSearchTool(context.Background()),
	}

	// ReAct Agent 初始化
	aquaRAAgent, err := react.ReactAgent(context.Background(), chatmodel, baseTools)
	if err != nil {
		log.Fatalf("failed to create AquaAgent: %v", err)
	}

	analSvc := service.NewAnalyseService(aquaRAAgent, frameRepo)

	// 加载 Koopman 预测模型；加载失败不致命，Predict 端点会返回 503
	koopmanModel, err := koopman.LoadModel(cfg.Prediction.ModelPath)
	if err != nil {
		log.Printf("Koopman模型加载失败: %v，预测服务不可用", err)
		koopmanModel = &koopman.Model{} // 空模型, IsLoaded()=false
	}

	predictSvc := service.NewPredictService(
		koopmanModel,
		sensorRepo,
		cfg.Prediction.Horizon,
		cfg.Prediction.StateWeight,
		cfg.Prediction.CtrlWeight,
	)

	// 后台: 订阅 NATS sensor 事件，持续更新预测状态缓存
	if sensorNATS != nil {
		go func() {
			nc := sensorNATS.Conn()
			if nc == nil {
				return
			}
			sub, err := nc.Subscribe("sensor.update", func(msg *nats.Msg) {
				var data domain.SensorData
				if json.Unmarshal(msg.Data, &data) == nil {
					predictSvc.UpdateState(&data)
				}
			})
			if err != nil {
				log.Printf("NATS sensor.update 订阅失败: %v", err)
				return
			}
			defer sub.Unsubscribe()
			select {} // 阻塞保持订阅
		}()
	}

	hdl := handler.NewHandler(analSvc, predictSvc)

	h := server.Default(server.WithHostPorts(":" + cfg.Server.Port))
	h.Use(hertzlogger.New())
	router.Register(h, hdl)

	log.Printf("AquaMind server starting on port %s", cfg.Server.Port)

	if err := h.Run(); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
