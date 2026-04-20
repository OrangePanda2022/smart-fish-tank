import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:aqua/domain/models/tank_history.dart';

class HistoryChart extends StatelessWidget {
  final TankHistoryData data;

  const HistoryChart({super.key, required this.data});

  @override
  Widget build(BuildContext context) {
    final lineColor = Theme.of(context).colorScheme.primary;
    final normalMin = data.type.normalMin;
    final normalMax = data.type.normalMax;

    return Container(
      height: 300,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surfaceContainerHighest,
        borderRadius: BorderRadius.circular(32),
      ),
      child: LineChart(
        LineChartData(
          gridData: FlGridData(
            show: true,
            drawVerticalLine: false,
            horizontalInterval: _calculateInterval(normalMin, normalMax),
            getDrawingHorizontalLine: (value) {
              return FlLine(
                color: Theme.of(context).colorScheme.outline.withAlpha(51),
                strokeWidth: 1,
              );
            },
          ),
          titlesData: FlTitlesData(
            leftTitles: AxisTitles(
              sideTitles: SideTitles(
                showTitles: true,
                reservedSize: 40,
                interval: _calculateInterval(normalMin, normalMax),
                getTitlesWidget: (value, meta) {
                  return Text(
                    value.toStringAsFixed(1),
                    style: TextStyle(
                      fontSize: 12,
                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                    ),
                  );
                },
              ),
            ),
            bottomTitles: AxisTitles(
              sideTitles: SideTitles(
                showTitles: true,
                reservedSize: 30,
                interval: 7,
                getTitlesWidget: (value, meta) {
                  final index = value.toInt();
                  if (index >= 0 && index < data.dataPoints.length) {
                    final date = data.dataPoints[index].timestamp;
                    return Padding(
                      padding: const EdgeInsets.only(top: 8),
                      child: Text(
                        '${date.month}/${date.day}',
                        style: TextStyle(
                          fontSize: 10,
                          color: Theme.of(context).colorScheme.onSurfaceVariant,
                        ),
                      ),
                    );
                  }
                  return const SizedBox();
                },
              ),
            ),
            topTitles: const AxisTitles(
              sideTitles: SideTitles(showTitles: false),
            ),
            rightTitles: const AxisTitles(
              sideTitles: SideTitles(showTitles: false),
            ),
          ),
          borderData: FlBorderData(show: false),
          minY: normalMin - 2,
          maxY: normalMax + 2,
          lineBarsData: [
            LineChartBarData(
              spots: data.dataPoints.asMap().entries.map((entry) {
                return FlSpot(entry.key.toDouble(), entry.value.value);
              }).toList(),
              isCurved: true,
              curveSmoothness: 0.3,
              color: lineColor,
              barWidth: 3,
              dotData: const FlDotData(show: false),
              belowBarData: BarAreaData(
                show: true,
                color: lineColor.withAlpha(51),
              ),
            ),
          ],
          lineTouchData: LineTouchData(
            touchTooltipData: LineTouchTooltipData(
              getTooltipColor: (touchedSpot) =>
                  Theme.of(context).colorScheme.primaryContainer,
              fitInsideHorizontally: true,
              fitInsideVertically: true,
              getTooltipItems: (touchedSpots) {
                return touchedSpots.map((spot) {
                  final index = spot.x.toInt();
                  final date = data.dataPoints[index].timestamp;
                  return LineTooltipItem(
                    '${date.month}/${date.day} ${date.hour}:00\n',
                    TextStyle(
                      color: Theme.of(context).colorScheme.onPrimaryContainer,
                      fontWeight: FontWeight.bold,
                    ),
                    children: [
                      TextSpan(
                        text: '${spot.y.toStringAsFixed(1)}${data.unit}',
                        style: TextStyle(
                          color: Theme.of(
                            context,
                          ).colorScheme.onPrimaryContainer,
                          fontWeight: FontWeight.normal,
                        ),
                      ),
                    ],
                  );
                }).toList();
              },
            ),
          ),
          extraLinesData: ExtraLinesData(
            horizontalLines: [
              HorizontalLine(
                y: normalMin,
                color: Colors.green.withAlpha(128),
                strokeWidth: 2,
                dashArray: [8, 4],
                label: HorizontalLineLabel(
                  show: true,
                  alignment: Alignment.topRight,
                  padding: const EdgeInsets.only(right: 8, bottom: 4),
                  style: const TextStyle(fontSize: 10, color: Colors.green),
                  labelResolver: (line) =>
                      'Min: ${normalMin.toStringAsFixed(1)}',
                ),
              ),
              HorizontalLine(
                y: normalMax,
                color: Colors.green.withAlpha(128),
                strokeWidth: 2,
                dashArray: [8, 4],
                label: HorizontalLineLabel(
                  show: true,
                  alignment: Alignment.bottomRight,
                  padding: const EdgeInsets.only(right: 8, top: 4),
                  style: const TextStyle(fontSize: 10, color: Colors.green),
                  labelResolver: (line) =>
                      'Max: ${normalMax.toStringAsFixed(1)}',
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  double _calculateInterval(double min, double max) {
    final range = max - min;
    if (range <= 2) return 0.5;
    if (range <= 10) return 2;
    return 5;
  }
}
