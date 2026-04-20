import 'package:flutter/material.dart';

enum SettingItemType { normal, switchTile }

class SettingItem {
  final String title;
  final String? subtitle;
  final IconData icon;
  final SettingItemType type;
  dynamic value;
  final Function(dynamic)? onChanged;

  SettingItem({
    required this.title,
    this.subtitle,
    required this.icon,
    required this.type,
    this.value,
    this.onChanged,
  });
}

class SettingGroup {
  final String title;
  final List<SettingItem> items;

  SettingGroup({required this.title, required this.items});
}