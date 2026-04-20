package react

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

func ReactAgent(ctx context.Context, model model.ToolCallingChatModel, tools []tool.BaseTool) (*react.Agent, error) {

	raAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: model,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})
	return raAgent, err
}
