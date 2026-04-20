package flows

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/claude"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func AquaGraph(ctx context.Context, chatModel *claude.ChatModel, tools []tool.InvokableTool) (compose.Runnable[[]*schema.Message, []*schema.Message], error) {

	var toolsInfo []*schema.ToolInfo
	for _, t := range tools {
		if info, err := t.Info(ctx); err == nil {
			toolsInfo = append(toolsInfo, info)
		} else {
			fmt.Printf("Error fetching tool info: %v\n", err)
		}
	}

	err := chatModel.BindTools(toolsInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to bind tools: %w", err)
	}

	g := compose.NewGraph[[]*schema.Message, []*schema.Message]()

	g.AddLambdaNode("LLMEngine", LLMNode(ctx, chatModel, nil))
	g.AddLambdaNode("ToolExecutor", ToolNode(tools))

	g.AddEdge(compose.START, "LLMEngine")
	g.AddBranch("LLMEngine", LoopBranch())
	g.AddEdge("ToolExecutor", "LLMEngine")

	return g.Compile(ctx)
}
