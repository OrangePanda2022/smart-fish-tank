import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/repositories/aquarium_repository.dart';
import '../widgets/app_background.dart';
import '../widgets/headers.dart';
import '../widgets/prediction_card.dart';
import '../widgets/section_header.dart';
import '../widgets/sensor_metric_grid.dart';
import '../widgets/wide_action_button.dart';
import 'history_screen.dart';

class DataCenterScreen extends StatelessWidget {
  const DataCenterScreen({
    required this.dashboard,
    required this.repository,
    super.key,
  });

  final AquariumDashboard dashboard;
  final AquariumRepository repository;

  @override
  Widget build(BuildContext context) {
    return AppBackground(
      child: SafeArea(
        bottom: false,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(12, 28, 12, 112),
          children: [
            const DataHeader(),
            const SizedBox(height: 22),
            PredictionCard(dashboard: dashboard),
            const SizedBox(height: 26),
            SectionHeader(
              title: '实时传感器数据',
              actionText: '更新时间：09:41',
              showChevron: false,
              compact: true,
              onAction: () {},
            ),
            const SizedBox(height: 14),
            SensorMetricGrid(sensor: dashboard.latest),
            const SizedBox(height: 18),
            WideActionButton(
              text: '全部数据',
              icon: CupertinoIcons.chevron_right,
              onTap: () {
                Navigator.of(context).push(
                  MaterialPageRoute<void>(
                    builder: (_) => HistoryScreen(
                      dashboard: dashboard,
                      repository: repository,
                    ),
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }
}
