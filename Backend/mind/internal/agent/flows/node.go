package flows

import (
	"context"
	"fmt"
	"mind/internal/agent/llm"

	"github.com/cloudwego/eino-ext/components/model/claude"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func LLMNode(ctx context.Context, model *claude.ChatModel, msgs []*schema.Message) *compose.Lambda {
	llmNode := compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		// 调用大模型
		resp, err := llm.GenerateAnswer(ctx, model, msgs)
		if err != nil {
			return nil, err
		}
		// 将大模型的回复追加到上下文
		return append(msgs, resp), nil
	})
	return llmNode
}

func ToolNode(tools []tool.InvokableTool) *compose.Lambda {
	toolExecNode := compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		lastMsg := msgs[len(msgs)-1]

		// 遍历大模型请求的所有 ToolCall
		for _, call := range lastMsg.ToolCalls {
			var matchedTool tool.InvokableTool
			for _, t := range tools {
				if tool, err := t.Info(ctx); err == nil && tool.Name == call.Function.Name {
					matchedTool = t
					break
				} else if err != nil {
					fmt.Printf("Error fetching tool info: %v\n", err)
				}
			}

			if matchedTool != nil {
				// 具体工具逻辑
				toolResult, err := matchedTool.InvokableRun(ctx, call.Function.Arguments)

				resultStr := toolResult
				if err != nil {
					// 将报错告诉 LLM
					resultStr = fmt.Sprintf("Error executing tool %s: %v", call.Function.Name, err)
				}

				// 将执行结果追加到历史
				toolMsg := schema.ToolMessage(resultStr, call.ID)
				msgs = append(msgs, toolMsg)
			}
		}
		return msgs, nil
	})
	return toolExecNode
}

func LoopBranch() *compose.GraphBranch {
	loopCondition := compose.NewGraphBranch(
		func(ctx context.Context, msgs []*schema.Message) (string, error) {
			lastMsg := msgs[len(msgs)-1]

			if len(lastMsg.ToolCalls) > 0 {
				return "ToolExecutor", nil
			}

			return compose.END, nil
		},
		map[string]bool{"ToolExecutor": true, compose.END: true})
	return loopCondition
}
