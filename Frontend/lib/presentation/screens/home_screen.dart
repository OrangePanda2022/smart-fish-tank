import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../domain/entities/analysis_report.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/repositories/aquarium_repository.dart';
import '../widgets/ai_analysis_card.dart';
import '../widgets/app_background.dart';
import '../widgets/aquarium_status_card.dart';
import '../widgets/device_grid.dart';
import '../widgets/headers.dart';
import '../widgets/section_header.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({
    required this.dashboard,
    required this.repository,
    required this.onRefresh,
    required this.onAnalyze,
    super.key,
  });

  final AquariumDashboard dashboard;
  final AquariumRepository repository;
  final VoidCallback onRefresh;
  final Future<AnalysisReport> Function(String tankId) onAnalyze;

  @override
  Widget build(BuildContext context) {
    return AppBackground(
      child: SafeArea(
        bottom: false,
        child: RefreshIndicator(
          onRefresh: () async => onRefresh(),
          color: AppColors.primary,
          child: ListView(
            padding: const EdgeInsets.fromLTRB(20, 16, 20, 112),
            children: [
              HomeHeader(tank: dashboard.tank),
              const SizedBox(height: 18),
              AquariumStatusCard(dashboard: dashboard),
              const SizedBox(height: 22),
              AiAnalysisCard(
                dashboard: dashboard,
                repository: repository,
                onAnalyze: onAnalyze,
                onAnalysisComplete: onRefresh,
              ),
              const SizedBox(height: 22),
              SectionHeader(
                title: context.tr('设备控制'),
                actionText: context.tr('全部设备'),
                onAction: () {},
              ),
              const SizedBox(height: 14),
              DeviceGrid(sensor: dashboard.latest),
            ],
          ),
        ),
      ),
    );
  }
}
