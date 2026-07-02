import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../view_models/chart_series.dart';

class HistoryChartPainter extends CustomPainter {
  HistoryChartPainter({required this.series});

  final ChartSeries series;

  @override
  void paint(Canvas canvas, Size size) {
    const left = 38.0;
    const top = 8.0;
    const right = 8.0;
    const bottom = 34.0;
    final chart = Rect.fromLTRB(
      left,
      top,
      size.width - right,
      size.height - bottom,
    );
    final gridPaint = Paint()
      ..color = AppColors.grid
      ..strokeWidth = 1;
    final textPainter = TextPainter(
      textDirection: TextDirection.ltr,
      textAlign: TextAlign.right,
    );

    for (var i = 0; i <= 4; i++) {
      final y = chart.top + chart.height * i / 4;
      canvas.drawLine(Offset(chart.left, y), Offset(chart.right, y), gridPaint);
      final value = series.high - (series.high - series.low) * i / 4;
      textPainter.text = TextSpan(
        text: value.toStringAsFixed(series.axisPrecision),
        style: const TextStyle(color: AppColors.secondaryText, fontSize: 11),
      );
      textPainter.layout(maxWidth: 34);
      textPainter.paint(canvas, Offset(0, y - 7));
    }

    for (var i = 0; i <= 6; i++) {
      final x = chart.left + chart.width * i / 6;
      canvas.drawLine(Offset(x, chart.top), Offset(x, chart.bottom), gridPaint);
      final timestamp = _timestampAt(i / 6);
      textPainter.text = TextSpan(
        text: _formatAxisDate(timestamp),
        style: const TextStyle(color: AppColors.secondaryText, fontSize: 11),
      );
      textPainter.layout();
      textPainter.paint(
        canvas,
        Offset(x - textPainter.width / 2, chart.bottom + 12),
      );
    }

    final targetPaint = Paint()
      ..color = series.metric == 'temperature'
          ? AppColors.alert.withAlpha(125)
          : AppColors.alert.withAlpha(125)
      ..strokeWidth = 1.2;
    final lowLine = _mapY(series.safeLow, chart);
    final highLine = _mapY(series.safeHigh, chart);
    _drawDashed(
      canvas,
      Offset(chart.left, lowLine),
      Offset(chart.right, lowLine),
      targetPaint,
    );
    _drawDashed(
      canvas,
      Offset(chart.left, highLine),
      Offset(chart.right, highLine),
      targetPaint,
    );

    if (series.values.isEmpty) {
      textPainter.text = const TextSpan(
        text: '暂无匹配数据',
        style: TextStyle(
          color: AppColors.secondaryText,
          fontSize: 14,
          fontWeight: FontWeight.w700,
        ),
      );
      textPainter.textAlign = TextAlign.center;
      textPainter.layout(maxWidth: chart.width);
      textPainter.paint(
        canvas,
        Offset(
          chart.left + (chart.width - textPainter.width) / 2,
          chart.top + (chart.height - textPainter.height) / 2,
        ),
      );
      return;
    }

    final fillPath = Path();
    final linePath = Path();
    for (var i = 0; i < series.values.length; i++) {
      final x = _mapX(i, chart);
      final y = _mapY(series.values[i], chart);
      if (i == 0) {
        linePath.moveTo(x, y);
        fillPath.moveTo(x, chart.bottom);
        fillPath.lineTo(x, y);
      } else {
        linePath.lineTo(x, y);
        fillPath.lineTo(x, y);
      }
    }
    fillPath
      ..lineTo(chart.right, chart.bottom)
      ..close();

    final fillPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [series.color.withAlpha(50), series.color.withAlpha(0)],
      ).createShader(chart);
    canvas.drawPath(fillPath, fillPaint);

    final linePaint = Paint()
      ..color = series.color
      ..strokeWidth = 2
      ..style = PaintingStyle.stroke
      ..strokeCap = StrokeCap.round
      ..strokeJoin = StrokeJoin.round;
    canvas.drawPath(linePath, linePaint);

    final highlightIndex = math.min(
      series.values.length - 1,
      series.values.length ~/ 2,
    );
    final hx = _mapX(highlightIndex, chart);
    final hy = _mapY(series.values[highlightIndex], chart);
    canvas.drawCircle(Offset(hx, hy), 4, Paint()..color = series.color);
    canvas.drawLine(
      Offset(hx, chart.top),
      Offset(hx, chart.bottom),
      Paint()
        ..color = AppColors.secondaryText.withAlpha(55)
        ..strokeWidth = 1,
    );

    final bubbleLeft = (hx + 8).clamp(chart.left, chart.right - 82);
    final bubbleTop = (hy - 44).clamp(chart.top, chart.bottom - 48);
    final bubble = RRect.fromRectAndRadius(
      Rect.fromLTWH(bubbleLeft.toDouble(), bubbleTop.toDouble(), 82, 48),
      const Radius.circular(8),
    );
    canvas.drawRRect(
      bubble,
      Paint()
        ..color = Colors.white
        ..style = PaintingStyle.fill,
    );
    canvas.drawRRect(
      bubble,
      Paint()
        ..color = AppColors.grid
        ..style = PaintingStyle.stroke,
    );
    final highlightTime = series.timestamps.isEmpty
        ? null
        : series.timestamps[highlightIndex];
    textPainter.text = TextSpan(
      text:
          '${_formatBubbleTime(highlightTime)}\n${series.valueLabel(series.values[highlightIndex])}',
      style: const TextStyle(color: AppColors.ink, fontSize: 11, height: 1.4),
    );
    textPainter.textAlign = TextAlign.left;
    textPainter.layout(maxWidth: 74);
    textPainter.paint(canvas, Offset(bubbleLeft + 8, bubbleTop + 7));
  }

  double _mapX(int index, Rect chart) {
    if (series.values.length == 1) return chart.center.dx;
    return chart.left + chart.width * index / (series.values.length - 1);
  }

  double _mapY(double value, Rect chart) {
    final normalized = (value - series.low) / (series.high - series.low);
    return chart.bottom - chart.height * normalized.clamp(0, 1);
  }

  DateTime? _timestampAt(double fraction) {
    if (series.timestamps.isEmpty) return null;
    if (series.timestamps.length == 1) return series.timestamps.first;
    final index = (fraction * (series.timestamps.length - 1)).round();
    return series.timestamps[index.clamp(0, series.timestamps.length - 1)];
  }

  String _formatAxisDate(DateTime? value) {
    if (value == null) return '';
    return '${_two(value.month)}/${_two(value.day)}';
  }

  String _formatBubbleTime(DateTime? value) {
    if (value == null) return '--';
    return '${_two(value.month)}/${_two(value.day)} ${_two(value.hour)}:${_two(value.minute)}';
  }

  String _two(int value) => value.toString().padLeft(2, '0');

  void _drawDashed(Canvas canvas, Offset start, Offset end, Paint paint) {
    const dash = 5.0;
    const gap = 5.0;
    final total = (end - start).distance;
    final direction = (end - start) / total;
    var distance = 0.0;
    while (distance < total) {
      final from = start + direction * distance;
      final to = start + direction * math.min(distance + dash, total);
      canvas.drawLine(from, to, paint);
      distance += dash + gap;
    }
  }

  @override
  bool shouldRepaint(covariant HistoryChartPainter oldDelegate) {
    return oldDelegate.series != series;
  }
}

class SparklinePainter extends CustomPainter {
  SparklinePainter({required this.color});

  final Color color;

  @override
  void paint(Canvas canvas, Size size) {
    final values = [0.2, 0.35, 0.72, 0.48, 0.4, 0.55, 0.62];
    final path = Path();
    for (var i = 0; i < values.length; i++) {
      final x = size.width * i / (values.length - 1);
      final y = size.height - values[i] * size.height;
      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }
    canvas.drawPath(
      path,
      Paint()
        ..color = color
        ..strokeWidth = 2
        ..style = PaintingStyle.stroke
        ..strokeCap = StrokeCap.round,
    );
  }

  @override
  bool shouldRepaint(covariant SparklinePainter oldDelegate) {
    return oldDelegate.color != color;
  }
}
