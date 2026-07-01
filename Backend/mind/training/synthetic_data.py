"""合成鱼缸传感器数据生成器

生成符合鱼缸物理规律的 9 维状态时间序列 + 5 维控制输入，
用于在没有真实历史数据时训练 Koopman 模型。

9 维状态: [temperature, ph, oxygen, ammonia, water_level, tds, nitrate, nitrite, chloride]
5 维控制: [heater_power, aerator_power, pump_power, feeder_dispense, light_intensity]
"""
import numpy as np
from config import STATE_DIM, CONTROL_DIM, DELTA_T

# 默认稳态参数
STEADY_STATE = np.array([25.0, 7.0, 7.0, 0.1, 75.0, 250.0, 10.0, 0.05, 50.0])


def generate_synthetic_data(n_steps: int = 10000, seed: int = 42) -> tuple:
    """生成合成状态转移对 (X, Y, U) 供 EDMD 训练。

    Args:
        n_steps: 总步数（含首尾）
        seed: 随机种子

    Returns:
        X: (n_steps-1, 9) 时刻 k 的状态
        Y: (n_steps-1, 9) 时刻 k+1 的状态
        U: (n_steps-1, 5) 时刻 k 的控制输入
    """
    rng = np.random.default_rng(seed)
    state = STEADY_STATE.copy()
    states = np.zeros((n_steps, STATE_DIM))
    controls = np.zeros((n_steps, CONTROL_DIM))

    for t in range(n_steps):
        # 保存当前状态
        states[t] = state

        # 生成控制输入
        u = _generate_control(rng, state, t)
        controls[t] = u

        # 状态转移
        state = _state_transition(state, u, rng, t)

    X = states[:-1]
    Y = states[1:]
    U = controls[:-1]

    return X, Y, U


def _generate_control(rng: np.random.Generator, state: np.ndarray, t: int) -> np.ndarray:
    """生成控制输入，带有简单的逻辑约束。"""
    # 基础随机控制
    u = np.array([
        rng.uniform(0.2, 0.8),   # heater: 温度低时加大
        rng.uniform(0.2, 0.8),   # aerator: 氧低时加大
        rng.uniform(0.0, 0.4),   # pump: 低频换水
        rng.uniform(0.0, 0.3),   # feeder: 适量
        rng.uniform(0.3, 0.7),   # light: 中等
    ])

    # 温度低时加大加热
    if state[0] < 23:
        u[0] = rng.uniform(0.6, 1.0)
    elif state[0] > 27:
        u[0] = rng.uniform(0.0, 0.2)

    # 氧含量低时加大增氧
    if state[2] < 5.5:
        u[1] = rng.uniform(0.6, 1.0)

    # 氨氮高时加大换水 + 减少喂食
    if state[3] > 0.3:
        u[2] = rng.uniform(0.3, 0.6)
        u[3] = rng.uniform(0.0, 0.05)

    # 亚硝酸盐高时加大换水
    if state[7] > 0.3:
        u[2] = max(u[2], rng.uniform(0.3, 0.5))

    return u


def _state_transition(state: np.ndarray, u: np.ndarray,
                      rng: np.random.Generator, t: int) -> np.ndarray:
    """模拟一步状态转移（5 分钟步长）。"""
    dt_h = DELTA_T / 3600.0  # 转换为小时
    s = state.copy()

    # 0: 温度 — 正弦日周期 + 加热器效应 + 噪声
    day_phase = 2 * np.pi * (t * DELTA_T) / 86400.0
    temp_drift = 0.3 * np.sin(day_phase) * dt_h  # 日周期缓慢漂移
    heater_effect = 2.0 * (u[0] - 0.5) * dt_h     # 加热器功率偏差
    s[0] += temp_drift + heater_effect + rng.normal(0, 0.05)

    # 1: pH — 慢漂移 + 噪声
    s[1] += rng.normal(0, 0.02) + 0.001 * (7.0 - s[1])  # 均值回归

    # 2: 溶解氧 — 与温度反相关 + 增氧机效应 + 噪声
    o2_temp = -0.1 * (s[0] - 25.0) * dt_h
    aerator_effect = 0.5 * (u[1] - 0.5) * dt_h
    s[2] += o2_temp + aerator_effect + rng.normal(0, 0.05)

    # 3: 氨氮 — 缓慢积累（喂食） + 换水稀释 + 硝化消耗 + 噪声
    feed_input = 0.02 * u[3] * dt_h       # 喂食产生氨
    pump_dilution = -0.3 * u[2] * s[3] * dt_h  # 换水稀释
    bio_consumption = -0.005 * s[3] * dt_h      # 生物过滤消耗
    s[3] += feed_input + pump_dilution + bio_consumption + rng.normal(0, 0.005)
    s[3] = max(0.0, s[3])

    # 4: 水位 — 缓慢蒸发 + 补水效应 + 噪声
    evaporation = -0.01 * dt_h
    pump_effect = 2.0 * u[2] * dt_h
    s[4] += evaporation + pump_effect + rng.normal(0, 0.1)

    # 5: TDS — 与换水反相关 + 缓慢积累 + 噪声
    tds_dilution = -5.0 * u[2] * dt_h   # 换水降低 TDS
    tds_accumulation = 0.5 * dt_h        # 蒸发和矿物质溶解缓慢增加
    s[5] += tds_dilution + tds_accumulation + rng.normal(0, 1.0)
    s[5] = max(50.0, s[5])

    # 6: 硝酸盐 — 氮循环积累 + 换水稀释 + 噪声
    # 氨 → 亚硝酸盐 → 硝酸盐（硝化反应终产物，缓慢积累）
    nitrate_from_cycle = 0.01 * s[3] * dt_h  # 氨氮转化
    nitrate_dilution = -2.0 * u[2] * s[6] / 50.0 * dt_h  # 换水稀释
    s[6] += nitrate_from_cycle + nitrate_dilution + rng.normal(0, 0.2)
    s[6] = max(0.0, s[6])

    # 7: 亚硝酸盐 — 氮循环中间体 + 换水稀释 + 噪声
    # 氨 → 亚硝酸盐（快速） → 硝酸盐（较慢），亚硝酸盐是中间态
    nitrite_from_ammonia = 0.05 * s[3] * dt_h   # 氨氮转化为亚硝酸盐
    nitrite_to_nitrate = -0.02 * s[7] * dt_h    # 亚硝酸盐继续转化为硝酸盐
    nitrite_dilution = -0.5 * u[2] * s[7] / 0.5 * dt_h  # 换水稀释
    s[7] += nitrite_from_ammonia + nitrite_to_nitrate + nitrite_dilution + rng.normal(0, 0.005)
    s[7] = max(0.0, s[7])

    # 8: 氯离子 — 与换水反相关 + 缓慢积累 + 噪声
    chloride_dilution = -3.0 * u[2] * dt_h
    chloride_accumulation = 0.2 * dt_h
    s[8] += chloride_dilution + chloride_accumulation + rng.normal(0, 0.3)
    s[8] = max(0.0, s[8])

    return s


if __name__ == "__main__":
    # 快速测试
    X, Y, U = generate_synthetic_data(100)
    print(f"X shape: {X.shape}, Y shape: {Y.shape}, U shape: {U.shape}")
    print(f"X sample: {X[0]}")
    print(f"Y sample: {Y[0]}")
    print(f"U sample: {U[0]}")
