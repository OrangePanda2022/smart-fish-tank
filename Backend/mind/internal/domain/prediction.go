package domain

import "time"

// StateVector 是9维鱼缸状态，顺序: [temperature, ph, oxygen, ammonia, water_level, tds, nitrate, nitrite, chloride]
type StateVector [9]float64

// ControlVector 是5维归一化控制输入，顺序: [heater, aerator, pump, feeder, light]
// 每个分量归一化到 [0, 1]
type ControlVector [5]float64

// SafeRange 定义单个参数的安全范围
type SafeRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// SafeRanges 5个核心参数 + 4个扩展监测参数的安全范围集合
type SafeRanges struct {
	Temperature SafeRange `json:"temperature"` // 22-28 °C
	PH          SafeRange `json:"ph"`           // 6.5-8.0
	Oxygen      SafeRange `json:"oxygen"`       // 5-9 mg/L
	Ammonia     SafeRange `json:"ammonia"`      // 0-0.5 mg/L
	WaterLevel  SafeRange `json:"water_level"`  // 60-90 %
	TDS         SafeRange `json:"tds"`          // 100-500 ppm
	Nitrate     SafeRange `json:"nitrate"`      // 0-50 mg/L
	Nitrite     SafeRange `json:"nitrite"`      // 0-0.5 mg/L
	Chloride    SafeRange `json:"chloride"`      // 0-250 mg/L
}

// DefaultSafeRanges 返回标准鱼缸安全范围
func DefaultSafeRanges() SafeRanges {
	return SafeRanges{
		Temperature: SafeRange{Min: 22, Max: 28},
		PH:          SafeRange{Min: 6.5, Max: 8.0},
		Oxygen:      SafeRange{Min: 5, Max: 9},
		Ammonia:     SafeRange{Min: 0, Max: 0.5},
		WaterLevel:  SafeRange{Min: 60, Max: 90},
		TDS:         SafeRange{Min: 100, Max: 500},
		Nitrate:     SafeRange{Min: 0, Max: 50},
		Nitrite:     SafeRange{Min: 0, Max: 0.5},
		Chloride:    SafeRange{Min: 0, Max: 250},
	}
}

// PredictionPoint 单个预测时刻的状态
type PredictionPoint struct {
	State     StateVector `json:"state"`
	Time      time.Time   `json:"time"`
	IsAnomaly bool        `json:"is_anomaly"` // 任一维度超安全范围则为 true
}

// PredictResponse 预测端点的统一返回结构
type PredictResponse struct {
	TankID            string            `json:"tank_id"`
	Predictions       []PredictionPoint `json:"predictions"`          // 自由演化轨迹（不做任何控制）
	MPCActions        []ControlAction   `json:"mpc_actions"`          // MPC 最优控制建议（第一步）
	MPCPredicted      []PredictionPoint `json:"mpc_predicted"`        // 施加最优控制后的轨迹
	AnomalyStep       int               `json:"anomaly_step"`         // 自由演化首次异常步，-1=无异常
	AnomalyWarn       string            `json:"anomaly_warning"`      // 预警文本，空=安全
	ExtendedAnomalies map[string]string `json:"extended_anomalies,omitempty"` // 扩展参数异常描述（当前读数）
	Confidence        float64           `json:"confidence"`           // 0-1 模型置信度
	GeneratedAt       time.Time         `json:"generated_at"`
}

// ControlAction 语义化的设备动作
type ControlAction struct {
	Device string  `json:"device"` // heater/aerator/pump/feeder/light
	Action string  `json:"action"` // heat/cool/on/off/feed/dim/brighten/maintain
	Value  float64 `json:"value"`  // 归一化值 0-1
}

// KoopmanModelConfig JSON 模型文件的 Go 映射结构
// Python 训练流水线导出的 JSON 必须与此结构完全对应
type KoopmanModelConfig struct {
	Version       string             `json:"version"`
	TrainedAt     string             `json:"trained_at"`
	DeltaT        float64            `json:"delta_t"`       // 训练时步长(秒)
	StateDim      int                `json:"state_dim"`     // 5
	ControlDim    int                `json:"control_dim"`   // 5
	LiftDim       int                `json:"lift_dim"`      // 提升空间维度
	KMatrix       [][]float64        `json:"k_matrix"`      // lift_dim x lift_dim
	BMatrix       [][]float64        `json:"b_matrix"`      // lift_dim x control_dim
	ObservableCfg []ObservableConfig `json:"observables"`
	SafeRanges    SafeRanges         `json:"safe_ranges"`
	Normalization NormConfig         `json:"normalization"`
}

// ObservableConfig 描述一个可观测量函数
type ObservableConfig struct {
	Type    string    `json:"type"`    // identity/quadratic/cross/trig
	Indices []int     `json:"indices"` // 应用的状态索引
	Params  []float64 `json:"params"`  // trig 频率参数等
}

// NormConfig 存储归一化参数
type NormConfig struct {
	StateMean []float64 `json:"state_mean"`
	StateStd  []float64 `json:"state_std"`
	CtrlMean  []float64 `json:"ctrl_mean"`
	CtrlStd   []float64 `json:"ctrl_std"`
}

// DeviceActionMap 控制向量索引 → 语义设备动作映射
var DeviceActionMap = [5]struct {
	Device string
	Decr   string // 减小方向的动作
	Incr   string // 增大方向的动作
}{
	{Device: "heater", Decr: "cool", Incr: "heat"},
	{Device: "aerator", Decr: "off", Incr: "on"},
	{Device: "pump", Decr: "off", Incr: "on"},
	{Device: "feeder", Decr: "hold", Incr: "feed"},
	{Device: "light", Decr: "dim", Incr: "brighten"},
}
