"""将训练好的 Koopman 模型导出为 JSON，供 Go 侧加载推理

导出格式必须严格对应 Go 侧 domain.KoopmanModelConfig 结构体。
"""
import json
import numpy as np
from datetime import datetime, timezone

from config import DELTA_T, STATE_DIM, CONTROL_DIM


def export_model(K: np.ndarray, B: np.ndarray,
                obs, scaler_X, scaler_U,
                safe_ranges: dict,
                output_path: str):
    """将模型写入 JSON 文件。

    Args:
        K: (lift_dim, lift_dim) Koopman 矩阵
        B: (lift_dim, control_dim) 控制矩阵
        obs: ObservableDictionary 实例
        scaler_X: 拟合的 StandardScaler（状态）
        scaler_U: 拟合的 StandardScaler（控制）
        safe_ranges: dict，格式 {"temperature": {"min": 22, "max": 28}, ...}
        output_path: 输出文件路径
    """
    # 处理 std=0 的情况（某维度方差为零，防止 Go 侧除零）
    state_std = scaler_X.scale_.tolist()
    state_std = [s if s > 0 else 1.0 for s in state_std]

    ctrl_std = scaler_U.scale_.tolist()
    ctrl_std = [s if s > 0 else 1.0 for s in ctrl_std]

    model = {
        "version": "1.0.0",
        "trained_at": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
        "delta_t": DELTA_T,
        "state_dim": STATE_DIM,
        "control_dim": CONTROL_DIM,
        "lift_dim": obs.lift_dim,
        "k_matrix": K.tolist(),
        "b_matrix": B.tolist(),
        "observables": obs.get_config(),
        "safe_ranges": safe_ranges,
        "normalization": {
            "state_mean": scaler_X.mean_.tolist(),
            "state_std": state_std,
            "ctrl_mean": scaler_U.mean_.tolist(),
            "ctrl_std": ctrl_std,
        }
    }

    with open(output_path, 'w', encoding='utf-8') as f:
        json.dump(model, f, indent=2, ensure_ascii=False)

    print(f"模型已导出到 {output_path}")
    print(f"  K矩阵: {K.shape}")
    print(f"  B矩阵: {B.shape}")
    print(f"  提升维度: {obs.lift_dim}")
