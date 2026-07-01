package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"tank/internal/stream"
	"time"

	"github.com/gin-gonic/gin"
)

// PostFrame 处理ESP32-CAM推送的JPEG帧
// POST /api/v1/tank/:tank_id/stream/frame
func (h *Handler) PostFrame(c *gin.Context) {
	tankID := c.Param("tank_id")
	if tankID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error", "error": "missing tank_id"})
		return
	}

	// 读取raw body（JPEG二进制数据）
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error", "error": "empty body"})
		return
	}

	// 校验JPEG SOI标记（0xFF 0xD8）
	if len(body) < 2 || body[0] != 0xFF || body[1] != 0xD8 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error", "error": "not a JPEG"})
		return
	}

	// TODO: 设备认证 — 验证X-Device-Key请求头
	// 当前MVP阶段不做设备认证

	// 写入帧缓冲区
	h.frameBuffer.PushFrame(tankID, body)

	// 可选：通过NATS发布帧元数据事件（不含JPEG数据）
	if h.natsConn != nil {
		meta := map[string]interface{}{
			"tank_id":   tankID,
			"size":      len(body),
			"timestamp": time.Now().UnixMilli(),
		}
		data, _ := json.Marshal(meta)
		if err := h.natsConn.Publish("tank.video."+tankID, data); err != nil {
			log.Printf("NATS publish error: %s\n", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// StreamVideo 返回MJPEG直播流，浏览器可用<img>标签直接显示
// GET /api/v1/tank/:tank_id/stream
func (h *Handler) StreamVideo(c *gin.Context) {
	tankID := c.Param("tank_id")
	if tankID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error", "error": "missing tank_id"})
		return
	}

	// 验证用户对该鱼缸的所有权
	userID, err := getUserID(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = h.tankSvc.GetTankByTankID(c.Request.Context(), userID, tankID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"msg": "error", "error": "forbidden"})
		return
	}

	// 注册为观看者
	frameCh := h.frameBuffer.RegisterViewer(tankID)
	defer h.frameBuffer.DeregisterViewer(tankID, frameCh)

	// 设置MJPEG流响应头
	mjpeg := stream.NewMJPEGWriter(c.Writer)
	c.Writer.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", mjpeg.Boundary()))
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	// 流式推送帧循环
	for {
		select {
		case jpeg, ok := <-frameCh:
			if !ok {
				// channel已关闭（鱼缸帧存储被删除）
				return
			}
			if err := mjpeg.WriteFrame(jpeg); err != nil {
				// 客户端断开连接
				return
			}
		case <-c.Request.Context().Done():
			// HTTP连接关闭
			return
		}
	}
}
