import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../domain/entities/sensor_reading.dart';
import 'device_control.dart';

List<DeviceControl> buildDeviceCatalog(
  BuildContext context,
  SensorReading sensor,
) {
  return [
    DeviceControl(
      context.tr('灯光'),
      'ON',
      CupertinoIcons.lightbulb_fill,
      AppColors.yellow,
      id: 'light',
      description: context.tr('模拟自然昼夜节律，当前处于日间补光。'),
      primaryAction: context.tr('调节亮度'),
      maintenance: context.tr('建议每天 08:00-20:00 开启'),
      metrics: [
        DeviceMetric(label: context.tr('亮度'), value: '76%'),
        DeviceMetric(label: context.tr('模式'), value: context.tr('日光')),
        DeviceMetric(label: context.tr('定时'), value: context.tr('20:00 关闭')),
      ],
    ),
    DeviceControl(
      context.tr('水泵'),
      'ON',
      Icons.waves_rounded,
      AppColors.sky,
      id: 'pump',
      description: context.tr('保持水体循环，辅助过滤系统稳定工作。'),
      primaryAction: context.tr('调整流量'),
      maintenance: context.tr('滤芯建议 14 天后检查'),
      metrics: [
        DeviceMetric(label: context.tr('流量'), value: context.tr('中档')),
        DeviceMetric(
          label: context.tr('水位'),
          value: '${sensor.waterLevel.toStringAsFixed(0)}%',
        ),
        DeviceMetric(label: context.tr('运行'), value: context.tr('连续')),
      ],
    ),
    DeviceControl(
      context.tr('氧气'),
      'ON',
      CupertinoIcons.circle_grid_hex,
      AppColors.bubble,
      id: 'aerator',
      description: context.tr('根据溶氧水平维持鱼缸含氧量。'),
      primaryAction: context.tr('增强供氧'),
      maintenance: context.tr('气石状态良好'),
      metrics: [
        DeviceMetric(
          label: context.tr('溶氧'),
          value: sensor.oxygen.toStringAsFixed(1),
        ),
        DeviceMetric(label: context.tr('强度'), value: '65%'),
        DeviceMetric(label: context.tr('模式'), value: context.tr('智能')),
      ],
    ),
    DeviceControl(
      context.tr('加热'),
      '${sensor.temperature.toStringAsFixed(1)}°C',
      CupertinoIcons.thermometer,
      AppColors.orange,
      id: 'heater',
      description: context.tr('自动维持适合观赏鱼的水温区间。'),
      primaryAction: context.tr('设置温度'),
      maintenance: context.tr('目标温度 26.0°C'),
      metrics: [
        DeviceMetric(
          label: context.tr('当前水温'),
          value: '${sensor.temperature.toStringAsFixed(1)}°C',
        ),
        DeviceMetric(label: context.tr('目标'), value: '26.0°C'),
        DeviceMetric(label: context.tr('功率'), value: '42%'),
      ],
    ),
    DeviceControl(
      context.tr('喂食'),
      context.tr('08:00 已喂食'),
      Icons.fastfood_rounded,
      AppColors.feed,
      id: 'feeder',
      description: context.tr('自动投喂已按计划完成。'),
      primaryAction: context.tr('立即投喂'),
      maintenance: context.tr('饲料余量充足'),
      metrics: [
        DeviceMetric(label: context.tr('今日次数'), value: '1/2'),
        DeviceMetric(label: context.tr('下次'), value: '18:00'),
        DeviceMetric(label: context.tr('份量'), value: context.tr('标准')),
      ],
    ),
    DeviceControl(
      context.tr('换水'),
      context.tr('剩余 7 天'),
      CupertinoIcons.drop_fill,
      AppColors.primary,
      id: 'water_change',
      description: context.tr('记录换水周期，提醒维持稳定水质。'),
      primaryAction: context.tr('记录换水'),
      maintenance: context.tr('下次换水 7 天后'),
      metrics: [
        DeviceMetric(label: context.tr('周期'), value: context.tr('14 天')),
        DeviceMetric(label: context.tr('上次'), value: context.tr('7 天前')),
        DeviceMetric(label: context.tr('建议量'), value: '20%'),
      ],
    ),
    DeviceControl(
      context.tr('摄像头'),
      context.tr('在线'),
      CupertinoIcons.camera_fill,
      AppColors.camera,
      id: 'camera',
      description: context.tr('远程查看鱼缸画面和设备状态。'),
      primaryAction: context.tr('查看画面'),
      maintenance: context.tr('网络连接稳定'),
      metrics: [
        DeviceMetric(label: context.tr('清晰度'), value: '1080P'),
        DeviceMetric(label: context.tr('延迟'), value: context.tr('低')),
        DeviceMetric(label: context.tr('存储'), value: context.tr('本地')),
      ],
    ),
    DeviceControl(
      context.tr('过滤'),
      'ON',
      CupertinoIcons.slider_horizontal_3,
      AppColors.teal,
      id: 'filter',
      description: context.tr('过滤杂质并辅助稳定氨氮、硝酸盐指标。'),
      primaryAction: context.tr('切换模式'),
      maintenance: context.tr('建议 3 天后清洗滤棉'),
      metrics: [
        DeviceMetric(label: context.tr('模式'), value: context.tr('自动')),
        DeviceMetric(label: context.tr('氨氮'), value: context.tr('稳定')),
        DeviceMetric(label: context.tr('滤棉'), value: context.tr('良好')),
      ],
    ),
    DeviceControl(
      context.tr('备用插座'),
      context.tr('离线'),
      Icons.power_rounded,
      AppColors.muted,
      id: 'outlet',
      description: context.tr('可接入扩展设备，当前未连接。'),
      primaryAction: context.tr('重新连接'),
      maintenance: context.tr('检查电源与网络'),
      isOnline: false,
      isEnabled: false,
      metrics: [
        DeviceMetric(label: context.tr('状态'), value: context.tr('离线')),
        DeviceMetric(label: context.tr('信号'), value: context.tr('无')),
        DeviceMetric(label: context.tr('负载'), value: '--'),
      ],
    ),
  ];
}
