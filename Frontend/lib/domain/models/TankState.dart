import 'package:flutter/material.dart';

class TankStat {
  final String title;
  final String value;
  final String unit;
  final String status;
  final String badgeText;
  final IconData iconName;

  TankStat({
    required this.title,
    required this.value,
    this.unit = '',
    this.status = '',
    this.badgeText = '',
    required this.iconName,
  });
}