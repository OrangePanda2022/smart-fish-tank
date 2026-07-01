package koopman

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sync"

	"mind/internal/domain"

	"gonum.org/v1/gonum/mat"
)

// Model 封装加载的 Koopman 模型，支持并发安全读取和热重载
type Model struct {
	mu        sync.RWMutex
	cfg       domain.KoopmanModelConfig
	K         *mat.Dense // lift_dim x lift_dim Koopman 矩阵
	B         *mat.Dense // lift_dim x control_dim 控制矩阵
	loaded    bool
	untrained bool // K ≈ I，模型未训练
}

// IsUntrained 返回模型是否为占位符（K ≈ 单位阵），未经过真实数据训练
func (m *Model) IsUntrained() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.untrained
}

// detectUntrained 检查 K 矩阵是否接近单位阵
func detectUntrained(K *mat.Dense, liftDim int) bool {
	tolerance := 0.01
	for i := 0; i < liftDim; i++ {
		for j := 0; j < liftDim; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if math.Abs(K.At(i, j)-expected) > tolerance {
				return false
			}
		}
	}
	return true
}

// LoadModel 从 JSON 文件加载模型并构建 K、B 矩阵
func LoadModel(path string) (*Model, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取模型文件失败: %w", err)
	}

	var cfg domain.KoopmanModelConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析模型JSON失败: %w", err)
	}

	// 验证并修复矩阵维度
	// K 矩阵：行不足补零行，行短补零列；行/列过多则截断
	kRows := len(cfg.KMatrix)
	kCols := 0
	if kRows > 0 {
		kCols = len(cfg.KMatrix[0])
	}
	if kRows != cfg.LiftDim || kCols != cfg.LiftDim {
		log.Printf("K矩阵维度 %dx%d 与 lift_dim %d 不匹配，自动补齐/截断",
			kRows, kCols, cfg.LiftDim)
		// 补齐/截断行
		for len(cfg.KMatrix) < cfg.LiftDim {
			cfg.KMatrix = append(cfg.KMatrix, make([]float64, cfg.LiftDim))
		}
		cfg.KMatrix = cfg.KMatrix[:cfg.LiftDim]
		// 补齐/截断每行的列
		for i, row := range cfg.KMatrix {
			if len(row) < cfg.LiftDim {
				extended := make([]float64, cfg.LiftDim)
				copy(extended, row)
				cfg.KMatrix[i] = extended
			} else if len(row) > cfg.LiftDim {
				cfg.KMatrix[i] = row[:cfg.LiftDim]
			}
		}
	}
	// B 矩阵：行不足补零行，行短补零列
	bRows := len(cfg.BMatrix)
	bCols := 0
	if bRows > 0 {
		bCols = len(cfg.BMatrix[0])
	}
	ctrlDim := cfg.ControlDim
	if ctrlDim == 0 {
		ctrlDim = len(cfg.BMatrix[0]) // 从矩阵推断
	}
	if bRows != cfg.LiftDim || bCols != ctrlDim {
		log.Printf("B矩阵维度 %dx%d 与 lift_dim %d / control_dim %d 不匹配，自动补齐/截断",
			bRows, bCols, cfg.LiftDim, ctrlDim)
		for len(cfg.BMatrix) < cfg.LiftDim {
			cfg.BMatrix = append(cfg.BMatrix, make([]float64, ctrlDim))
		}
		cfg.BMatrix = cfg.BMatrix[:cfg.LiftDim]
		for i, row := range cfg.BMatrix {
			if len(row) < ctrlDim {
				extended := make([]float64, ctrlDim)
				copy(extended, row)
				cfg.BMatrix[i] = extended
			} else if len(row) > ctrlDim {
				cfg.BMatrix[i] = row[:ctrlDim]
			}
		}
	}

	k := mat.NewDense(cfg.LiftDim, cfg.LiftDim, flatten(cfg.KMatrix))
	bRowsB, bColsB := cfg.LiftDim, cfg.ControlDim
	if bColsB == 0 {
		bColsB = len(cfg.BMatrix[0])
	}
	b := mat.NewDense(bRowsB, bColsB, flatten(cfg.BMatrix))

	untrained := detectUntrained(k, cfg.LiftDim)
	if untrained {
		log.Printf("检测到 K ≈ I，模型为占位符（未训练），将使用启发式控制")
	}

	return &Model{cfg: cfg, K: k, B: b, loaded: true, untrained: untrained}, nil
}

// HotReload 从磁盘重新加载模型，失败不影响现有模型
func (m *Model) HotReload(path string) error {
	newModel, err := LoadModel(path)
	if err != nil {
		return fmt.Errorf("热重载失败: %w", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg = newModel.cfg
	m.K = newModel.K
	m.B = newModel.B
	m.loaded = true
	return nil
}

// IsLoaded 返回模型是否已成功加载
func (m *Model) IsLoaded() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.loaded
}

// Config 返回模型配置的副本
func (m *Model) Config() domain.KoopmanModelConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg
}

// Lift 将归一化后的状态向量提升到可观测量空间 ψ(x)
// 必须与 Python 训练流水线的 ObservableDictionary.lift() 完全一致
func (m *Model) Lift(s domain.StateVector) []float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	obs := m.cfg.ObservableCfg
	lifted := make([]float64, 0, m.cfg.LiftDim)

	for _, ob := range obs {
		switch ob.Type {
		case "identity":
			for _, idx := range ob.Indices {
				lifted = append(lifted, s[idx])
			}
		case "quadratic":
			for _, idx := range ob.Indices {
				lifted = append(lifted, s[idx]*s[idx])
			}
		case "cross":
			// 索引对 xi*xj, i < j
			for i := 0; i < len(ob.Indices); i++ {
				for j := i + 1; j < len(ob.Indices); j++ {
					lifted = append(lifted, s[ob.Indices[i]]*s[ob.Indices[j]])
				}
			}
		case "trig":
			// sin(omega * xi), cos(omega * xi) 对每个频率参数
			for _, idx := range ob.Indices {
				for _, omega := range ob.Params {
					lifted = append(lifted, math.Sin(omega*s[idx]), math.Cos(omega*s[idx]))
				}
			}
		}
	}

	return lifted
}

// flatten 将二维 float64 数组展平为一维
func flatten(matrix [][]float64) []float64 {
	flat := make([]float64, 0, len(matrix)*len(matrix[0]))
	for _, row := range matrix {
		flat = append(flat, row...)
	}
	return flat
}
