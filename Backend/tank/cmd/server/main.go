package main

import (
	"context"
	"log"
	"net/http"
	"tank/internal/controller"
	"tank/internal/infra/db"
	"tank/internal/infra/mq"
	"tank/internal/repo"
	"tank/internal/service"
)

func main() {
	db, _ := db.NewSQLite("./tank.db")
	tankRepo := repo.NewTankRepo(db)
	tankSvc := service.NewTankService(tankRepo)
	handler := controller.NewHandler(tankSvc)
	router := controller.NewRouter(handler)
	nats, err := mq.NewNATSConsumer("nats://localhost:4222", tankRepo)
	if err != nil {
		log.Printf("Failed to create NATS consumer: %s\n", err)
	} else {
		go func() {
			if err := nats.Start(context.Background()); err != nil {
				log.Printf("NATS consumer error: %s\n", err)
			}
		}()
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:8084",
		Handler: router,
	}

	// 启动服务器
	log.Println("Server started at :8084")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen error: %s\n", err)
	}

}
