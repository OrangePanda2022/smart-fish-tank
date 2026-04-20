import 'package:flutter/material.dart';

enum StatType { temperature, ph, waterLevel }

extension StatTypeExtension on StatType {
  String get title {
    switch (this) {
      case StatType.temperature:
        return '水温';
      case StatType.ph:
        return '酸碱度';
      case StatType.waterLevel:
        return '水位';
    }
  }

  String get unit {
    switch (this) {
      case StatType.temperature:
        return '°C';
      case StatType.ph:
        return '';
      case StatType.waterLevel:
        return '%';
    }
  }

  IconData get icon {
    switch (this) {
      case StatType.temperature:
        return Icons.thermostat;
      case StatType.ph:
        return Icons.science;
      case StatType.waterLevel:
        return Icons.water_drop;
    }
  }

  double get normalMin {
    switch (this) {
      case StatType.temperature:
        return 24.0;
      case StatType.ph:
        return 6.5;
      case StatType.waterLevel:
        return 80.0;
    }
  }

  double get normalMax {
    switch (this) {
      case StatType.temperature:
        return 28.0;
      case StatType.ph:
        return 7.8;
      case StatType.waterLevel:
        return 100.0;
    }
  }
}

class HistoryDataPoint {
  final DateTime timestamp;
  final double value;

  HistoryDataPoint({required this.timestamp, required this.value});
}

class TankHistoryData {
  final StatType type;
  final String unit;
  final List<HistoryDataPoint> dataPoints;

  TankHistoryData({
    required this.type,
    required this.unit,
    required this.dataPoints,
  });

  double get maxValue =>
      dataPoints.map((e) => e.value).reduce((a, b) => a > b ? a : b);
  double get minValue =>
      dataPoints.map((e) => e.value).reduce((a, b) => a < b ? a : b);
  double get avgValue =>
      dataPoints.map((e) => e.value).reduce((a, b) => a + b) /
      dataPoints.length;
}
