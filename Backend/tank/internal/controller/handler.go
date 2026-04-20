package controller

import (
	"fmt"
	"tank/internal/dto/request"
	"tank/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tankSvc *service.TankService
}

func NewHandler(tankSvc *service.TankService) *Handler {
	return &Handler{tankSvc: tankSvc}
}

func (h *Handler) CreateTank(c *gin.Context) {
	var req request.CreateTankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		return
	}
	userID, err := getUserID(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = h.tankSvc.CreateTank(c.Request.Context(), userID, req.TankName, req.TankSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func (h *Handler) GetTankByTankID(c *gin.Context) {
	tankID := c.Param("tank_id")
	if tankID == "" {
		return
	}
	userID, err := getUserID(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	tank, err := h.tankSvc.GetTankByTankID(c.Request.Context(), userID, tankID)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"msg":  "ok",
		"tank": tank,
	})
}

func (h *Handler) GetTanksByUserID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	tanks, err := h.tankSvc.GetTanksByUserID(c.Request.Context(), userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"msg":   "ok",
		"tanks": tanks,
	})
}

func (h *Handler) UpdateTank(c *gin.Context) {
	tankID := c.Param("tank_id")
	if tankID == "" {
		return
	}
	userID, err := getUserID(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	var req request.UpdateTankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		return
	}
	err = h.tankSvc.UpdateTankByTankID(c.Request.Context(), userID, tankID, req.TankName, req.TankSize)
	tank, err := h.tankSvc.GetTankByTankID(c.Request.Context(), userID, tankID)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{
		"msg":  "ok",
		"tank": tank,
	})

}

func (h *Handler) DeleteTankByTankID(c *gin.Context) {
	tankID := c.Param("tank_id")
	if tankID == "" {
		return
	}
	userID, err := getUserID(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	h.tankSvc.DeleteTank(c.Request.Context(), userID, tankID)
	c.JSON(204, gin.H{
		"msg": "No Content",
	})
}

func getUserID(c *gin.Context) (string, error) {
	// userID, exists := c.Get("user_id")
	// if !exists {
	// 	return "", errors.New("userID not found in context")
	// }

	// id, ok := userID.(string)
	// if !ok {
	// 	return "", errors.New("userID type invalid")
	// }

	return "019ceb4d-95ef-75cd-8364-0144f0b984a7", nil
}
