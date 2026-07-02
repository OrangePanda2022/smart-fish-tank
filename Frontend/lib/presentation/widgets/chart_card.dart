import 'package:flutter/cupertino.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../view_models/chart_series.dart';
import 'painters.dart';
import 'soft_card.dart';

class ChartCard extends StatelessWidget {
  const ChartCard({required this.series, super.key});

  final ChartSeries series;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      radius: 12,
      padding: const EdgeInsets.all(18),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Text(series.title, style: AppText.titleSmall),
              const SizedBox(width: 6),
              const Icon(
                CupertinoIcons.question_circle,
                color: AppColors.secondaryText,
                size: 18,
              ),
              const Spacer(),
              const Icon(
                CupertinoIcons.arrow_up_left_arrow_down_right,
                color: AppColors.secondaryText,
                size: 22,
              ),
            ],
          ),
          const SizedBox(height: 26),
          SizedBox(
            height: 300,
            child: CustomPaint(
              painter: HistoryChartPainter(series: series),
              child: const SizedBox.expand(),
            ),
          ),
          const SizedBox(height: 18),
          ChartStats(series: series),
        ],
      ),
    );
  }
}

class ChartStats extends StatelessWidget {
  const ChartStats({required this.series, super.key});

  final ChartSeries series;

  @override
  Widget build(BuildContext context) {
    if (series.isEmpty) {
      return Container(
        width: double.infinity,
        padding: const EdgeInsets.symmetric(vertical: 18),
        decoration: BoxDecoration(
          color: AppColors.paleBlue.withAlpha(140),
          borderRadius: BorderRadius.circular(14),
        ),
        child: const Text(
          '暂无统计数据',
          textAlign: TextAlign.center,
          style: AppText.caption,
        ),
      );
    }

    final unit = series.unit;
    final items = [
      ('平均值', '${series.average.toStringAsFixed(series.precision)}$unit'),
      ('最高值', '${series.max.toStringAsFixed(series.precision)}$unit'),
      ('最低值', '${series.min.toStringAsFixed(series.precision)}$unit'),
      (
        '波动范围',
        '±${(series.max - series.min).toStringAsFixed(series.precision)}$unit',
      ),
    ];

    return Container(
      padding: const EdgeInsets.symmetric(vertical: 18),
      decoration: BoxDecoration(
        color: AppColors.paleBlue.withAlpha(140),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Row(
        children: items
            .map(
              (item) => Expanded(
                child: Column(
                  children: [
                    Text(item.$1, style: AppText.caption),
                    const SizedBox(height: 10),
                    FittedBox(
                      fit: BoxFit.scaleDown,
                      child: Text(item.$2, style: AppText.statValue),
                    ),
                  ],
                ),
              ),
            )
            .toList(),
      ),
    );
  }
}
