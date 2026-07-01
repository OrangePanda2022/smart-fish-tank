"""Extended Dynamic Mode Decomposition (EDMD) 实现

通过将状态提升到可观测量空间 psi(x)，使非线性动力系统在提升空间中近似线性，
从而学习 Koopman 算子 K 和控制矩阵 B。

关键约束：此文件中的 ObservableDictionary 必须与 Go 侧 Model.Lift() 完全一致。
"""
import numpy as np
from sklearn.preprocessing import StandardScaler
from config import STATE_DIM, CONTROL_DIM, TRIG_FREQS


class ObservableDictionary:
    """定义提升函数 psi(x)，与 Go 侧 Model.Lift() 完全对齐。

    提升字典包含:
    - identity: x1, x2, ..., xn（5项）
    - quadratic: x1^2, x2^2, ..., xn^2（5项）
    - cross: xi*xj for i < j（10项）
    - trig: sin(w_k * xi), cos(w_k * xi) 对每个频率 w_k 和每个状态 xi
            （5 * len(TRIG_FREQS) * 2 项）

    总提升维度 = 5 + 5 + 10 + 5*len(TRIG_FREQS)*2
    """

    def __init__(self, state_dim=STATE_DIM, trig_freqs=None):
        self.state_dim = state_dim
        self.trig_freqs = trig_freqs or list(TRIG_FREQS)
        self._compute_lift_dim()

    def _compute_lift_dim(self):
        n = self.state_dim
        self.lift_dim = n                         # identity
        self.lift_dim += n                        # quadratic
        self.lift_dim += n * (n - 1) // 2         # cross terms
        self.lift_dim += n * len(self.trig_freqs) * 2  # sin+cos

    def lift(self, X: np.ndarray) -> np.ndarray:
        """对状态矩阵应用可观测量字典。

        Args:
            X: (n_samples, state_dim) 归一化后的状态矩阵
        Returns:
            Psi: (n_samples, lift_dim) 提升后的状态矩阵
        """
        n_samples = X.shape[0]
        n = self.state_dim
        Psi = np.zeros((n_samples, self.lift_dim))

        col = 0
        # identity: x1, x2, ..., xn
        Psi[:, col:col+n] = X
        col += n

        # quadratic: x1^2, x2^2, ..., xn^2
        Psi[:, col:col+n] = X ** 2
        col += n

        # cross terms: xi*xj for i < j
        for i in range(n):
            for j in range(i + 1, n):
                Psi[:, col] = X[:, i] * X[:, j]
                col += 1

        # trigonometric: sin(w_k * xi), cos(w_k * xi)
        for omega in self.trig_freqs:
            for i in range(n):
                Psi[:, col] = np.sin(omega * X[:, i])
                col += 1
                Psi[:, col] = np.cos(omega * X[:, i])
                col += 1

        assert col == self.lift_dim, f"列数 {col} != lift_dim {self.lift_dim}"
        return Psi

    def get_config(self):
        """返回与 Go 侧 domain.ObservableConfig 对应的序列化配置"""
        configs = []

        configs.append({
            "type": "identity",
            "indices": list(range(self.state_dim)),
            "params": []
        })
        configs.append({
            "type": "quadratic",
            "indices": list(range(self.state_dim)),
            "params": []
        })
        configs.append({
            "type": "cross",
            "indices": list(range(self.state_dim)),
            "params": []
        })
        for omega in self.trig_freqs:
            configs.append({
                "type": "trig",
                "indices": list(range(self.state_dim)),
                "params": [omega]
            })

        return configs


def edmd(X: np.ndarray, Y: np.ndarray, U: np.ndarray, obs: ObservableDictionary):
    """带控制输入的 EDMD。

    给定状态对 (x_k, x_{k+1}) 和控制 u_k:
      1. 提升状态: Psi_k = psi(x_k)
      2. 最小二乘求解: [K | B] = Psi_{k+1} @ pinv([Psi_k; u_k^T])

    Args:
        X: (n_samples, state_dim) 时刻 k 的状态
        Y: (n_samples, state_dim) 时刻 k+1 的状态
        U: (n_samples, control_dim) 时刻 k 的控制输入
        obs: ObservableDictionary 实例

    Returns:
        K: (lift_dim, lift_dim) Koopman 矩阵
        B: (lift_dim, control_dim) 控制矩阵
        scaler_X: 拟合的 StandardScaler（状态归一化）
        scaler_U: 拟合的 StandardScaler（控制归一化）
    """
    # 归一化状态
    scaler_X = StandardScaler()
    X_norm = scaler_X.fit_transform(X)
    Y_norm = scaler_X.transform(Y)

    # 归一化控制
    scaler_U = StandardScaler()
    U_norm = scaler_U.fit_transform(U)

    # 提升状态
    Psi_X = obs.lift(X_norm)  # (n_samples, lift_dim)
    Psi_Y = obs.lift(Y_norm)  # (n_samples, lift_dim)

    # 增广回归矩阵: [Psi_k | u_k]
    Aug = np.hstack([Psi_X, U_norm])  # (n_samples, lift_dim + ctrl_dim)

    # 最小二乘求解: [K | B] = Psi_Y^T @ pinv(Aug^T)
    KB = Psi_Y.T @ np.linalg.pinv(Aug.T)  # (lift_dim, lift_dim + ctrl_dim)

    K = KB[:, :obs.lift_dim]
    B = KB[:, obs.lift_dim:]

    return K, B, scaler_X, scaler_U
