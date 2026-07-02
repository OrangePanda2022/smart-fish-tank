import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../domain/entities/sensor_reading.dart';

class HistoryMetric {
  const HistoryMetric({
    required this.key,
    required this.label,
    required this.title,
    required this.unit,
    required this.precision,
    required this.axisPrecision,
    required this.safeLow,
    required this.safeHigh,
    required this.color,
    required this.valueOf,
  });

  final String key;
  final String label;
  final String title;
  final String unit;
  final int precision;
  final int axisPrecision;
  final double safeLow;
  final double safeHigh;
  final Color color;
  final double Function(SensorReading reading) valueOf;
}

class ChartSeries {
  ChartSeries({
    required this.metric,
    required this.title,
    required this.values,
    required this.timestamps,
    required this.color,
    required this.unit,
    required this.precision,
    required this.axisPrecision,
    required this.safeLow,
    required this.safeHigh,
  }) : min = values.isEmpty ? 0 : values.reduce(math.min),
       max = values.isEmpty ? 0 : values.reduce(math.max),
       average = values.isEmpty
           ? 0
           : values.reduce((a, b) => a + b) / values.length {
    final range = max - min;
    final padding = math.max(range * 0.5, _minimumPadding(metric));
    low = math.min(safeLow, min - padding);
    high = math.max(safeHigh, max + padding);
  }

  final String metric;
  final String title;
  final List<double> values;
  final List<DateTime> timestamps;
  final Color color;
  final String unit;
  final int precision;
  final int axisPrecision;
  final double safeLow;
  final double safeHigh;
  final double min;
  final double max;
  final double average;
  late final double low;
  late final double high;

  bool get isEmpty => values.isEmpty;

  String valueLabel(double value) {
    return '${value.toStringAsFixed(precision)}$unit';
  }

  static const metrics = <HistoryMetric>[
    HistoryMetric(
      key: 'ph',
      label: 'pH值',
      title: 'pH值',
      unit: '',
      precision: 2,
      axisPrecision: 1,
      safeLow: 6.5,
      safeHigh: 7.8,
      color: AppColors.primary,
      valueOf: _phOf,
    ),
    HistoryMetric(
      key: 'temperature',
      label: '温度',
      title: '温度 (°C)',
      unit: '°C',
      precision: 1,
      axisPrecision: 0,
      safeLow: 23.5,
      safeHigh: 29,
      color: AppColors.orange,
      valueOf: _temperatureOf,
    ),
    HistoryMetric(
      key: 'oxygen',
      label: '溶解氧',
      title: '溶解氧 (mg/L)',
      unit: ' mg/L',
      precision: 1,
      axisPrecision: 0,
      safeLow: 6,
      safeHigh: 10,
      color: AppColors.teal,
      valueOf: _oxygenOf,
    ),
    HistoryMetric(
      key: 'tds',
      label: 'TDS',
      title: 'TDS (ppm)',
      unit: ' ppm',
      precision: 0,
      axisPrecision: 0,
      safeLow: 120,
      safeHigh: 350,
      color: AppColors.purple,
      valueOf: _tdsOf,
    ),
    HistoryMetric(
      key: 'ammonia',
      label: '氨氮',
      title: '氨氮 (mg/L)',
      unit: ' mg/L',
      precision: 2,
      axisPrecision: 2,
      safeLow: 0,
      safeHigh: 0.08,
      color: AppColors.alert,
      valueOf: _ammoniaOf,
    ),
    HistoryMetric(
      key: 'nitrate',
      label: '硝酸盐',
      title: '硝酸盐 (mg/L)',
      unit: ' mg/L',
      precision: 1,
      axisPrecision: 0,
      safeLow: 0,
      safeHigh: 20,
      color: AppColors.green,
      valueOf: _nitrateOf,
    ),
    HistoryMetric(
      key: 'nitrite',
      label: '亚硝酸盐',
      title: '亚硝酸盐 (mg/L)',
      unit: ' mg/L',
      precision: 2,
      axisPrecision: 2,
      safeLow: 0,
      safeHigh: 0.1,
      color: AppColors.yellow,
      valueOf: _nitriteOf,
    ),
    HistoryMetric(
      key: 'chloride',
      label: '氯离子',
      title: '氯离子 (mg/L)',
      unit: ' mg/L',
      precision: 1,
      axisPrecision: 0,
      safeLow: 0,
      safeHigh: 10,
      color: AppColors.sky,
      valueOf: _chlorideOf,
    ),
    HistoryMetric(
      key: 'waterLevel',
      label: '水位',
      title: '水位 (%)',
      unit: '%',
      precision: 0,
      axisPrecision: 0,
      safeLow: 60,
      safeHigh: 100,
      color: AppColors.periwinkle,
      valueOf: _waterLevelOf,
    ),
  ];

  static HistoryMetric metricFor(String key) {
    return metrics.firstWhere(
      (metric) => metric.key == key,
      orElse: () => metrics.first,
    );
  }

  factory ChartSeries.from(List<SensorReading> history, String metric) {
    final option = metricFor(metric);
    final sorted = [...history]
      ..sort((a, b) => a.timestamp.compareTo(b.timestamp));
    final points = sorted
        .map((reading) => (reading.timestamp, option.valueOf(reading)))
        .where((point) => point.$2 > 0)
        .toList();

    return ChartSeries(
      metric: option.key,
      title: option.title,
      values: points.map((point) => point.$2).toList(),
      timestamps: points.map((point) => point.$1).toList(),
      color: option.color,
      unit: option.unit,
      precision: option.precision,
      axisPrecision: option.axisPrecision,
      safeLow: option.safeLow,
      safeHigh: option.safeHigh,
    );
  }
}

double _minimumPadding(String metric) {
  return switch (metric) {
    'ph' => 0.35,
    'temperature' => 1.4,
    'ammonia' || 'nitrite' => 0.03,
    'tds' => 20,
    'waterLevel' => 8,
    _ => 1,
  };
}

double _temperatureOf(SensorReading reading) => reading.temperature;
double _phOf(SensorReading reading) => reading.ph;
double _oxygenOf(SensorReading reading) => reading.oxygen;
double _ammoniaOf(SensorReading reading) => reading.ammonia;
double _waterLevelOf(SensorReading reading) => reading.waterLevel;
double _tdsOf(SensorReading reading) => reading.tds;
double _nitrateOf(SensorReading reading) => reading.nitrate;
double _nitriteOf(SensorReading reading) => reading.nitrite;
double _chlorideOf(SensorReading reading) => reading.chloride;
