import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/presentation/widgets/ImageCard.dart';
import 'package:aqua/presentation/widgets/StateBadge.dart';
import 'package:aqua/presentation/widgets/StateCard.dart';
import 'package:aqua/presentation/widgets/ActionButton.dart';
import 'package:aqua/presentation/widgets/AnalysisButton.dart';
import 'package:aqua/providers/home_provider.dart';
import 'package:aqua/domain/models/tank_history.dart';
import 'package:aqua/presentation/pages/history.dart';

class HomePage extends ConsumerWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final tankStats = ref.watch(tankStatsProvider);
    final quickActions = ref.watch(quickActionsProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('智慧鱼缸')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const ImageCard(
              imageUrl: 'lib/assets/img/aqua.png',
              width: double.infinity,
              height: 380,
            ),
            const SizedBox(height: 16),
            _buildTankInfoCard(context),
            const SizedBox(height: 16),
            _buildTankCardLayout(context, tankStats),
            const SizedBox(height: 2),
            _buildQuickActions(context, ref, quickActions),
            const SizedBox(height: 24),
            const AnalysisButton(),
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Widget _buildTankInfoCard(BuildContext context) {
    return Card(
      color: Theme.of(context).colorScheme.surface,
      borderOnForeground: false,
      elevation: 0,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(32)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            const Icon(Icons.water),
            const SizedBox(width: 12),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: const [
                Text(
                  '主鱼缸',
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18),
                ),
                Text('容量: 120L'),
              ],
            ),
            const Spacer(),
            StatusBadge(
              text: '系统运行正常',
              backgroundColor: Colors.green.shade100,
              textColor: Colors.green.shade900,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildQuickActions(
    BuildContext context,
    WidgetRef ref,
    List quickActions,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '快捷操作',
          style: TextStyle(
            fontWeight: FontWeight.bold,
            fontSize: 18,
            color: Colors.black87,
          ),
        ),
        const SizedBox(height: 12),
        Row(
          children: quickActions.map((action) {
            return ActionButton(
              icon: action.iconName,
              label: action.label,
              active: action.active,
              onTap: () {},
            );
          }).toList(),
        ),
      ],
    );
  }

  Widget _buildTankCardLayout(BuildContext context, List tankStats) {
    if (tankStats.isEmpty) return const SizedBox();

    final statTypes = [StatType.temperature, StatType.ph, StatType.waterLevel];

    return Column(
      children: [
        StatCard(
          title: tankStats[0].title,
          value: tankStats[0].value,
          unit: tankStats[0].unit,
          icon: tankStats[0].iconName,
          backgroundColor: darken(
            Theme.of(context).colorScheme.primaryContainer,
            .05,
          ),
          status: tankStats[0].status,
          badgeText: tankStats[0].badgeText,
          statType: statTypes[0],
          onTap: () => _navigateToHistory(context, statTypes[0]),
        ),
        const SizedBox(height: 16),
        ..._buildRemainingCards(context, tankStats, statTypes),
      ],
    );
  }

  List<Widget> _buildRemainingCards(
    BuildContext context,
    List tankStats,
    List<StatType> statTypes,
  ) {
    List<Widget> widgets = [];

    for (int i = 1; i < tankStats.length; i += 2) {
      final left = tankStats[i];
      final right = (i + 1 < tankStats.length) ? tankStats[i + 1] : null;

      widgets.add(
        Row(
          children: [
            Expanded(
              child: StatCard(
                title: left.title,
                value: left.value,
                unit: left.unit,
                icon: left.iconName,
                backgroundColor: Theme.of(
                  context,
                ).colorScheme.secondaryContainer,
                status: left.status,
                badgeText: left.badgeText,
                statType: statTypes[i],
                onTap: () => _navigateToHistory(context, statTypes[i]),
              ),
            ),
            const SizedBox(width: 12),
            if (right != null)
              Expanded(
                child: StatCard(
                  title: right.title,
                  value: right.value,
                  unit: right.unit,
                  icon: right.iconName,
                  backgroundColor: darken(
                    Theme.of(context).colorScheme.secondaryContainer,
                    .025,
                  ),
                  status: right.status,
                  badgeText: right.badgeText,
                  statType: statTypes[i + 1],
                  onTap: () => _navigateToHistory(context, statTypes[i + 1]),
                ),
              )
            else
              const Expanded(child: SizedBox()),
          ],
        ),
      );

      widgets.add(const SizedBox(height: 12));
    }

    return widgets;
  }

  void _navigateToHistory(BuildContext context, StatType type) {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => TankHistoryPage(statType: type)),
    );
  }

  Color darken(Color color, [double amount = .1]) {
    final hsl = HSLColor.fromColor(color);
    final hslDark = hsl.withLightness((hsl.lightness - amount).clamp(0.0, 1.0));
    return hslDark.toColor();
  }
}
