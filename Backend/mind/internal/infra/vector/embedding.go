package vector

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
)

// 创建 Embedding Model
func CreateEmbeddingModel(ctx context.Context) (*ark.Embedder, error) {
	EmbeddingModel, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		Model: "doubao-embedding-text-240715",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding model: %w", err)
	}
	return EmbeddingModel, nil
}

// 生成 Embedding 向量
func GenerateEmbedding(ctx context.Context, embedder *ark.Embedder, input []string) ([][]float64, error) {
	embedding, err := embedder.EmbedStrings(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}
	return embedding, nil
}
