package agent

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func InitAgent(ctx context.Context, model model.BaseChatModel, tools []tool.BaseTool, handlers []adk.ChatModelAgentMiddleware) (*adk.ChatModelAgent, error) {

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MindAgent",
		Model:       model,
		Description: "An Agent will provide suggest for a fish tank",
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		Handlers: handlers,
	})
	return agent, err
}

func InitRunner(ctx context.Context, agent *adk.ChatModelAgent, tools []tool.InvokableTool) *adk.Runner {

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})
	return runner
}

func RunAgent(ctx context.Context, runner *adk.Runner, input []adk.Message) ([]*schema.Message, error) {
	var res []*schema.Message
	iter := runner.Run(ctx, input)
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return nil, event.Err
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			msg := event.Output.MessageOutput.Message
			if msg != nil && len(msg.ToolCalls) == 0 && msg.Content != "" {
				res = append(res, msg)
			}
		}
	}
	return res, nil
}
