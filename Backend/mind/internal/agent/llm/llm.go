package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/claude"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// 创建 Model

func CreateModel(ctx context.Context) (*ark.ChatModel, error) {
	return ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("API_KEY"),
		Model:  "doubao-seed-2-0-mini-260215",
	})
}

// 创建消息模板
func CreateMessageTemplate() *prompt.DefaultChatTemplate {
	return prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}。请用{style}语气回答问题。"),
		schema.UserMessage("问题: {question}"),
	)
}

// 生成回答
func GenerateAnswer(ctx context.Context, model *claude.ChatModel, promptTexts []*schema.Message) (*schema.Message, error) {
	response, err := model.Generate(ctx, promptTexts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate answer: %w", err)
	}

	return response, nil
}
