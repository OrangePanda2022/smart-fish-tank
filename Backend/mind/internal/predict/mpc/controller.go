package mpc

import (
	"fmt"
	"math"
	"time"

	"mind/internal/domain"
	"mind/internal/predict/koopman"

	"gonum.org/v1/gonum/mat"
)

// Controller 封装 MPC 控制器
type Controller struct {
	model     *koopman.Model
	solver    *QPSolver
	cost      *CostMatrices
	horizon   int
	predictor *koopman.Predictor
}

// NewController 创建 MPC 控制器
func NewController(m *koopman.Model, horizon int, stateWeight, ctrlWeight float64) *Controller {
	return &Controller{
		model:     m,
		solver:    NewQPSolver(),
		cost:      NewCostMatrices(stateWeight, ctrlWeight),
		horizon:   horizon,
		predictor: koopman.NewPredictor(m),
	}
}

// Solve 计算最优控制并返回完整的 PredictResponse
//
// 1. 计算设定点
// 2. 构建 + 求解 QP（或使用启发式控制，若模型未训练）
// 3. 提取最优控制序列
// 4. 用最优控制做多步预测
// 5. 计算总代价
// 6. 将第一步控制向量转为语义动作
// 7. 同时跑自由演化预测
// 8. 组装 PredictResponse
func (c *Controller) Solve(
	currentState domain.StateVector,
	safeRanges domain.SafeRanges,
) (*domain.PredictResponse, error) {
	if !c.model.IsLoaded() {
		return nil, fmt.Errorf("Koopman模型未加载，无法执行MPC")
	}

	cfg := c.model.Config()
	setpoint := ComputeSetpoint(safeRanges)

	// 未训练模型（K≈I）→ QP 无法产出有效控制，改用启发式
	useHeuristic := c.model.IsUntrained()

	var controlSeq []domain.ControlVector
	var mpcPredicted []domain.PredictionPoint
	var totalCost float64

	if useHeuristic {
		// 启发式控制：根据状态偏离安全范围程度计算控制强度
		heuristicU := heuristicControl(currentState, safeRanges)
		controlSeq = make([]domain.ControlVector, c.horizon)
		for k := 0; k < c.horizon; k++ {
			controlSeq[k] = heuristicU
		}

		// 用启发式控制做多步预测
		var err error
		mpcPredicted, err = c.predictor.PredictMultiStep(
			currentState, controlSeq, c.horizon, cfg.DeltaT)
		if err != nil {
			return nil, fmt.Errorf("启发式预测轨迹失败: %w", err)
		}
		totalCost = 0 // 启发式无代价值
	} else {
		// 构建 + 求解 QP
		H, f, lb, ub, err := BuildQP(c.model, c.cost, currentState, setpoint, c.horizon)
		if err != nil {
			return nil, fmt.Errorf("构建QP失败: %w", err)
		}

		Uopt, err := c.solver.Solve(H, f, lb, ub)
		if err != nil {
			return nil, fmt.Errorf("求解QP失败: %w", err)
		}

		// 提取控制序列
		ctrlDim := c.model.Config().ControlDim
		controlSeq = make([]domain.ControlVector, c.horizon)
		for k := 0; k < c.horizon; k++ {
			for i := 0; i < ctrlDim; i++ {
				controlSeq[k][i] = Uopt[k*ctrlDim+i]
			}
		}

		// 用最优控制做多步预测
		var err2 error
		mpcPredicted, err2 = c.predictor.PredictMultiStep(
			currentState, controlSeq, c.horizon, cfg.DeltaT)
		if err2 != nil {
			return nil, fmt.Errorf("MPC预测轨迹失败: %w", err2)
		}

		// 计算总代价: cost = 0.5·UᵀHU + fᵀU
		uVec := mat.NewVecDense(len(Uopt), Uopt)
		var hu mat.VecDense
		hu.MulVec(H, uVec)
		totalCost = 0.5*dot(Uopt, hu.RawVector().Data) + dot(Uopt, f.RawVector().Data)
	}

	_ = totalCost // 代价值可用于日志/监控，暂不写入响应

	// 将第一步控制向量转为语义动作
	firstCtrl := controlSeq[0]
	actions := controlToActions(firstCtrl)

	// 自由演化预测（不做控制 = 设备处于默认"关"状态）
	// 不能传 nil，否则 NormalizeControl 将物理零向量归一化为大负数，
	// 通过 B 矩阵产生漂移。自由演化应使用 ctrl_mean（归一化后为零向量），
	// 表示"设备不施加额外控制"，B·0 = 0 无偏移。
	freeCtrl := make([]domain.ControlVector, c.horizon)
	for k := 0; k < c.horizon; k++ {
		copy(freeCtrl[k][:], cfg.Normalization.CtrlMean[:])
	}
	freePredicted, err := c.predictor.PredictMultiStep(
		currentState, freeCtrl, c.horizon, cfg.DeltaT)
	if err != nil {
		return nil, fmt.Errorf("自由演化预测失败: %w", err)
	}

	// 寻找自由演化中首次异常步
	anomalyStep := -1
	anomalyWarn := ""
	for k, pred := range freePredicted {
		if pred.IsAnomaly {
			anomalyStep = k
			anomalyWarn = fmt.Sprintf("第 %d 步预测异常 (%.0f秒后)，建议立即采取行动",
				k+1, float64(k+1)*cfg.DeltaT)
			break
		}
	}

	// 估算置信度：基于预测偏差
	confidence := estimateConfidence(mpcPredicted, setpoint, safeRanges)

	return &domain.PredictResponse{
		Predictions:  freePredicted,
		MPCActions:   actions,
		MPCPredicted: mpcPredicted,
		AnomalyStep:  anomalyStep,
		AnomalyWarn:  anomalyWarn,
		Confidence:   confidence,
		GeneratedAt:  time.Now(),
	}, nil
}

// estimateConfidence 基于MPC轨迹偏离设定点的程度估算置信度
func estimateConfidence(predicted []domain.PredictionPoint, setpoint domain.StateVector, sr domain.SafeRanges) float64 {
	confidence := 1.0
	rangeWidths := []float64{
		sr.Temperature.Max - sr.Temperature.Min,
		sr.PH.Max - sr.PH.Min,
		sr.Oxygen.Max - sr.Oxygen.Min,
		sr.Ammonia.Max - sr.Ammonia.Min,
		sr.WaterLevel.Max - sr.WaterLevel.Min,
		sr.TDS.Max - sr.TDS.Min,
		sr.Nitrate.Max - sr.Nitrate.Min,
		sr.Nitrite.Max - sr.Nitrite.Min,
		sr.Chloride.Max - sr.Chloride.Min,
	}

	for _, pred := range predicted {
		for i := 0; i < len(rangeWidths) && i < len(pred.State) && i < len(setpoint); i++ {
			dev := math.Abs(pred.State[i] - setpoint[i])
			if rangeWidths[i] > 0 && dev/rangeWidths[i] > 2 {
				confidence *= 0.9
			}
		}
	}
	if confidence < 0.1 {
		confidence = 0.1
	}
	return confidence
}

// controlToActions 将归一化控制向量转换为语义设备动作
//
// 启发式控制下偏离安全范围的参数不会恰好落在 0.5，因此阈值需区分：
//   - < 0.35: decrease（含偏离上限的情况，如温度过高 heater 值 ~0.35）
//   - 0.35 - 0.65: maintain（仅在安全区中心附近）
//   - > 0.65: increase（含偏离下限的情况，如氧含量偏低 aerator 值 ~0.64）
//
// 这样安全区内的绝大部分映射为 maintain，仅边界附近和越界值才触发纠正动作
func controlToActions(u domain.ControlVector) []domain.ControlAction {
	var actions []domain.ControlAction
	for i, val := range u {
		mapping := domain.DeviceActionMap[i]
		if val < 0.35 {
			actions = append(actions, domain.ControlAction{
				Device: mapping.Device,
				Action: mapping.Decr,
				Value:  val,
			})
		} else if val > 0.65 {
			actions = append(actions, domain.ControlAction{
				Device: mapping.Device,
				Action: mapping.Incr,
				Value:  val,
			})
		} else {
			actions = append(actions, domain.ControlAction{
				Device: mapping.Device,
				Action: "maintain",
				Value:  val,
			})
		}
	}
	return actions
}

// heuristicControl 未训练模型的启发式控制——根据状态偏离安全范围计算控制强度
//
// 对每个设备，计算其主控状态量偏离安全范围的程度：
//   - 安全区内 → 0.3~0.7（维持），中心=0.5
//   - 低于安全下界 → 0.5 + 0.5·ratio^0.7，ratio=偏离/到紧急线距离
//   - 高于安全上界 → 0.5 - 0.5·ratio^0.7
//
// ratio=0.5（偏离到紧急线一半）时，强度≈0.5+0.5·0.5^0.7≈0.81→约80%
func heuristicControl(currentState domain.StateVector, safeRanges domain.SafeRanges) domain.ControlVector {
	var u domain.ControlVector
	emergMin, emergMax := EmergencyThresholds()

	// heater → temperature (index 0)
	u[0] = deviceControlIntensity(currentState[0],
		safeRanges.Temperature.Min, safeRanges.Temperature.Max,
		emergMin[0], emergMax[0])

	// aerator → oxygen (index 2)
	u[1] = deviceControlIntensity(currentState[2],
		safeRanges.Oxygen.Min, safeRanges.Oxygen.Max,
		emergMin[2], emergMax[2])

	// pump → water_level (index 4)，同时兼顾氨氮稀释
	pumpWL := deviceControlIntensity(currentState[4],
		safeRanges.WaterLevel.Min, safeRanges.WaterLevel.Max,
		emergMin[4], emergMax[4])
	// 氨氮过高时也需要换水稀释
	pumpNH3 := 0.0
	if currentState[3] > safeRanges.Ammonia.Max {
		ratio := (currentState[3] - safeRanges.Ammonia.Max) / (emergMax[3] - safeRanges.Ammonia.Max)
		ratio = math.Min(ratio, 1.0)
		pumpNH3 = 0.5 + 0.5*math.Pow(ratio, 0.7)
	}
	u[2] = math.Max(pumpWL, pumpNH3)

	// feeder → 氨氮的反向控制：氨氮高→不喂，氨氮低→适量喂
	if currentState[3] >= safeRanges.Ammonia.Max {
		// 氨氮超标 → 停止投喂
		u[3] = math.Max(0.05, 0.3-0.3*(currentState[3]-safeRanges.Ammonia.Max)/(emergMax[3]-safeRanges.Ammonia.Max))
	} else {
		// 氨氮正常/偏低 → 适量投喂，越低越多
		u[3] = 0.5 + 0.3*(safeRanges.Ammonia.Max-currentState[3])/safeRanges.Ammonia.Max
	}

	// light → 次要，维持为主；温度过高时适当调暗辅助降温
	u[4] = 0.5
	if currentState[0] > safeRanges.Temperature.Max {
		ratio := (currentState[0] - safeRanges.Temperature.Max) / (emergMax[0] - safeRanges.Temperature.Max)
		ratio = math.Min(ratio, 1.0)
		u[4] = 0.5 - 0.3*math.Pow(ratio, 0.7)
	}

	return u
}

// deviceControlIntensity 计算单个设备的控制强度 [0,1]
//
//	value 在安全区内 → 0.3~0.7（安全区中心=0.5）
//	value 低于安全下界 → 0.5 + 0.5·ratio^0.7（加强控制）
//	value 高于安全上界 → 0.5 - 0.5·ratio^0.7（减弱控制）
func deviceControlIntensity(value, safeMin, safeMax, emergMin, emergMax float64) float64 {
	// 在安全范围内
	if value >= safeMin && value <= safeMax {
		rangeWidth := safeMax - safeMin
		if rangeWidth <= 0 {
			return 0.5
		}
		normalizedPos := (value - safeMin) / rangeWidth // 0=下界, 1=上界
		return 0.3 + 0.4*normalizedPos                 // 0.3~0.7
	}

	// 低于安全下界 → 需要加大控制
	if value < safeMin {
		if emergMin >= safeMin {
			return 1.0 // 无紧急线，全开
		}
		ratio := (safeMin - value) / (safeMin - emergMin)
		ratio = math.Min(ratio, 1.0)
		return 0.5 + 0.5*math.Pow(ratio, 0.7)
	}

	// 高于安全上界 → 需要减少控制
	if emergMax <= safeMax {
		return 0.0
	}
	ratio := (value - safeMax) / (emergMax - safeMax)
	ratio = math.Min(ratio, 1.0)
	return 0.5 - 0.5*math.Pow(ratio, 0.7)
}
