package vector

import (
	"context"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// 新建一个 Milvus 客户端实例
func NewMilvusClient(ctx context.Context, dbname, addr, username, password string) (*milvusclient.Client, error) {
	return milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address:  addr,
		Username: username,
		Password: password,
		DBName:   dbname,
	})
}
