package koopman

import (
	"fmt"

	"mind/internal/domain"
)

// ComputeLiftedState 对原始状态做 z-score 归一化后提升到可观测量空间
func ComputeLiftedState(m *Model, raw domain.StateVector) ([]float64, error) {
	if !m.IsLoaded() {
		return nil, fmt.Errorf("Koopman模型未加载")
	}
	cfg := m.Config()
	normalized := normalizeState(raw, cfg.Normalization)
	return m.Lift(normalized), nil
}

// normalizeState 使用训练统计量做 z-score 归一化
func normalizeState(s domain.StateVector, norm domain.NormConfig) domain.StateVector {
	var out domain.StateVector
	stateDim := len(norm.StateMean)
	for i := 0; i < stateDim && i < len(out); i++ {
		if norm.StateStd[i] > 0 {
			out[i] = (s[i] - norm.StateMean[i]) / norm.StateStd[i]
		} else {
			out[i] = s[i] - norm.StateMean[i]
		}
	}
	return out
}

// denormalizeState 反 z-score 归一化
func denormalizeState(s domain.StateVector, norm domain.NormConfig) domain.StateVector {
	var out domain.StateVector
	stateDim := len(norm.StateMean)
	for i := 0; i < stateDim && i < len(out); i++ {
		out[i] = s[i]*norm.StateStd[i] + norm.StateMean[i]
	}
	return out
}

// ExtractPhysicalState 从提升态前 stateDim 维提取并反归一化物理量
func ExtractPhysicalState(lifted []float64, norm domain.NormConfig) domain.StateVector {
	stateDim := len(norm.StateMean)
	var normalized domain.StateVector
	copy(normalized[:], lifted[:stateDim])
	return denormalizeState(normalized, norm)
}

// NormalizeControl 归一化控制向量
func NormalizeControl(u domain.ControlVector, norm domain.NormConfig) domain.ControlVector {
	var out domain.ControlVector
	for i := 0; i < 5; i++ {
		if norm.CtrlStd[i] > 0 {
			out[i] = (u[i] - norm.CtrlMean[i]) / norm.CtrlStd[i]
		} else {
			out[i] = u[i] - norm.CtrlMean[i]
		}
	}
	return out
}

// SensorDataToState 将领域 SensorData 转换为 9D StateVector
func SensorDataToState(data *domain.SensorData) domain.StateVector {
	return domain.StateVector{
		data.Temperature,
		data.PH,
		data.Oxygen,
		data.Ammonia,
		data.WaterLevel,
		data.TDS,
		data.Nitrate,
		data.Nitrite,
		data.Chloride,
	}
}
