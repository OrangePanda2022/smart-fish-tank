package service

import (
	"context"
	"encoding/json"
	"fmt"
	"mind/internal/agent/prompts"
	"mind/internal/domain"
	"regexp"
	"strings"
	"time"

	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

type AnalyseService struct {
	aquaReActAgent *react.Agent
}

func NewAnalyseService(aquaReActAgent *react.Agent) *AnalyseService {
	return &AnalyseService{
		aquaReActAgent: aquaReActAgent,
	}
}

func (s *AnalyseService) Analyse(ctx context.Context, tankID string, overrideData []*domain.SensorData) (*domain.AnalyseResponse, error) {

	report, err := s.generateReport(ctx, tankID, overrideData)
	if err != nil {
		fmt.Println(err)
		report = s.generateDefaultReport(tankID, "生成报告失败")
	}

	response := generateResponse(tankID, report)

	return response, nil
}

func (s *AnalyseService) AnalyseStream(ctx context.Context, tankID string, overrideData []*domain.SensorData) (*schema.StreamReader[*schema.Message], error) {
	return s.generateReportStream(ctx, tankID, overrideData)
}

func (s *AnalyseService) generateReport(ctx context.Context, tankID string, sensorData []*domain.SensorData) (*domain.Decision, error) {

	initialMsgs := []*schema.Message{
		prompts.SysPrompt,
		prompts.UserPrompt,
		schema.UserMessage(fmt.Sprintf("鱼缸 ID 为 %s", tankID)),
	}

	// if len(sensorData) > 0 {
	// 	dataBytes, _ := json.Marshal(sensorData)
	// 	initialMsgs = append(initialMsgs, schema.UserMessage(fmt.Sprintf("传感器数据: %s", string(dataBytes))))
	// }

	url := "https://aqua.cn-nb1.rains3.com/VID_20260416_161729.mp4"

	video := schema.MessageInputPart{
		Type: schema.ChatMessagePartTypeVideoURL,
		Video: &schema.MessageInputVideo{
			MessagePartCommon: schema.MessagePartCommon{
				URL:      &url,
				MIMEType: "video/mp4",
			},
		},
	}

	input := schema.Message{
		Role: schema.User,
		UserInputMultiContent: []schema.MessageInputPart{
			video,
		},
	}

	initialMsgs = append(initialMsgs, &input)

	// 这里直接调用 ReAct MOE 来生成分析报告
	// finalMsg, _ := s.aquaMoE.Generate(ctx, initialMsgs)

	// 这里直接调用 ReAct Agent 来生成分析报告
	finalMsg, err := s.aquaReActAgent.Generate(ctx, initialMsgs)

	if err != nil {
		return nil, err
	}

	// 这里直接调用 Agent 来生成分析报告
	// finalMsgs, _ := agent.RunAgent(ctx, s.aquaRunner, initialMsgs)

	// 这里直接调用 Graph 来生成分析报告
	// finalMsgs, err := s.aquaGraph.Invoke(ctx, initialMsgs)

	// if err != nil {
	// 	return nil, fmt.Errorf("Agent 执行失败: %w", err)
	// }

	// rawOutput := finalMsgs[len(finalMsgs)-1].Content

	rawOutput := finalMsg.Content

	cleanJSON := extractRawJSON(rawOutput)

	var decision domain.Decision
	if err := json.Unmarshal([]byte(cleanJSON), &decision); err != nil {
		return nil, fmt.Errorf("大模型未按要求输出 JSON, 解析失败: %w, 原始输出: %s", err, rawOutput)
	}

	return &decision, nil
}

func (s *AnalyseService) generateReportStream(ctx context.Context, tankID string, sensorData []*domain.SensorData) (*schema.StreamReader[*schema.Message], error) {

	initialMsgs := []*schema.Message{
		prompts.UserPrompt,
		schema.UserMessage(fmt.Sprintf("鱼缸 ID 为 %s", tankID)),
	}

	url := "https://aqua.cn-nb1.rains3.com/VID_20260416_161729.mp4"

	video := schema.MessageInputPart{
		Type: schema.ChatMessagePartTypeVideoURL,
		Video: &schema.MessageInputVideo{
			MessagePartCommon: schema.MessagePartCommon{
				URL:      &url,
				MIMEType: "video/mp4",
			},
		},
	}

	input := schema.Message{
		Role: schema.User,
		UserInputMultiContent: []schema.MessageInputPart{
			video,
		},
	}

	initialMsgs = append(initialMsgs, &input)

	return s.aquaReActAgent.Stream(ctx, initialMsgs)
}

func (s *AnalyseService) generateDefaultReport(tankID, summary string) *domain.Decision {
	return &domain.Decision{
		StatusScore: 50,
		Summary:     summary,
		Actions:     []domain.Action{{Device: "system", Action: "check"}},
		Reasoning:   "无法获取完整数据，建议人工检查设备状态。",
	}
}

func generateResponse(tankID string, decision *domain.Decision) *domain.AnalyseResponse {
	return &domain.AnalyseResponse{
		ReportID: generateReportID(),
		TankID:   tankID,
		Decision: *decision,
	}
}

func generateReportID() string {
	return fmt.Sprintf("REP_%d", time.Now().Unix())
}

func extractRawJSON(raw string) string {
	content := strings.TrimSpace(raw)

	re := regexp.MustCompile(`(?s)\x60\x60\x60(?:json)?\s*(.*?)\s*\x60\x60\x60`)
	matches := re.FindStringSubmatch(content)

	var jsonStr string
	if len(matches) > 1 {
		jsonStr = strings.TrimSpace(matches[1])
	} else {
		start := strings.Index(content, "{")
		end := strings.LastIndex(content, "}")

		if start != -1 && end != -1 && end > start {
			jsonStr = content[start : end+1]
		} else {
			jsonStr = content
		}
	}

	return jsonStr
}
