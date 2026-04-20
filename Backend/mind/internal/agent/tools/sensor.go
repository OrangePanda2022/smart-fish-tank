package tools

import (
	"context"
	"encoding/json"
	"log"
	"mind/internal/repo"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type SensorParams struct {
	TankID string `json:"tank_id" jsonschema:"description=需要鱼缸的唯一标识符 TankID"`
}

func NewSensorTool(ctx context.Context, sensorRepo *repo.SensorRepo) tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "SensorQueryTool",
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
		}, func(ctx context.Context, query *SensorParams) (string, error) {
			log.Println("sensor_tool 被调用")
			data, err := sensorRepo.GetLatestSensorData(query.TankID, 10)
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
