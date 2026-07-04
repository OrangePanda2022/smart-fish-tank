import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/entities/sensor_reading.dart';

class AquariumStatusCard extends StatelessWidget {
  const AquariumStatusCard({required this.dashboard, super.key});

  final AquariumDashboard dashboard;

  @override
  Widget build(BuildContext context) {
    final sensor = dashboard.latest;
    final score = dashboard.score;

    return ClipRRect(
      borderRadius: BorderRadius.circular(26),
      child: SizedBox(
        height: 290,
        child: Stack(
          fit: StackFit.expand,
          children: [
            Image.asset('assets/images/aquarium_hero.png', fit: BoxFit.cover),
            DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    AppColors.navy.withAlpha(125),
                    Colors.transparent,
                    AppColors.navy.withAlpha(180),
                  ],
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(24, 22, 24, 20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          score == null
                              ? context.tr('鱼缸状态  未分析')
                              : context.tr('鱼缸状态  良好'),
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 23,
                            fontWeight: FontWeight.w800,
                          ),
                        ),
                      ),
                      Container(
                        width: 10,
                        height: 10,
                        decoration: const BoxDecoration(
                          color: AppColors.success,
                          shape: BoxShape.circle,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        context.tr('设备在线'),
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 16,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 18),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Flexible(
                        child: FittedBox(
                          fit: BoxFit.scaleDown,
                          alignment: Alignment.centerLeft,
                          child: Text(
                            score == null ? context.tr('未分析') : '$score',
                            maxLines: 1,
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 64,
                              height: 0.95,
                              fontWeight: FontWeight.w900,
                            ),
                          ),
                        ),
                      ),
                      if (score != null) ...[
                        const SizedBox(width: 8),
                        Padding(
                          padding: const EdgeInsets.only(bottom: 8),
                          child: Text(
                            context.tr('分'),
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 22,
                              fontWeight: FontWeight.w800,
                            ),
                          ),
                        ),
                      ],
                    ],
                  ),
                  const SizedBox(height: 8),
                  score == null
                      ? Text(
                          context.tr('点击立即分析生成健康评分'),
                          style: TextStyle(
                            color: Colors.white.withAlpha(215),
                            fontSize: 15,
                            fontWeight: FontWeight.w700,
                          ),
                        )
                      : const StarRating(),
                  const Spacer(),
                  StatusMetricStrip(sensor: sensor),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class StarRating extends StatelessWidget {
  const StarRating({super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: List.generate(
        5,
        (index) => Icon(
          index == 4 ? Icons.star_half_rounded : Icons.star_rounded,
          color: AppColors.star,
          size: 25,
        ),
      ),
    );
  }
}

class StatusMetricStrip extends StatelessWidget {
  const StatusMetricStrip({required this.sensor, super.key});

  final SensorReading sensor;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 70,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        color: AppColors.navy.withAlpha(145),
        borderRadius: BorderRadius.circular(22),
      ),
      child: Row(
        children: [
          StatusMetric(
            icon: CupertinoIcons.thermometer,
            value: '${sensor.temperature.toStringAsFixed(1)}°C',
            label: context.tr('水温'),
          ),
          StatusMetric(
            icon: CupertinoIcons.drop_fill,
            value: sensor.ph.toStringAsFixed(1),
            label: 'pH',
          ),
          StatusMetric(
            icon: CupertinoIcons.circle_grid_hex_fill,
            value: '${sensor.oxygen.toStringAsFixed(1)} mg/L',
            label: context.tr('溶氧'),
          ),
          StatusMetric(
            icon: CupertinoIcons.waveform_path_ecg,
            value: context.tr('已喂食'),
            label: context.tr('今天 08:00'),
          ),
        ],
      ),
    );
  }
}

class StatusMetric extends StatelessWidget {
  const StatusMetric({
    required this.icon,
    required this.value,
    required this.label,
    super.key,
  });

  final IconData icon;
  final String value;
  final String label;

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Row(
        children: [
          Container(
            width: 38,
            height: 38,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: Colors.white.withAlpha(42),
              shape: BoxShape.circle,
            ),
            child: Icon(icon, color: Colors.white, size: 20),
          ),
          const SizedBox(width: 8),
          Flexible(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                FittedBox(
                  fit: BoxFit.scaleDown,
                  alignment: Alignment.centerLeft,
                  child: Text(
                    value,
                    maxLines: 1,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 17,
                      fontWeight: FontWeight.w900,
                    ),
                  ),
                ),
                const SizedBox(height: 3),
                Text(
                  label,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(
                    color: Colors.white.withAlpha(190),
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
