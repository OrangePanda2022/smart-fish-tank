package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mind/internal/repo"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type TankParams struct {
	TankID string `json:"tank_id" jsonschema:"description=需要鱼缸的唯一标识符 TankID"`
}

func NewTankTool(ctx context.Context, tankRepo repo.TankRepository) tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "TankQueryTool",
			Desc: "基于输入的 TankID，在后端数据库中进行检索，返回相关的数据",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"tank_id": {
						Type:     schema.String,
						Desc:     "需要鱼缸的唯一标识符 TankID",
						Required: true,
					},
				},
			),
		}, func(ctx context.Context, query *TankParams) (string, error) {
			log.Println("tank_tool 被调用")
			if tankRepo == nil {
				return "", fmt.Errorf("tank repo 未初始化（NATS 连接失败）")
			}
			data, err := tankRepo.GetTank(query.TankID)
			if err != nil {
				return "", err
			}
			result, err := json.Marshal(data)
			if err != nil {
				return "", err
			}
			return string(result), nil
		},
	)
}
