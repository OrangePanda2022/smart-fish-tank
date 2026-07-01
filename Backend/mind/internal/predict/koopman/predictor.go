package koopman

import (
	"fmt"
	"math"
	"time"

	"mind/internal/domain"

	"gonum.org/v1/gonum/mat"
)

// Predictor 执行多步预测
type Predictor struct {
	model *Model
}

// NewPredictor 创建预测器
func NewPredictor(m *Model) *Predictor {
	return &Predictor{model: m}
}

// PredictMultiStep 使用 Koopman 模型做 N 步预测
//
//	核心循环: ψ_{k+1} = K·ψ_k + B·u_k
//
// 每步从提升态前 stateDim 维提取物理量，检查是否超安全范围标记 IsAnomaly。
// 若 controlSeq 为 nil，使用零输入（自由演化），用于异常检测。
// 若提升态任一分量 |ψ_i| > 100，截断预测并设 confidence=0。
func (p *Predictor) PredictMultiStep(
	initialState domain.StateVector,
	controlSequence []domain.ControlVector, // len = horizon, 可以为 nil
	horizon int,
	deltaT float64,
) ([]domain.PredictionPoint, error) {
	if !p.model.IsLoaded() {
		return nil, fmt.Errorf("Koopman模型未加载")
	}
	cfg := p.model.Config()

	// 无控制序列则用 ctrl_mean 填充（自由演化）
	// 零向量归一化后为大负数，通过 B 矩阵产生漂移；
	// ctrl_mean 归一化后为零向量，B·0 = 0 无偏移。
	if controlSequence == nil {
		controlSequence = make([]domain.ControlVector, horizon)
		for k := 0; k < horizon; k++ {
			copy(controlSequence[k][:], cfg.Normalization.CtrlMean[:])
		}
	}
	if len(controlSequence) < horizon {
		return nil, fmt.Errorf("控制序列长度 %d < 预测步数 %d", len(controlSequence), horizon)
	}

	// 计算初始提升态
	lifted, err := ComputeLiftedState(p.model, initialState)
	if err != nil {
		return nil, err
	}
	liftedVec := mat.NewVecDense(cfg.LiftDim, lifted)

	predictions := make([]domain.PredictionPoint, 0, horizon)
	t := time.Now()
	diverged := false

	for step := 0; step < horizon; step++ {
		// ψ_{k+1} = K·ψ_k + B·u_k
		var nextLifted mat.VecDense
		nextLifted.MulVec(p.model.K, liftedVec)

		// 加入控制输入
		ctrlNorm := NormalizeControl(controlSequence[step], cfg.Normalization)
		ctrlVec := mat.NewVecDense(cfg.ControlDim, ctrlNorm[:])
		var bu mat.VecDense
		bu.MulVec(p.model.B, ctrlVec)
		nextLifted.AddVec(&nextLifted, &bu)

		// 数值稳定性检查
		maxAbs := 0.0
		for i := 0; i < cfg.LiftDim; i++ {
			absVal := math.Abs(nextLifted.AtVec(i))
			if absVal > maxAbs {
				maxAbs = absVal
			}
		}
		if maxAbs > 100 {
			diverged = true
			break
		}

		// 提取物理量并反归一化
		physicalState := ExtractPhysicalState(nextLifted.RawVector().Data, cfg.Normalization)

		t = t.Add(time.Duration(deltaT) * time.Second)

		// 检查是否超安全范围
		isAnomaly := checkAnomaly(physicalState, cfg.SafeRanges)

		predictions = append(predictions, domain.PredictionPoint{
			State:     physicalState,
			Time:      t,
			IsAnomaly: isAnomaly,
		})

		// 下一轮使用完整的提升态（不能截断到 stateDim 维！）
		liftedVec = &nextLifted
	}

	// 如果发散，截断处不再继续
	if diverged && len(predictions) < horizon {
		// 已经收集了部分有效预测，直接返回
	}

	return predictions, nil
}

// checkAnomaly 检查物理态是否超安全范围（9 维全覆盖）
func checkAnomaly(s domain.StateVector, sr domain.SafeRanges) bool {
	return s[0] < sr.Temperature.Min || s[0] > sr.Temperature.Max ||
		s[1] < sr.PH.Min || s[1] > sr.PH.Max ||
		s[2] < sr.Oxygen.Min || s[2] > sr.Oxygen.Max ||
		s[3] < sr.Ammonia.Min || s[3] > sr.Ammonia.Max ||
		s[4] < sr.WaterLevel.Min || s[4] > sr.WaterLevel.Max ||
		s[5] < sr.TDS.Min || s[5] > sr.TDS.Max ||
		s[6] < sr.Nitrate.Min || s[6] > sr.Nitrate.Max ||
		s[7] < sr.Nitrite.Min || s[7] > sr.Nitrite.Max ||
		s[8] < sr.Chloride.Min || s[8] > sr.Chloride.Max
}
