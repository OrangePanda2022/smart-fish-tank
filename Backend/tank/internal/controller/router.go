package controller

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		tankGroup := v1.Group("/tank")
		{
			tankGroup.PATCH("/", handler.CreateTank)
			tankGroup.GET("/", handler.GetTanksByUserID)
			tankGroup.GET("/:tank_id", handler.GetTankByTankID)
			tankGroup.PUT("/:tank_id", handler.UpdateTank)
			tankGroup.DELETE("/:tank_id", handler.DeleteTankByTankID)
			// 直播流：ESP32-CAM 推帧 / 浏览器拉 MJPEG
			tankGroup.POST("/:tank_id/stream/frame", handler.PostFrame)
			tankGroup.GET("/:tank_id/stream", handler.StreamVideo)
		}
	}

	return r
}
