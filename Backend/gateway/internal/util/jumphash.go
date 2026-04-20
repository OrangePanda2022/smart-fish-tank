package util

import (
	"hash/fnv"
)

// JumpHash 使用 Google 的 Jump Consistent Hash 算法
// key: 用户唯一标识 (如 UserID)
// numBuckets: 当前可用后端实例的数量
func JumpHash(key string, numBuckets int) int {
	// 1. 快速将字符串转换为 uint64
	h := fnv.New64a()
	h.Write([]byte(key))
	hash := h.Sum64()

	// 2. 核心 Jump Hash 算法 (无需额外内存)
	var b int64 = -1
	var j int64 = 0

	for j < int64(numBuckets) {
		b = j
		hash = hash*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(1<<31) / float64((hash>>33)+1)))
	}
	return int(b)
}
