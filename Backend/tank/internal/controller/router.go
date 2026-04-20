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
		}
	}

	return r
}
