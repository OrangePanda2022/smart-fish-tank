package stream

import (
	"sync"
)

// FrameBuffer 管理每个鱼缸的JPEG帧环形缓冲区和观看者注册
type FrameBuffer struct {
	mu      sync.RWMutex
	stores  map[string]*tankFrameStore // key: tank_id
	ringCap int                        // 环形缓冲区容量
	viewBuf int                        // 观看者channel缓冲区大小
}

// tankFrameStore 存储单个鱼缸的帧数据和观看者集合
type tankFrameStore struct {
	ring    [][]byte                  // 固定大小JPEG帧环形缓冲区
	head    int                       // 最新帧的索引位置
	count   int                       // 缓冲区中已有帧的数量
	viewers map[chan []byte]struct{}  // 观看者通知channel集合
}

// NewFrameBuffer 创建帧缓冲区
func NewFrameBuffer(ringCap, viewBuf int) *FrameBuffer {
	return &FrameBuffer{
		stores:  make(map[string]*tankFrameStore),
		ringCap: ringCap,
		viewBuf: viewBuf,
	}
}

// PushFrame 写入一帧JPEG数据，并通知所有观看者
func (fb *FrameBuffer) PushFrame(tankID string, jpeg []byte) {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	store, ok := fb.stores[tankID]
	if !ok {
		store = &tankFrameStore{
			ring:    make([][]byte, fb.ringCap),
			viewers: make(map[chan []byte]struct{}),
		}
		fb.stores[tankID] = store
	}

	// 复制JPEG数据到环形缓冲区
	store.ring[store.head] = make([]byte, len(jpeg))
	copy(store.ring[store.head], jpeg)

	// 推进head指针
	store.head = (store.head + 1) % fb.ringCap
	if store.count < fb.ringCap {
		store.count++
	}

	// 非阻塞通知所有观看者
	for ch := range store.viewers {
		select {
		case ch <- jpeg:
		default:
			// 观看者读取慢，丢弃此帧
		}
	}
}

// RegisterViewer 注册一个观看者，返回帧通知channel
func (fb *FrameBuffer) RegisterViewer(tankID string) chan []byte {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	store, ok := fb.stores[tankID]
	if !ok {
		store = &tankFrameStore{
			ring:    make([][]byte, fb.ringCap),
			viewers: make(map[chan []byte]struct{}),
		}
		fb.stores[tankID] = store
	}

	ch := make(chan []byte, fb.viewBuf)
	store.viewers[ch] = struct{}{}
	return ch
}

// DeregisterViewer 注销观看者，关闭channel
func (fb *FrameBuffer) DeregisterViewer(tankID string, ch chan []byte) {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	store, ok := fb.stores[tankID]
	if !ok {
		return
	}

	delete(store.viewers, ch)
	close(ch)
}

// GetLatestFrame 获取最新的一帧JPEG数据（用于快照端点）
func (fb *FrameBuffer) GetLatestFrame(tankID string) []byte {
	fb.mu.RLock()
	defer fb.mu.RUnlock()

	store, ok := fb.stores[tankID]
	if !ok || store.count == 0 {
		return nil
	}

	// head指向下一个写入位置，最新帧在head-1
	idx := (store.head - 1 + fb.ringCap) % fb.ringCap
	return store.ring[idx]
}

// RemoveTank 删除鱼缸的帧存储，关闭所有观看者channel
func (fb *FrameBuffer) RemoveTank(tankID string) {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	store, ok := fb.stores[tankID]
	if !ok {
		return
	}

	for ch := range store.viewers {
		close(ch)
	}
	delete(fb.stores, tankID)
}
