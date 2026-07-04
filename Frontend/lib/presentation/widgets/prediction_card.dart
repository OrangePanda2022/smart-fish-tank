import 'package:flutter/cupertino.dart';

import '../../core/i18n/app_localizations.dart';
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
    final suggestion = _suggestionText(context, dashboard);

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
                      children: [
                        Text(context.tr('AI 趋势预测'), style: AppText.blueHeading),
                        const Spacer(),
                        Text(context.tr('未来24小时'), style: AppText.caption),
                      ],
                    ),
                    const SizedBox(height: 18),
                    TrendRow(
                      label: context.tr('pH值'),
                      value:
                          '${sensor.ph.toStringAsFixed(1)} → ${forecast.ph.toStringAsFixed(1)}',
                      color: AppColors.success,
                    ),
                    TrendRow(
                      label: context.tr('温度'),
                      value:
                          '${sensor.temperature.toStringAsFixed(1)}°C → ${forecast.temperature.toStringAsFixed(1)}°C',
                      color: AppColors.orange,
                    ),
                    TrendRow(
                      label: context.tr('溶氧'),
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

  String _suggestionText(BuildContext context, AquariumDashboard dashboard) {
    if (dashboard.prediction.hasWarning) {
      return '${context.tr('AI 建议')}\n${dashboard.prediction.warning}';
    }

    final actions = dashboard.prediction.actions
        .where((action) => action.action != 'maintain')
        .map(
          (action) =>
              '${_deviceName(context, action.device)}${_actionName(context, action.action)}',
        )
        .take(2)
        .join(context.l10n.isEnglish ? ', ' : '，');
    if (actions.isNotEmpty) {
      return context.l10n.isEnglish
          ? '${context.tr('AI 建议')}\nRecommend $actions'
          : '${context.tr('AI 建议')}\n建议$actions';
    }

    return '${context.tr('AI 建议')}\n${context.tr('未来趋势稳定，保持当前设备策略')}';
  }

  String _deviceName(BuildContext context, String device) {
    return switch (device) {
      'heater' => context.tr('加热设备'),
      'aerator' => context.tr('增氧设备'),
      'pump' => context.tr('水泵'),
      'feeder' => context.tr('喂食器'),
      'light' => context.tr('灯光'),
      _ => device,
    };
  }

  String _actionName(BuildContext context, String action) {
    return switch (action) {
      'cool' => context.tr('降温'),
      'heat' => context.tr('升温'),
      'on' => context.tr('开启'),
      'off' => context.tr('关闭'),
      'feed' => context.l10n.isEnglish ? ' feed' : '喂食',
      'dim' => context.tr('调暗'),
      'brighten' => context.tr('调亮'),
      'hold' => context.tr('保持'),
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
