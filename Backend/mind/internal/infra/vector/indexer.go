package vector

import (
	"context"

	"github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// 新建一个 Milvus Indexer 实例
func NewIndexer(ctx context.Context, client *milvusclient.Client, collection string, embedder embedding.Embedder) (*milvus2.Indexer, error) {
	return milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		Client:     client,
		Collection: collection,
		Vector: &milvus2.VectorConfig{
			Dimension:    512,
			MetricType:   milvus2.COSINE,
			IndexBuilder: milvus2.NewHNSWIndexBuilder().WithM(16).WithEfConstruction(200),
		},
		Embedding: embedder,
	})
}
