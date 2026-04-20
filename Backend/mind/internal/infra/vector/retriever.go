package vector

import (
	"context"

	"github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// 新建一个 Milvus Retriever 实例
func NewRetriever(ctx context.Context, client *milvusclient.Client, collection string, embedder embedding.Embedder) (*milvus2.Retriever, error) {
	return milvus2.NewRetriever(ctx, &milvus2.RetrieverConfig{
		Client:     client,
		Collection: collection,
		TopK:       5,
		SearchMode: search_mode.NewApproximate(milvus2.COSINE),
		Embedding:  embedder,
	})
}
