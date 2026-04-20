package moe

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func NewExpert(ctx context.Context, raAgent *react.Agent) *host.Specialist {
	return &host.Specialist{
		AgentMeta: host.AgentMeta{
			Name:        "analyse_expert",
			IntendedUse: "Analyse the data of the fish tank",
		},
		Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (output *schema.Message, err error) {
			return raAgent.Generate(ctx, input, opts...)
		},
	}
}

func NewHost(ctx context.Context, model model.ToolCallingChatModel) (*host.Host, error) {
	return &host.Host{
		ToolCallingModel: model,
		SystemPrompt:     "你可以分析鱼缸状态并且给出操作建议，当用户提问你时，你可以调用专家来帮你分析数据，专家会给你分析结果，你需要根据分析结果给用户一个清晰的分析报告，报告需要包含以下内容：1. 当前鱼缸的状态 2. 可能存在的问题 3. 推荐的操作建议",
	}, nil
}

func NewSummarizer(ctx context.Context, model model.BaseChatModel) *host.Summarizer {
	return &host.Summarizer{
		ChatModel:    model,
		SystemPrompt: "总结其他专家的输出",
	}
}

func NewMOEAgent(ctx context.Context, h host.Host, experts []*host.Specialist, summarizer *host.Summarizer) (*host.MultiAgent, error) {
	hostMOEAgent, err := host.NewMultiAgent(ctx, &host.MultiAgentConfig{
		Host:        h,
		Specialists: experts,
		Summarizer:  summarizer,
	})
	return hostMOEAgent, err
}
