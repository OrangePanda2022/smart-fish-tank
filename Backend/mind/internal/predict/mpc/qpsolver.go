package mpc

import (
	"fmt"
	"math"

	"mind/internal/domain"
	"mind/internal/predict/koopman"

	"gonum.org/v1/gonum/mat"
)

// QPSolver 使用投影梯度法求解带 box 约束的二次规划
//
//	min  0.5·UᵀHU + fᵀU
//	s.t. lb ≤ U ≤ ub
type QPSolver struct {
	maxIter  int
	tol      float64
	stepSize float64
}

// NewQPSolver 创建 QP 求解器
func NewQPSolver() *QPSolver {
	return &QPSolver{
		maxIter:  500,
		tol:      1e-6,
		stepSize: 0.01,
	}
}

// Solve 求解 QP，返回最优控制序列的扁平化向量
func (s *QPSolver) Solve(H *mat.Dense, f *mat.VecDense, lb, ub []float64) ([]float64, error) {
	n := f.Len()
	if r, c := H.Dims(); r != n || c != n {
		return nil, fmt.Errorf("Hessian维度不匹配: H是%dx%d, f是%d", r, c, n)
	}

	// 初始化为边界中点
	U := make([]float64, n)
	for i := range U {
		U[i] = (lb[i] + ub[i]) / 2
	}

	for iter := 0; iter < s.maxIter; iter++ {
		// 计算梯度: g = H·U + f
		uVec := mat.NewVecDense(n, U)
		var g mat.VecDense
		g.MulVec(H, uVec)
		g.AddVec(&g, f)

		gradNorm := mat.Norm(&g, 2)
		if gradNorm < 1e-12 {
			break
		}

		// 回溯线搜索
		alpha := s.stepSize
		for step := 0; step < 20; step++ {
			newU := make([]float64, n)
			for i := 0; i < n; i++ {
				newU[i] = U[i] - alpha*g.AtVec(i)
				// 投影到 box 约束
				newU[i] = math.Max(lb[i], math.Min(ub[i], newU[i]))
			}

			// Armijo 充分下降条件
			newVec := mat.NewVecDense(n, newU)
			var newG mat.VecDense
			newG.MulVec(H, newVec)
			newObj := 0.5*dot(newU, newG.RawVector().Data) + dot(newU, f.RawVector().Data)
			curObj := 0.5*dot(U, g.RawVector().Data) + dot(U, f.RawVector().Data)

			if newObj < curObj-1e-4*alpha*gradNorm*gradNorm {
				U = newU
				break
			}
			alpha *= 0.5
		}

		// 检查投影梯度收敛
		converged := true
		uVec = mat.NewVecDense(n, U)
		var pg mat.VecDense
		pg.MulVec(H, uVec)
		pg.AddVec(&pg, f)

		for i := 0; i < n; i++ {
			pgVal := pg.AtVec(i)
			if U[i] <= lb[i] && pgVal > 0 {
				pgVal = 0
			}
			if U[i] >= ub[i] && pgVal < 0 {
				pgVal = 0
			}
			if math.Abs(pgVal) > s.tol {
				converged = false
				break
			}
		}
		if converged {
			break
		}
	}

	return U, nil
}

// BuildQP 从 Koopman 动力学和代价函数构建 QP 矩阵
//
//	Koopman 提升动力学: ψ_{k+1} = K·ψ_k + B·u_k
//	物理态提取: z_k = C·ψ_k (C 取前 stateDim 维)
//
// 预测表达:
//
//	Z = Z_free + Γ·U
//
// 代价归约为标准 QP:
//
//	min 0.5·Uᵀ·(ΓᵀQ̄Γ + R̄)·U + (Z_freeᵀQ̄Γ)ᵀ·U
func BuildQP(
	model *koopman.Model,
	cost *CostMatrices,
	initialState domain.StateVector,
	setpoint domain.StateVector,
	horizon int,
) (H *mat.Dense, f *mat.VecDense, lb []float64, ub []float64, err error) {
	cfg := model.Config()
	liftDim := cfg.LiftDim
	ctrlDim := cfg.ControlDim
	stateDim := cfg.StateDim

	// C 矩阵: 从提升态提取物理量 (5 x lift_dim)，前5列为单位阵
	C := mat.NewDense(stateDim, liftDim, nil)
	for i := 0; i < stateDim; i++ {
		C.Set(i, i, 1.0)
	}

	// 计算自由响应（零控制下提升态演化）
	liftedInit, err := koopman.ComputeLiftedState(model, initialState)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("计算初始提升态失败: %w", err)
	}
	psiFree := make([]*mat.VecDense, horizon+1)
	psiFree[0] = mat.NewVecDense(liftDim, liftedInit)

	for k := 0; k < horizon; k++ {
		var next mat.VecDense
		next.MulVec(model.K, psiFree[k])
		psiFree[k+1] = &next
	}

	// 构建 Γ 矩阵: (N·stateDim) x (N·ctrlDim)
	// Γ[k,j] 块 = C·K^{k-j}·B  对 k > j，0 对 k <= j
	totalStateDim := horizon * stateDim
	totalCtrlDim := horizon * ctrlDim
	Gamma := mat.NewDense(totalStateDim, totalCtrlDim, nil)

	for k := 0; k < horizon; k++ {
		for j := 0; j <= k; j++ {
			var Kpow mat.Dense
			matPow(&Kpow, model.K, k-j)

			var KpowB mat.Dense
			KpowB.Mul(&Kpow, model.B)

			var CKpowB mat.Dense
			CKpowB.Mul(C, &KpowB)

			// 放入块 (k*stateDim, j*ctrlDim)
			for r := 0; r < stateDim; r++ {
				for c := 0; c < ctrlDim; c++ {
					Gamma.Set(k*stateDim+r, j*ctrlDim+c, CKpowB.At(r, c))
				}
			}
		}
	}

	// 自由响应向量 Z_free (N·stateDim)
	zFree := make([]float64, totalStateDim)
	for k := 0; k < horizon; k++ {
		physicalState := koopman.ExtractPhysicalState(
			psiFree[k+1].RawVector().Data, cfg.Normalization)
		for i := 0; i < stateDim; i++ {
			zFree[k*stateDim+i] = physicalState[i] - setpoint[i]
		}
	}
	zFreeVec := mat.NewVecDense(totalStateDim, zFree)

	// Q̄: 块对角 (Q 重复 horizon-1 次, 然后 QN)
	Qbar := mat.NewDiagDense(totalStateDim, nil)
	for k := 0; k < horizon-1; k++ {
		for i := 0; i < stateDim; i++ {
			Qbar.SetDiag(k*stateDim+i, cost.Q.At(i, i))
		}
	}
	for i := 0; i < stateDim; i++ {
		Qbar.SetDiag((horizon-1)*stateDim+i, cost.QN.At(i, i))
	}

	// R̄: 块对角 (R 重复 horizon 次)
	Rbar := mat.NewDiagDense(totalCtrlDim, nil)
	for k := 0; k < horizon; k++ {
		for i := 0; i < ctrlDim; i++ {
			Rbar.SetDiag(k*ctrlDim+i, cost.R.At(i, i))
		}
	}

	// H = Γᵀ·Q̄·Γ + R̄
	// 注意乘法顺序: 先 Q̄·Γ (stateDim×ctrlDim)，再 Γᵀ·(Q̄·Γ)
	var GammaQ mat.Dense
	GammaQ.Mul(Qbar, Gamma)

	var GQG mat.Dense
	GQG.Mul(Gamma.T(), &GammaQ)

	var Hdense mat.Dense
	Hdense.Add(&GQG, Rbar)

	// f = Γᵀ·Q̄·Z_free
	var QbarZ mat.VecDense
	QbarZ.MulVec(Qbar, zFreeVec)

	var fVec mat.VecDense
	fVec.MulVec(Gamma.T(), &QbarZ)

	// 控制边界
	ctrlLb, ctrlUb := ControlBounds()
	lb = make([]float64, totalCtrlDim)
	ub = make([]float64, totalCtrlDim)
	for k := 0; k < horizon; k++ {
		for i := 0; i < ctrlDim; i++ {
			lb[k*ctrlDim+i] = ctrlLb[i]
			ub[k*ctrlDim+i] = ctrlUb[i]
		}
	}

	return &Hdense, &fVec, lb, ub, nil
}

// matPow 计算 M^n（密集矩阵幂）
func matPow(result *mat.Dense, M *mat.Dense, n int) {
	r, c := M.Dims()
	if n == 0 {
		*result = *mat.NewDense(r, c, nil)
		for i := 0; i < r; i++ {
			result.Set(i, i, 1.0)
		}
		return
	}
	// 确保 result 有正确的尺寸，零值 mat.Dense 的 Copy 不会自动扩展
	*result = *mat.NewDense(r, c, nil)
	result.Copy(M)
	if n == 1 {
		return
	}
	for i := 1; i < n; i++ {
		var tmp mat.Dense
		tmp.Mul(result, M)
		result.Copy(&tmp)
	}
}

// dot 计算两个 float64 切片的内积
func dot(a, b []float64) float64 {
	var s float64
	for i := range a {
		s += a[i] * b[i]
	}
	return s
}
