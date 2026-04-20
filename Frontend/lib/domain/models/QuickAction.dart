import 'package:flutter/material.dart';

class QuickAction {
  final String label;
  final IconData iconName;
  final bool active;

  QuickAction({
    required this.label,
    required this.iconName,
    this.active = false,
  });
}