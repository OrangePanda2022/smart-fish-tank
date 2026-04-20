package middleware

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

type safeToolMiddleware struct {
	*adk.BaseChatModelAgentMiddleware
}

func (m *safeToolMiddleware) WrapInvokableToolCall(ctx context.Context, endpoint adk.InvokableToolCallEndpoint, _ *adk.ToolContext) (adk.InvokableToolCallEndpoint, error) {
	return func(ctx context.Context, args string, opts ...tool.Option) (string, error) {
		result, err := endpoint(ctx, args, opts...)
		if err != nil {
			// 中断错误不转换，需要继续传播
			if _, ok := compose.IsInterruptRerunError(err); ok {
				return "", err
			}
			// 将错误转换为字符串，而不是返回错误
			return fmt.Sprintf("Error executing tool %v", err), nil
		}
		return result, nil
	}, nil
}
