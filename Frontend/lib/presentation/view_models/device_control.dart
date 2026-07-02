import 'package:flutter/material.dart';

class DeviceControl {
  const DeviceControl(
    this.title,
    this.status,
    this.icon,
    this.color, {
    this.id = '',
    this.description = '',
    this.primaryAction = '',
    this.metrics = const [],
    this.isOnline = true,
    this.isEnabled = true,
    this.maintenance = '',
  });

  final String title;
  final String status;
  final IconData icon;
  final Color color;
  final String id;
  final String description;
  final String primaryAction;
  final List<DeviceMetric> metrics;
  final bool isOnline;
  final bool isEnabled;
  final String maintenance;
}

class DeviceMetric {
  const DeviceMetric({required this.label, required this.value});

  final String label;
  final String value;
}
