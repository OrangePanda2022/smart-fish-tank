package result

import (
	errorHandler "auth/internal/error"
	"time"

	"github.com/gin-gonic/gin"
)

const requestIDKey = "request_id"

// Body 定义统一返回结构，保证前后端契约稳定
type Body struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Body{
		Code:      "OK",
		Message:   "success",
		Data:      data,
		RequestID: getRequestID(c),
		Timestamp: time.Now().Unix(),
	})
}

// Error 返回统一错误响应
func Error(c *gin.Context, err error) {
	ae := errorHandler.AsAppError(err)
	c.JSON(ae.HTTPStatus, Body{
		Code:      ae.Code,
		Message:   ae.Message,
		RequestID: getRequestID(c),
		Timestamp: time.Now().Unix(),
	})
}

func getRequestID(c *gin.Context) string {
	v, ok := c.Get(requestIDKey)
	if !ok {
		return ""
	}
	id, ok := v.(string)
	if !ok {
		return ""
	}
	return id
}
