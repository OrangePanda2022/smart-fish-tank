import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/sensor_reading.dart';
import '../view_models/sensor_metric.dart';
import 'soft_card.dart';

class SensorMetricGrid extends StatelessWidget {
  const SensorMetricGrid({required this.sensor, super.key});

  final SensorReading sensor;

  @override
  Widget build(BuildContext context) {
    final metrics = [
      SensorMetric(
        context.tr('pH值'),
        sensor.ph.toStringAsFixed(2),
        '',
        CupertinoIcons.drop,
        AppColors.teal,
      ),
      SensorMetric(
        context.tr('温度'),
        sensor.temperature.toStringAsFixed(1),
        '°C',
        CupertinoIcons.thermometer,
        AppColors.orange,
      ),
      SensorMetric(
        'TDS',
        sensor.tds.toStringAsFixed(0),
        'ppm',
        CupertinoIcons.square_list,
        AppColors.mint,
      ),
      SensorMetric(
        context.tr('溶氧'),
        sensor.oxygen.toStringAsFixed(1),
        'mg/L',
        CupertinoIcons.drop_fill,
        AppColors.primary,
      ),
      const SensorMetric(
        'ORP',
        '152',
        'mV',
        CupertinoIcons.wand_stars,
        AppColors.green,
      ),
      SensorMetric(
        context.tr('氨氮'),
        sensor.ammonia.toStringAsFixed(2),
        'mg/L',
        CupertinoIcons.lab_flask,
        AppColors.purple,
      ),
      SensorMetric(
        context.tr('亚硝酸盐'),
        sensor.nitrite.toStringAsFixed(2),
        'mg/L',
        CupertinoIcons.drop_triangle,
        AppColors.violet,
      ),
      SensorMetric(
        context.tr('硝酸盐'),
        sensor.nitrate.toStringAsFixed(1),
        'mg/L',
        CupertinoIcons.square_stack,
        AppColors.mint,
      ),
      SensorMetric(
        context.tr('氯离子'),
        sensor.chloride.toStringAsFixed(1),
        'mg/L',
        CupertinoIcons.circle_grid_3x3,
        AppColors.periwinkle,
      ),
    ];

    return LayoutBuilder(
      builder: (context, constraints) {
        final isNarrow = constraints.maxWidth < 340;

        return GridView.builder(
          padding: EdgeInsets.zero,
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          itemCount: metrics.length,
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: isNarrow ? 2 : 3,
            mainAxisSpacing: 10,
            crossAxisSpacing: 10,
            childAspectRatio: isNarrow ? 1.05 : 0.86,
          ),
          itemBuilder: (context, index) =>
              SensorMetricTile(metric: metrics[index], compact: isNarrow),
        );
      },
    );
  }
}

class SensorMetricTile extends StatelessWidget {
  const SensorMetricTile({
    required this.metric,
    this.compact = false,
    super.key,
  });

  final SensorMetric metric;
  final bool compact;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      radius: 12,
      padding: EdgeInsets.symmetric(
        horizontal: compact ? 10 : 8,
        vertical: compact ? 16 : 24,
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Icon(metric.icon, color: metric.color, size: 25),
          Text(metric.label, style: AppText.metricLabel),
          FittedBox(
            fit: BoxFit.scaleDown,
            child: RichText(
              text: TextSpan(
                style: const TextStyle(
                  color: AppColors.ink,
                  fontWeight: FontWeight.w900,
                ),
                children: [
                  TextSpan(
                    text: metric.value,
                    style: const TextStyle(fontSize: 25),
                  ),
                  if (metric.unit.isNotEmpty)
                    TextSpan(
                      text: ' ${metric.unit}',
                      style: const TextStyle(fontSize: 12),
                    ),
                ],
              ),
            ),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const CircleAvatar(radius: 4, backgroundColor: AppColors.success),
              const SizedBox(width: 5),
              Text(context.tr('正常'), style: AppText.okText),
            ],
          ),
        ],
      ),
    );
  }
}
