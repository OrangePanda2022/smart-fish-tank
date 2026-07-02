import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../domain/entities/sensor_reading.dart';
import 'device_control.dart';

List<DeviceControl> buildDeviceCatalog(SensorReading sensor) {
  return [
    DeviceControl(
      '灯光',
      'ON',
      CupertinoIcons.lightbulb_fill,
      AppColors.yellow,
      id: 'light',
      description: '模拟自然昼夜节律，当前处于日间补光。',
      primaryAction: '调节亮度',
      maintenance: '建议每天 08:00-20:00 开启',
      metrics: const [
        DeviceMetric(label: '亮度', value: '76%'),
        DeviceMetric(label: '模式', value: '日光'),
        DeviceMetric(label: '定时', value: '20:00 关闭'),
      ],
    ),
    DeviceControl(
      '水泵',
      'ON',
      Icons.waves_rounded,
      AppColors.sky,
      id: 'pump',
      description: '保持水体循环，辅助过滤系统稳定工作。',
      primaryAction: '调整流量',
      maintenance: '滤芯建议 14 天后检查',
      metrics: [
        const DeviceMetric(label: '流量', value: '中档'),
        DeviceMetric(
          label: '水位',
          value: '${sensor.waterLevel.toStringAsFixed(0)}%',
        ),
        const DeviceMetric(label: '运行', value: '连续'),
      ],
    ),
    DeviceControl(
      '氧气',
      'ON',
      CupertinoIcons.circle_grid_hex,
      AppColors.bubble,
      id: 'aerator',
      description: '根据溶氧水平维持鱼缸含氧量。',
      primaryAction: '增强供氧',
      maintenance: '气石状态良好',
      metrics: [
        DeviceMetric(label: '溶氧', value: sensor.oxygen.toStringAsFixed(1)),
        const DeviceMetric(label: '强度', value: '65%'),
        const DeviceMetric(label: '模式', value: '智能'),
      ],
    ),
    DeviceControl(
      '加热',
      '${sensor.temperature.toStringAsFixed(1)}°C',
      CupertinoIcons.thermometer,
      AppColors.orange,
      id: 'heater',
      description: '自动维持适合观赏鱼的水温区间。',
      primaryAction: '设置温度',
      maintenance: '目标温度 26.0°C',
      metrics: [
        DeviceMetric(
          label: '当前水温',
          value: '${sensor.temperature.toStringAsFixed(1)}°C',
        ),
        const DeviceMetric(label: '目标', value: '26.0°C'),
        const DeviceMetric(label: '功率', value: '42%'),
      ],
    ),
    const DeviceControl(
      '喂食',
      '08:00 已喂食',
      Icons.fastfood_rounded,
      AppColors.feed,
      id: 'feeder',
      description: '自动投喂已按计划完成。',
      primaryAction: '立即投喂',
      maintenance: '饲料余量充足',
      metrics: [
        DeviceMetric(label: '今日次数', value: '1/2'),
        DeviceMetric(label: '下次', value: '18:00'),
        DeviceMetric(label: '份量', value: '标准'),
      ],
    ),
    const DeviceControl(
      '换水',
      '剩余 7 天',
      CupertinoIcons.drop_fill,
      AppColors.primary,
      id: 'water_change',
      description: '记录换水周期，提醒维持稳定水质。',
      primaryAction: '记录换水',
      maintenance: '下次换水 7 天后',
      metrics: [
        DeviceMetric(label: '周期', value: '14 天'),
        DeviceMetric(label: '上次', value: '7 天前'),
        DeviceMetric(label: '建议量', value: '20%'),
      ],
    ),
    const DeviceControl(
      '摄像头',
      '在线',
      CupertinoIcons.camera_fill,
      AppColors.camera,
      id: 'camera',
      description: '远程查看鱼缸画面和设备状态。',
      primaryAction: '查看画面',
      maintenance: '网络连接稳定',
      metrics: [
        DeviceMetric(label: '清晰度', value: '1080P'),
        DeviceMetric(label: '延迟', value: '低'),
        DeviceMetric(label: '存储', value: '本地'),
      ],
    ),
    const DeviceControl(
      '过滤',
      'ON',
      CupertinoIcons.slider_horizontal_3,
      AppColors.teal,
      id: 'filter',
      description: '过滤杂质并辅助稳定氨氮、硝酸盐指标。',
      primaryAction: '切换模式',
      maintenance: '建议 3 天后清洗滤棉',
      metrics: [
        DeviceMetric(label: '模式', value: '自动'),
        DeviceMetric(label: '氨氮', value: '稳定'),
        DeviceMetric(label: '滤棉', value: '良好'),
      ],
    ),
    const DeviceControl(
      '备用插座',
      '离线',
      Icons.power_rounded,
      AppColors.muted,
      id: 'outlet',
      description: '可接入扩展设备，当前未连接。',
      primaryAction: '重新连接',
      maintenance: '检查电源与网络',
      isOnline: false,
      isEnabled: false,
      metrics: [
        DeviceMetric(label: '状态', value: '离线'),
        DeviceMetric(label: '信号', value: '无'),
        DeviceMetric(label: '负载', value: '--'),
      ],
    ),
  ];
}
