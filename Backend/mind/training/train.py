"""Koopman + MPC 预测模型训练入口

使用方法:
    cd training/
    pip install -r requirements.txt
    python train.py --tank-id tank-001 --hours 720
    python train.py --synthetic            # 使用合成数据训练

训练流程:
    1. 获取历史传感器数据（InfluxDB 或合成数据）
    2. 重采样到5分钟间隔并线性插值
    3. 使用 EDMD 学习 Koopman 算子 K 和控制矩阵 B
    4. 验证重建误差 (RMSE)
    5. 导出 JSON 模型文件供 Go 侧推理
"""
import argparse
import numpy as np

from influxdb_client import fetch_training_data, fetch_control_data
from edmd import ObservableDictionary, edmd
from model_export import export_model
from config import STATE_DIM, STATE_COLS, CONTROL_DIM, MODEL_OUTPUT_PATH, TRIG_FREQS

# 标准安全范围（与 Go 侧 domain.DefaultSafeRanges() 一致）
SAFE_RANGES = {
    "temperature": {"min": 22, "max": 28},
    "ph":          {"min": 6.5, "max": 8.0},
    "oxygen":      {"min": 5, "max": 9},
    "ammonia":     {"min": 0, "max": 0.5},
    "water_level": {"min": 60, "max": 90},
    "tds":         {"min": 100, "max": 500},
    "nitrate":     {"min": 0, "max": 50},
    "nitrite":     {"min": 0, "max": 0.5},
    "chloride":    {"min": 0, "max": 250},
}


def main():
    parser = argparse.ArgumentParser(description="训练鱼缸 Koopman 预测模型")
    parser.add_argument("--tank-id", default="019f1807-dbca-753d-afb8-6b2db160fa8d", help="鱼缸ID")
    parser.add_argument("--hours", type=int, default=720, help="历史数据时长(小时)，默认720=30天")
    parser.add_argument("--output", default=MODEL_OUTPUT_PATH, help="模型输出路径")
    parser.add_argument("--synthetic", action="store_true", help="使用合成数据训练（无需 InfluxDB）")
    args = parser.parse_args()

    # 1. 获取数据
    if args.synthetic:
        print("使用合成数据训练...")
        from synthetic_data import generate_synthetic_data
        X, Y, U = generate_synthetic_data(n_steps=10000)
    else:
        print(f"获取鱼缸 {args.tank_id} 的历史数据 ({args.hours} 小时)...")
        df_state = fetch_training_data(args.tank_id, args.hours)
        df_ctrl = fetch_control_data(args.tank_id, args.hours)

        if len(df_state) < 100:
            print(f"数据不足: 仅 {len(df_state)} 条记录，需要至少 100 条。尝试 --synthetic 模式。")
            return

        X = df_state[STATE_COLS].values[:-1]   # x_k
        Y = df_state[STATE_COLS].values[1:]    # x_{k+1}
        U = df_ctrl.values[:-1]               # u_k

    # 2. EDMD 训练
    obs = ObservableDictionary(trig_freqs=TRIG_FREQS)
    print(f"训练 EDMD (状态维度: {STATE_DIM}, 控制维度: {CONTROL_DIM}, 提升维度: {obs.lift_dim})...")
    K, B, scaler_X, scaler_U = edmd(X, Y, U, obs)

    # 3. 验证重建误差
    Psi_X = obs.lift(scaler_X.transform(X))
    Psi_Y_pred = (K @ Psi_X.T + B @ U.T).T
    Y_pred = scaler_X.inverse_transform(Psi_Y_pred[:, :STATE_DIM])
    rmse = np.sqrt(np.mean((Y - Y_pred) ** 2))
    print(f"验证 RMSE: {rmse:.4f}")

    # 逐维度 RMSE
    for i, col in enumerate(STATE_COLS):
        col_rmse = np.sqrt(np.mean((Y[:, i] - Y_pred[:, i]) ** 2))
        print(f"  {col}: RMSE = {col_rmse:.4f}")

    # 4. 导出模型
    export_model(K, B, obs, scaler_X, scaler_U, SAFE_RANGES, args.output)
    print("训练完成！")


if __name__ == "__main__":
    main()
