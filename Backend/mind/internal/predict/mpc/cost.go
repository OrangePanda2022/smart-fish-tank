package mpc

import (
	"mind/internal/domain"

	"gonum.org/v1/gonum/mat"
)

// CostMatrices 保存 MPC 代价函数的权重矩阵
type CostMatrices struct {
	Q  *mat.DiagDense // 状态偏差代价 (9x9)
	QN *mat.DiagDense // 终端状态代价 (9x9)
	R  *mat.DiagDense // 控制努力代价 (5x5)
}

// NewCostMatrices 创建代价矩阵
// stateWeight 越大 → 越积极纠正偏差
// ctrlWeight 越大 → 控制动作越平滑保守
func NewCostMatrices(stateWeight, ctrlWeight float64) *CostMatrices {
	// 有毒物质（ammonia, nitrite）权重加倍
	q := mat.NewDiagDense(9, []float64{
		stateWeight,     // temperature
		stateWeight,     // ph
		stateWeight,     // oxygen
		stateWeight * 2, // ammonia — 有毒，优先
		stateWeight,     // water_level
		stateWeight,     // tds
		stateWeight,     // nitrate
		stateWeight * 2, // nitrite — 有毒，优先
		stateWeight,     // chloride
	})
	// 终端代价更重，确保稳定
	qn := mat.NewDiagDense(9, []float64{
		stateWeight * 2,
		stateWeight * 2,
		stateWeight * 2,
		stateWeight * 4, // 氨氮终端权重
		stateWeight * 2,
		stateWeight * 2,
		stateWeight * 2,
		stateWeight * 4, // 亚硝酸盐终端权重
		stateWeight * 2,
	})
	// feeder 代价减半——鼓励适量喂食
	r := mat.NewDiagDense(5, []float64{
		ctrlWeight,       // heater
		ctrlWeight,       // aerator
		ctrlWeight,       // pump
		ctrlWeight * 0.5, // feeder — 降低代价鼓励喂食
		ctrlWeight,       // light
	})
	return &CostMatrices{Q: q, QN: qn, R: r}
}

// ComputeSetpoint 计算各安全范围中点作为9维设定点
// 有毒物质（ammonia, nitrite）特殊：取下界（目标是尽量降低有毒物质）
func ComputeSetpoint(sr domain.SafeRanges) domain.StateVector {
	return domain.StateVector{
		(sr.Temperature.Min + sr.Temperature.Max) / 2, // 25.0
		(sr.PH.Min + sr.PH.Max) / 2,                   // 7.25
		(sr.Oxygen.Min + sr.Oxygen.Max) / 2,            // 7.0
		sr.Ammonia.Min,                                 // 0.0 — 有毒物质目标最小
		(sr.WaterLevel.Min + sr.WaterLevel.Max) / 2,    // 75.0
		(sr.TDS.Min + sr.TDS.Max) / 2,                  // 300.0
		(sr.Nitrate.Min + sr.Nitrate.Max) / 2,           // 25.0
		sr.Nitrite.Min,                                  // 0.0 — 有毒物质目标最小
		(sr.Chloride.Min + sr.Chloride.Max) / 2,         // 125.0
	}
}
