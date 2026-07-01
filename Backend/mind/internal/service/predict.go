package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"mind/internal/domain"
	"mind/internal/predict/koopman"
	"mind/internal/predict/mpc"
	"mind/internal/repo"
)

// PredictService 预测+MPC 业务逻辑层
type PredictService struct {
	model      *koopman.Model
	predictor  *koopman.Predictor
	ctrl       *mpc.Controller
	sensorRepo repo.SensorRepository

	// NATS 实时状态缓存（保存完整 SensorData 供扩展异常检测）
	tankData map[string]*domain.SensorData
	mu       sync.RWMutex

	// 预测精度追踪
	predCount  int
	totalError float64
}

// NewPredictService 创建预测服务
func NewPredictService(
	m *koopman.Model,
	sensorRepo repo.SensorRepository,
	horizon int,
	stateWeight, ctrlWeight float64,
) *PredictService {
	return &PredictService{
		model:      m,
		predictor:  koopman.NewPredictor(m),
		ctrl:       mpc.NewController(m, horizon, stateWeight, ctrlWeight),
		sensorRepo: sensorRepo,
		tankData:   make(map[string]*domain.SensorData),
	}
}

// Predict 执行预测 + MPC，返回完整结果
// 调用方只需提供 tankID 和 horizon
func (s *PredictService) Predict(ctx context.Context, tankID string, horizon int) (*domain.PredictResponse, error) {
	if !s.model.IsLoaded() {
		return nil, fmt.Errorf("预测模型未加载")
	}

	// 获取当前状态：优先从 NATS 缓存取，否则从 sensorRepo 查
	currentState, latestData, err := s.resolveCurrentState(ctx, tankID)
	if err != nil {
		return nil, fmt.Errorf("获取当前状态失败: %w", err)
	}

	safeRanges := s.resolveSafeRanges()

	resp, err := s.ctrl.Solve(currentState, safeRanges)
	if err != nil {
		return nil, fmt.Errorf("MPC求解失败: %w", err)
	}
	resp.TankID = tankID

	// 扩展参数异常检测（检查当前传感器读数是否超安全范围）
	if latestData != nil {
		extAnomalies := checkExtendedAnomalies(latestData, safeRanges)
		if len(extAnomalies) > 0 {
			resp.ExtendedAnomalies = extAnomalies
			if resp.AnomalyStep == -1 {
				resp.AnomalyStep = 0 // 标记存在当前异常
				warnings := make([]string, 0, len(extAnomalies))
				for _, w := range extAnomalies {
					warnings = append(warnings, w)
				}
				resp.AnomalyWarn = "当前传感器读数存在扩展参数异常: " + strings.Join(warnings, "; ")
			}
		}
	}

	return resp, nil
}

// UpdateState 更新 NATS 实时状态缓存（保存完整 SensorData）
func (s *PredictService) UpdateState(data *domain.SensorData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tankData[data.TankID] = data
}

// UpdateAccuracy 追踪预测精度
func (s *PredictService) UpdateAccuracy(predicted, actual domain.StateVector) {
	var sumSq float64
	stateDim := len(predicted)
	for i := 0; i < stateDim && i < len(actual); i++ {
		diff := predicted[i] - actual[i]
		sumSq += diff * diff
	}
	s.predCount++
	s.totalError += sumSq
}

// ModelStatus 返回模型加载状态
func (s *PredictService) ModelStatus() (loaded bool, version string, trainedAt string) {
	if !s.model.IsLoaded() {
		return false, "", ""
	}
	cfg := s.model.Config()
	return true, cfg.Version, cfg.TrainedAt
}

// resolveCurrentState 获取当前鱼缸状态，同时返回原始 SensorData 供扩展异常检测
func (s *PredictService) resolveCurrentState(ctx context.Context, tankID string) (domain.StateVector, *domain.SensorData, error) {
	// 优先从 NATS 缓存获取
	s.mu.RLock()
	if data, ok := s.tankData[tankID]; ok {
		s.mu.RUnlock()
		state := koopman.SensorDataToState(data)
		return state, data, nil
	}
	s.mu.RUnlock()

	// SHIT2 保留 nil 守卫路线（未采用副本的内存 mock 降级）：NATS 不可用时 sensorRepo 为 nil
	if s.sensorRepo == nil {
		return domain.StateVector{}, nil, fmt.Errorf("sensor repo 不可用（NATS 连接失败）")
	}

	// 从 sensorRepo 获取最新数据
	dataList, err := s.sensorRepo.GetLatestSensorData(tankID, 1)
	if err != nil {
		return domain.StateVector{}, nil, fmt.Errorf("传感器数据查询失败: %w", err)
	}
	if len(dataList) == 0 {
		return domain.StateVector{}, nil, fmt.Errorf("鱼缸 %s 无传感器数据", tankID)
	}
	latest := dataList[0]
	state := koopman.SensorDataToState(latest)

	// 更新缓存
	s.UpdateState(latest)
	return state, latest, nil
}

// resolveSafeRanges 返回安全范围（优先用模型训练时的配置）
func (s *PredictService) resolveSafeRanges() domain.SafeRanges {
	if s.model.IsLoaded() {
		return s.model.Config().SafeRanges
	}
	return domain.DefaultSafeRanges()
}

// estimateConfidence 返回启发式置信度
// 置信度由两部分组成：
//   - 基础分：预测次数不足时降低（数据不够，统计不可靠）
//   - 误差分：历史预测误差越小置信度越高
func (s *PredictService) estimateConfidence() float64 {
	// 基础分：5次以下0.5，10次以上1.0，中间线性插值
	base := 0.5
	if s.predCount >= 10 {
		base = 1.0
	} else if s.predCount >= 5 {
		base = 0.5 + 0.5*float64(s.predCount-5)/5.0
	}
	// 误差分：无历史数据时取1.0（不额外惩罚），有数据时用sigmoid
	errorFactor := 1.0
	if s.predCount > 0 {
		avgError := s.totalError / float64(s.predCount)
		errorFactor = 1.0 / (1.0 + avgError)
	}
	return base * errorFactor
}

// predictTime 计算预测时间戳序列
func predictTime(horizon int, deltaT float64) []time.Time {
	times := make([]time.Time, horizon)
	t := time.Now()
	for k := 0; k < horizon; k++ {
		t = t.Add(time.Duration(deltaT) * time.Second)
		times[k] = t
	}
	return times
}

// checkExtendedAnomalies 检查扩展参数是否超出安全范围
func checkExtendedAnomalies(data *domain.SensorData, sr domain.SafeRanges) map[string]string {
	anomalies := make(map[string]string)
	if data.TDS < sr.TDS.Min || data.TDS > sr.TDS.Max {
		anomalies["tds"] = fmt.Sprintf("TDS异常 (%.0f ppm, 安全线 %.0f-%.0f)", data.TDS, sr.TDS.Min, sr.TDS.Max)
	}
	if data.Nitrate < sr.Nitrate.Min || data.Nitrate > sr.Nitrate.Max {
		anomalies["nitrate"] = fmt.Sprintf("硝酸根异常 (%.1f mg/L, 安全线 %.0f-%.0f)", data.Nitrate, sr.Nitrate.Min, sr.Nitrate.Max)
	}
	if data.Nitrite < sr.Nitrite.Min || data.Nitrite > sr.Nitrite.Max {
		anomalies["nitrite"] = fmt.Sprintf("亚硝酸盐异常 (%.2f mg/L, 安全线 %.1f-%.1f)", data.Nitrite, sr.Nitrite.Min, sr.Nitrite.Max)
	}
	if data.Chloride < sr.Chloride.Min || data.Chloride > sr.Chloride.Max {
		anomalies["chloride"] = fmt.Sprintf("氯离子异常 (%.1f mg/L, 安全线 %.0f-%.0f)", data.Chloride, sr.Chloride.Min, sr.Chloride.Max)
	}
	return anomalies
}
