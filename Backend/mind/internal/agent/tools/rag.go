package tools

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type RAGQueryParams struct {
	Query string `json:"query" jsonschema:"description=需要检索的水族专业知识、疾病特征或设备故障排查指南的关键词"`
}

func NewRAGTool(ctx context.Context, retriever *milvus2.Retriever) tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "RAGQueryTool",
			Desc: "基于输入的关键词，在水族专业知识库中进行检索，返回相关的知识片段",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"query": {
						Type:     schema.String,
						Desc:     "需要检索的水族专业知识、疾病特征或设备故障排查指南的关键词",
						Required: true,
					},
				},
			),
		}, func(ctx context.Context, query *RAGQueryParams) (string, error) {
			docs, err := retriever.Retrieve(ctx, query.Query)
			if err != nil {
				return "", fmt.Errorf("检索知识库失败: %w", err)
			}

			if len(docs) == 0 {
				return "未检索到相关知识。", nil
			}

			// 将检索到的知识切片拼接成字符串返回给大模型
			var result string
			for i, doc := range docs {
				result += fmt.Sprintf("[%d] %s\n", i+1, doc)
			}
			return result, nil
		},
	)
}
