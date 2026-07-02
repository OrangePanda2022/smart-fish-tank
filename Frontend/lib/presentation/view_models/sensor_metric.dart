import 'package:flutter/material.dart';

class SensorMetric {
  const SensorMetric(this.label, this.value, this.unit, this.icon, this.color);

  final String label;
  final String value;
  final String unit;
  final IconData icon;
  final Color color;
}
