package buffer

import (
	"sync"
)

type BytePool struct {
	Pool sync.Pool
}

func (p *BytePool) Get() []byte {
	return p.Pool.Get().([]byte)
}

func (p *BytePool) Put(b []byte) {
	p.Pool.Put(b)
}

func NewBufferPool() *BytePool {
	return &BytePool{
		Pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 32*1024) // 32KB 缓冲区
			},
		},
	}
}
