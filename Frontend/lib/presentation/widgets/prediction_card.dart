import 'package:flutter/cupertino.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import 'painters.dart';
import 'soft_card.dart';

class PredictionCard extends StatelessWidget {
  const PredictionCard({required this.dashboard, super.key});

  final AquariumDashboard dashboard;

  @override
  Widget build(BuildContext context) {
    final sensor = dashboard.latest;
    final forecast = dashboard.prediction.forecast;
    final suggestion = _suggestionText(dashboard);

    return SoftCard(
      padding: const EdgeInsets.fromLTRB(18, 18, 14, 16),
      child: Column(
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: const [
                        Text('AI 趋势预测', style: AppText.blueHeading),
                        Spacer(),
                        Text('未来24小时', style: AppText.caption),
                      ],
                    ),
                    const SizedBox(height: 18),
                    TrendRow(
                      label: 'pH值',
                      value:
                          '${sensor.ph.toStringAsFixed(1)} → ${forecast.ph.toStringAsFixed(1)}',
                      color: AppColors.success,
                    ),
                    TrendRow(
                      label: '温度',
                      value:
                          '${sensor.temperature.toStringAsFixed(1)}°C → ${forecast.temperature.toStringAsFixed(1)}°C',
                      color: AppColors.orange,
                    ),
                    TrendRow(
                      label: '溶氧',
                      value:
                          '${sensor.oxygen.toStringAsFixed(1)}mg/L → ${forecast.oxygen.toStringAsFixed(1)}mg/L',
                      color: AppColors.primary,
                    ),
                  ],
                ),
              ),
              SizedBox(
                width: 95,
                height: 105,
                child: Image.asset(
                  'assets/images/ai_robot.png',
                  fit: BoxFit.contain,
                ),
              ),
            ],
          ),
          const SizedBox(height: 14),
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(14),
            decoration: BoxDecoration(
              color: AppColors.paleBlue,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Row(
              children: [
                const Icon(CupertinoIcons.lightbulb, color: AppColors.primary),
                const SizedBox(width: 10),
                Expanded(
                  child: Text(
                    suggestion,
                    style: const TextStyle(
                      color: AppColors.ink,
                      fontSize: 14,
                      height: 1.45,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  String _suggestionText(AquariumDashboard dashboard) {
    if (dashboard.prediction.hasWarning) {
      return 'AI 建议\n${dashboard.prediction.warning}';
    }

    final actions = dashboard.prediction.actions
        .where((action) => action.action != 'maintain')
        .map(
          (action) =>
              '${_deviceName(action.device)}${_actionName(action.action)}',
        )
        .take(2)
        .join('，');
    if (actions.isNotEmpty) return 'AI 建议\n建议$actions';

    return 'AI 建议\n未来趋势稳定，保持当前设备策略';
  }

  String _deviceName(String device) {
    return switch (device) {
      'heater' => '加热设备',
      'aerator' => '增氧设备',
      'pump' => '水泵',
      'feeder' => '喂食器',
      'light' => '灯光',
      _ => device,
    };
  }

  String _actionName(String action) {
    return switch (action) {
      'cool' => '降温',
      'heat' => '升温',
      'on' => '开启',
      'off' => '关闭',
      'feed' => '喂食',
      'dim' => '调暗',
      'brighten' => '调亮',
      'hold' => '保持',
      _ => action,
    };
  }
}

class TrendRow extends StatelessWidget {
  const TrendRow({
    required this.label,
    required this.value,
    required this.color,
    super.key,
  });

  final String label;
  final String value;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Row(
        children: [
          SizedBox(width: 58, child: Text(label, style: AppText.metricLabel)),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                color: AppColors.ink,
                fontSize: 15,
                fontWeight: FontWeight.w800,
              ),
            ),
          ),
          SizedBox(
            width: 58,
            height: 20,
            child: CustomPaint(painter: SparklinePainter(color: color)),
          ),
        ],
      ),
    );
  }
}
