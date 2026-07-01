package mpc

import "mind/internal/domain"

// ControlBounds 返回控制输入的下界和上界，归一化到 [0,1]
func ControlBounds() (lb, ub domain.ControlVector) {
	return domain.ControlVector{0, 0, 0, 0, 0},
		domain.ControlVector{1, 1, 1, 1, 1}
}

// StateBounds 将安全范围转换为9维状态向量上下界
func StateBounds(sr domain.SafeRanges) (min, max domain.StateVector) {
	return domain.StateVector{sr.Temperature.Min, sr.PH.Min, sr.Oxygen.Min, sr.Ammonia.Min, sr.WaterLevel.Min, sr.TDS.Min, sr.Nitrate.Min, sr.Nitrite.Min, sr.Chloride.Min},
		domain.StateVector{sr.Temperature.Max, sr.PH.Max, sr.Oxygen.Max, sr.Ammonia.Max, sr.WaterLevel.Max, sr.TDS.Max, sr.Nitrate.Max, sr.Nitrite.Max, sr.Chloride.Max}
}

// EmergencyThresholds 定义9维状态超出安全范围的紧急阈值
// 当 QP 在安全约束下不可行时，放宽到这些阈值重新求解
func EmergencyThresholds() (min, max domain.StateVector) {
	return domain.StateVector{18, 5.5, 3, 1.0, 40, 50, 0, 0, 0},
		domain.StateVector{33, 9.5, 12, 1.0, 98, 800, 100, 2.0, 500}
}
