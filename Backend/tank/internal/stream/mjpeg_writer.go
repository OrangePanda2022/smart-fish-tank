package stream

import (
	"fmt"
	"io"
	"net/http"
)

// MJPEGWriter 将JPEG帧写入multipart/x-mixed-replace HTTP响应流
type MJPEGWriter struct {
	w        io.Writer
	flusher  http.Flusher
	boundary string
}

// NewMJPEGWriter 创建MJPEG写入器
func NewMJPEGWriter(w io.Writer) *MJPEGWriter {
	flusher, ok := w.(http.Flusher)
	if !ok {
		flusher = http.Flusher(nil)
	}
	return &MJPEGWriter{
		w:        w,
		flusher:  flusher,
		boundary: "frameboundary",
	}
}

// Boundary 返回multipart边界字符串，供Content-Type header使用
func (mw *MJPEGWriter) Boundary() string {
	return mw.boundary
}

// WriteFrame 写入一帧JPEG到multipart流
func (mw *MJPEGWriter) WriteFrame(jpeg []byte) error {
	// 写入multipart边界和帧头
	if _, err := fmt.Fprintf(mw.w, "--%s\r\n", mw.boundary); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(mw.w, "Content-Type: image/jpeg\r\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(mw.w, "Content-Length: %d\r\n\r\n", len(jpeg)); err != nil {
		return err
	}
	// 写入JPEG数据
	if _, err := mw.w.Write(jpeg); err != nil {
		return err
	}
	// 写入帧结束标记
	if _, err := fmt.Fprintf(mw.w, "\r\n"); err != nil {
		return err
	}
	// 刷新缓冲区，确保数据立即发送给客户端
	if mw.flusher != nil {
		mw.flusher.Flush()
	}
	return nil
}
