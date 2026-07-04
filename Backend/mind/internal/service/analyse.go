package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mind/internal/agent/prompts"
	"mind/internal/domain"
	"mind/internal/repo"
	"regexp"
	"strings"
	"time"

	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

type AnalyseService struct {
	aquaReActAgent *react.Agent
	frameRepo      repo.FrameRepository // nil 时降级为纯文本分析
}

func NewAnalyseService(aquaReActAgent *react.Agent, frameRepo repo.FrameRepository) *AnalyseService {
	return &AnalyseService{
		aquaReActAgent: aquaReActAgent,
		frameRepo:      frameRepo,
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

	if len(sensorData) > 0 {
		dataBytes, _ := json.Marshal(sensorData)
		initialMsgs = append(initialMsgs, schema.UserMessage(fmt.Sprintf("传感器数据: %s", string(dataBytes))))
	}

	// 多模态：拉取鱼缸最新帧并作为图片输入；无帧或出错则降级为纯文本（无回归）
	if s.frameRepo != nil {
		jpeg, ferr := s.frameRepo.GetLatestFrame(tankID)
		switch {
		case ferr != nil:
			log.Printf("获取鱼缸帧失败,降级纯文本: %v", ferr)
		case len(jpeg) > 0:
			b64 := base64.StdEncoding.EncodeToString(jpeg) // 裸 base64；ark 的 ensureDataURL 会拼成 data:URL，勿自带 data: 前缀
			imgMsg := &schema.Message{
				Role: schema.User,
				UserInputMultiContent: []schema.MessageInputPart{
					{Type: schema.ChatMessagePartTypeText, Text: "这是鱼缸当前的画面，请结合画面与传感器数据进行分析。"},
					{Type: schema.ChatMessagePartTypeImageURL, Image: &schema.MessageInputImage{
						MessagePartCommon: schema.MessagePartCommon{
							Base64Data: &b64,
							MIMEType:   "image/jpeg",
						},
					}},
				},
			}
			initialMsgs = append(initialMsgs, imgMsg)
		default:
			log.Printf("鱼缸 %s 暂无帧,降级纯文本分析", tankID)
		}
	}

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
	normalizeDecision(&decision)

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
		Zh: domain.DecisionContent{
			Summary:   summary,
			Actions:   []domain.Action{{Device: "系统", Action: "人工检查"}},
			Reasoning: "无法获取完整数据，建议人工检查设备状态。",
		},
		En: domain.DecisionContent{
			Summary:   "Report generation failed",
			Actions:   []domain.Action{{Device: "system", Action: "manual inspection"}},
			Reasoning: "Complete data could not be obtained. Please inspect the devices manually.",
		},
	}
}

func generateResponse(tankID string, decision *domain.Decision) *domain.AnalyseResponse {
	normalizeDecision(decision)
	return &domain.AnalyseResponse{
		ReportID: generateReportID(),
		TankID:   tankID,
		Decision: *decision,
	}
}

func normalizeDecision(decision *domain.Decision) {
	if decision == nil {
		return
	}

	legacy := domain.DecisionContent{
		Summary:   decision.Summary,
		Actions:   decision.Actions,
		Reasoning: decision.Reasoning,
	}
	if decision.Zh.Summary == "" {
		decision.Zh.Summary = legacy.Summary
	}
	if len(decision.Zh.Actions) == 0 {
		decision.Zh.Actions = legacy.Actions
	}
	if decision.Zh.Reasoning == "" {
		decision.Zh.Reasoning = legacy.Reasoning
	}
	if decision.En.Summary == "" {
		decision.En.Summary = decision.Zh.Summary
	}
	if len(decision.En.Actions) == 0 {
		decision.En.Actions = decision.Zh.Actions
	}
	if decision.En.Reasoning == "" {
		decision.En.Reasoning = decision.Zh.Reasoning
	}

	decision.Summary = decision.Zh.Summary
	decision.Actions = decision.Zh.Actions
	decision.Reasoning = decision.Zh.Reasoning
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
