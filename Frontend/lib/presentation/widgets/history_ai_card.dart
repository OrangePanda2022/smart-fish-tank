import 'package:flutter/cupertino.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../view_models/chart_series.dart';
import 'soft_card.dart';

class HistoryAiCard extends StatelessWidget {
  const HistoryAiCard({required this.series, super.key});

  final ChartSeries series;

  @override
  Widget build(BuildContext context) {
    final isTemperature = series.metric == 'temperature';
    return SoftCard(
      radius: 12,
      padding: const EdgeInsets.fromLTRB(18, 20, 18, 18),
      child: Stack(
        children: [
          Positioned(
            right: -4,
            top: -20,
            child: SizedBox(
              width: 74,
              height: 82,
              child: Image.asset('assets/images/ai_robot.png'),
            ),
          ),
          Padding(
            padding: const EdgeInsets.only(right: 54),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('✦ AI 分析', style: AppText.blueHeading),
                const SizedBox(height: 18),
                Text(
                  isTemperature
                      ? '近7天温度在${series.min.toStringAsFixed(1)}°C ~ ${series.max.toStringAsFixed(1)}°C之间波动，整体趋于稳定。\n\n每日12:00 ~ 16:00温度较高，可能影响鱼缸水质。'
                      : '近7天 pH 值整体保持稳定，波动范围在正常区间内。\n\n在05/15下午出现轻微上升趋势，可能与换水或喂食有关。',
                  style: AppText.body,
                ),
                const SizedBox(height: 18),
                const Text('建议：', style: AppText.blueHeading),
                const SizedBox(height: 10),
                AdviceRow(
                  text: isTemperature ? '建议在12:00 ~ 16:00开启散热风扇' : '保持当前换水频率',
                ),
                const SizedBox(height: 8),
                AdviceRow(
                  text: isTemperature
                      ? '可将目标温度设置为25 ~ 26°C'
                      : '可在每日16:00后适当增加曝气',
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class AdviceRow extends StatelessWidget {
  const AdviceRow({required this.text, super.key});

  final String text;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Icon(
          CupertinoIcons.check_mark_circled,
          color: AppColors.success,
          size: 18,
        ),
        const SizedBox(width: 8),
        Expanded(child: Text(text, style: AppText.body)),
      ],
    );
  }
}
