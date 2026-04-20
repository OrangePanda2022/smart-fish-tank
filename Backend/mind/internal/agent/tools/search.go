package tools

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/bingsearch"
	"github.com/cloudwego/eino/components/tool"
)

func NewSearchTool(ctx context.Context) tool.InvokableTool {
	bingSearchAPIKey := os.Getenv("BING_SEARCH_API_KEY")

	// 创建 Bing Search 工具
	bingSearchTool, err := bingsearch.NewTool(ctx, &bingsearch.Config{
		APIKey: bingSearchAPIKey,
		Cache:  5 * time.Minute,
	})

	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}

	return bingSearchTool
}
