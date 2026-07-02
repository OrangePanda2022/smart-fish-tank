import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_shadow.dart';
import '../../core/theme/app_text.dart';
import '../view_models/chart_series.dart';
import 'circle_icon_button.dart';

class HistoryFilters extends StatelessWidget {
  const HistoryFilters({
    required this.metric,
    required this.rangeIndex,
    required this.customRangeLabel,
    required this.onMetricChanged,
    required this.onRangeChanged,
    required this.onCustomRangePressed,
    super.key,
  });

  final String metric;
  final int rangeIndex;
  final String? customRangeLabel;
  final ValueChanged<String> onMetricChanged;
  final ValueChanged<int> onRangeChanged;
  final VoidCallback onCustomRangePressed;

  @override
  Widget build(BuildContext context) {
    final ranges = ['今天', '近7天', '近30天', '自定义'];
    final selectedMetric = ChartSeries.metricFor(metric);
    final rangeLabel = rangeIndex == 3 && customRangeLabel != null
        ? customRangeLabel!
        : ranges[rangeIndex];

    return Column(
      children: [
        Row(
          children: [
            Expanded(
              child: FilterSelect(
                label: '传感器',
                value: selectedMetric.label,
                icon: _iconForMetric(metric),
                onTap: () => _showMetricPicker(context),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: FilterSelect(
                label: '时间范围',
                value: rangeLabel,
                icon: CupertinoIcons.calendar,
                onTap: onCustomRangePressed,
              ),
            ),
            const SizedBox(width: 8),
            CircleIconButton(
              icon: CupertinoIcons.calendar,
              fill: Colors.white,
              color: AppColors.ink,
              onPressed: onCustomRangePressed,
            ),
          ],
        ),
        const SizedBox(height: 18),
        Row(
          children: List.generate(
            ranges.length,
            (index) => Expanded(
              child: Padding(
                padding: EdgeInsets.only(
                  right: index == ranges.length - 1 ? 0 : 8,
                ),
                child: ChoiceChipButton(
                  text: ranges[index],
                  selected: rangeIndex == index,
                  onTap: () {
                    if (index == 3) {
                      onCustomRangePressed();
                    } else {
                      onRangeChanged(index);
                    }
                  },
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  void _showMetricPicker(BuildContext context) {
    showModalBottomSheet<void>(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(18)),
      ),
      builder: (context) {
        return SafeArea(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(18, 12, 18, 18),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Container(
                  width: 42,
                  height: 4,
                  decoration: BoxDecoration(
                    color: AppColors.grid,
                    borderRadius: BorderRadius.circular(99),
                  ),
                ),
                const SizedBox(height: 14),
                ...ChartSeries.metrics.map(
                  (option) => ListTile(
                    leading: Icon(
                      _iconForMetric(option.key),
                      color: option.color,
                    ),
                    title: Text(option.label, style: AppText.cardTitle),
                    trailing: option.key == metric
                        ? const Icon(
                            CupertinoIcons.check_mark,
                            color: AppColors.primary,
                          )
                        : null,
                    onTap: () {
                      Navigator.of(context).pop();
                      onMetricChanged(option.key);
                    },
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  IconData _iconForMetric(String metric) {
    return switch (metric) {
      'temperature' => CupertinoIcons.thermometer,
      'oxygen' => CupertinoIcons.wind,
      'tds' => CupertinoIcons.waveform_path_ecg,
      'ammonia' => CupertinoIcons.exclamationmark_triangle,
      'nitrate' || 'nitrite' => CupertinoIcons.lab_flask,
      'chloride' => CupertinoIcons.drop_triangle,
      'waterLevel' => CupertinoIcons.gauge,
      _ => CupertinoIcons.drop_fill,
    };
  }
}

class FilterSelect extends StatelessWidget {
  const FilterSelect({
    required this.label,
    required this.value,
    required this.icon,
    required this.onTap,
    super.key,
  });

  final String label;
  final String value;
  final IconData icon;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: AppText.caption),
        const SizedBox(height: 10),
        GestureDetector(
          onTap: onTap,
          child: Container(
            height: 62,
            padding: const EdgeInsets.symmetric(horizontal: 14),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(12),
              boxShadow: AppShadow.soft,
            ),
            child: Row(
              children: [
                Icon(icon, color: AppColors.primary, size: 21),
                const SizedBox(width: 8),
                Expanded(child: Text(value, style: AppText.cardTitle)),
                const Icon(
                  CupertinoIcons.chevron_down,
                  color: AppColors.secondaryText,
                  size: 18,
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}

class ChoiceChipButton extends StatelessWidget {
  const ChoiceChipButton({
    required this.text,
    required this.selected,
    required this.onTap,
    super.key,
  });

  final String text;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 180),
        height: 48,
        alignment: Alignment.center,
        decoration: BoxDecoration(
          color: selected ? AppColors.primary : Colors.white,
          borderRadius: BorderRadius.circular(14),
          boxShadow: selected ? AppShadow.button : AppShadow.soft,
        ),
        child: Text(
          text,
          style: TextStyle(
            color: selected ? Colors.white : AppColors.ink,
            fontSize: 15,
            fontWeight: FontWeight.w700,
          ),
        ),
      ),
    );
  }
}
