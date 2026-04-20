package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"mind/internal/agent/llm"
	"mind/internal/agent/react"
	"mind/internal/agent/tools"
	"mind/internal/config"
	"mind/internal/handler"
	"mind/internal/repo"
	"mind/internal/router"
	"mind/internal/service"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertzlogger "github.com/hertz-contrib/logger/accesslog"
	"github.com/joho/godotenv"
)

func main() {
	cfg := config.Load()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tankRepo := repo.NewTankRepo()
	sensorRepo := repo.NewSensorRepo()
	// 向量数据库初始化
	// milvus, err := vector.NewMilvusClient(context.Background(), cfg.Milvus.DBName, cfg.Milvus.Addr, cfg.Milvus.UserName, cfg.Milvus.Password)
	// if err != nil {
	// 	log.Fatalf("failed to create Milvus client: %v", err)
	// }

	// 向量模型初始化
	// EmbeddingModel, err := vector.CreateEmbeddingModel(context.Background())
	// if err != nil {
	// 	log.Fatalf("failed to create embedding model: %v", err)
	// }

	// 创建检索器
	// retriever, err := vector.NewRetriever(context.Background(), milvus, cfg.Milvus.Collection, EmbeddingModel)
	// if err != nil {
	// 	log.Fatalf("failed to create retriever: %v", err)
	// }

	chatmodel, err := llm.CreateModel(context.Background())

	if err != nil {
		fmt.Printf("Failed to create chat model: %v\n", err)
		os.Exit(1)
	}

	baseTools := []tool.BaseTool{
		tools.NewTankTool(context.Background(), tankRepo),
		tools.NewSensorTool(context.Background(), sensorRepo),
		// tools.NewRAGTool(context.Background(), retriever),
		tools.NewWeatherTool(context.Background()),
		// tools.NewSearchTool(context.Background()),
	}

	// agent 初始化
	// auqaAgent, err := agent.InitAgent(context.Background(), chatmodel, baseTools)
	// if err != nil {
	// 	log.Fatalf("failed to create agent: %v", err)
	// }
	// runner 初始化
	// aquaRunner := agent.InitRunner(context.Background(), auqaAgent, invokableTools)

	// ReAct Agent 初始化
	aquaRAAgent, err := react.ReactAgent(context.Background(), chatmodel, baseTools)

	// MoE 初始化
	// hst, err := moe.NewHost(context.Background(), chatmodel)
	// if err != nil {
	// 	log.Fatalf("failed to create MoE host: %v", err)
	// }
	// experts := []*host.Specialist{
	// 	moe.NewExpert(context.Background(), aquaRAAgent),
	// }
	// summarizer := moe.NewSummarizer(context.Background(), chatmodel)
	// MoE, err := moe.NewMOEAgent(context.Background(), *hst, experts, summarizer)

	// Graph 初始化
	// _, err = flows.AquaGraph(context.Background(), chatmodel, invokableTools)

	if err != nil {
		log.Fatalf("failed to create AquaAgent: %v", err)
	}

	analSvc := service.NewAnalyseService(aquaRAAgent)
	hdl := handler.NewHandler(analSvc)

	h := server.Default(server.WithHostPorts(":" + cfg.Server.Port))
	h.Use(hertzlogger.New())
	router.Register(h, hdl)

	log.Printf("AquaMind server starting on port %s", cfg.Server.Port)

	if err := h.Run(); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
