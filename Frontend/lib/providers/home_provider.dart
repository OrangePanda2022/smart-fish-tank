import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/domain/models/QuickAction.dart';
import 'package:aqua/domain/models/TankState.dart';

final tankStatsProvider = NotifierProvider<TankStatsNotifier, List<TankStat>>(
  TankStatsNotifier.new,
);

class TankStatsNotifier extends Notifier<List<TankStat>> {
  @override
  List<TankStat> build() {
    return [
      TankStat(
        title: '水温',
        value: '26.5',
        unit: '°C',
        status: '当前温度适宜',
        badgeText: 'OFF',
        iconName: Icons.thermostat,
      ),
      TankStat(
        title: '酸碱',
        value: '7.2',
        status: '当前酸碱适中',
        badgeText: 'ON',
        iconName: Icons.science,
      ),
      TankStat(
        title: '水位',
        value: '95%',
        status: '当前水位正常',
        badgeText: 'OFF',
        iconName: Icons.water_drop,
      ),
    ];
  }
}

final quickActionsProvider =
    NotifierProvider<QuickActionsNotifier, List<QuickAction>>(
      QuickActionsNotifier.new,
    );

class QuickActionsNotifier extends Notifier<List<QuickAction>> {
  @override
  List<QuickAction> build() {
    return [
      QuickAction(label: '立即喂食', iconName: Icons.restaurant),
      QuickAction(label: '照明灯', iconName: Icons.lightbulb, active: true),
      QuickAction(label: '造浪泵', iconName: Icons.water),
    ];
  }
}
