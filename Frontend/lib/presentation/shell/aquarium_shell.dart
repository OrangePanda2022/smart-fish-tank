import 'package:flutter/material.dart';

import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/repositories/aquarium_repository.dart';
import '../../infrastructure/repositories/aqua_api_repository.dart';
import '../screens/data_center_screen.dart';
import '../screens/device_screen.dart';
import '../screens/home_screen.dart';
import '../screens/profile_screen.dart';
import '../widgets/aqua_bottom_nav.dart';

class AquariumShell extends StatefulWidget {
  const AquariumShell({super.key});

  @override
  State<AquariumShell> createState() => _AquariumShellState();
}

class _AquariumShellState extends State<AquariumShell> {
  final AquariumRepository _repository = AquaApiRepository();
  late Future<AquariumDashboard> _dashboardFuture;
  AquariumDashboard? _cachedDashboard;
  int _selectedIndex = 0;

  @override
  void initState() {
    super.initState();
    _dashboardFuture = _loadDashboard();
  }

  Future<AquariumDashboard> _loadDashboard() async {
    final dashboard = await _repository.loadDashboard();
    final resolved = _withCachedScore(dashboard);
    _cachedDashboard = resolved;
    return resolved;
  }

  AquariumDashboard _withCachedScore(AquariumDashboard dashboard) {
    final cached = _cachedDashboard;
    if (dashboard.score != null ||
        cached?.score == null ||
        cached?.tank.id != dashboard.tank.id) {
      return dashboard;
    }

    return AquariumDashboard(
      tank: dashboard.tank,
      latest: dashboard.latest,
      prediction: dashboard.prediction,
      score: cached!.score,
      analysisSummary: dashboard.analysisSummary,
    );
  }

  void _refresh() {
    setState(() {
      _dashboardFuture = _loadDashboard();
    });
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<AquariumDashboard>(
      future: _dashboardFuture,
      builder: (context, snapshot) {
        final dashboard =
            snapshot.data ?? _cachedDashboard ?? AquariumDashboard.empty();
        final pages = [
          HomeScreen(
            dashboard: dashboard,
            repository: _repository,
            onRefresh: _refresh,
            onAnalyze: _repository.runAnalysis,
          ),
          DataCenterScreen(dashboard: dashboard, repository: _repository),
          DeviceScreen(dashboard: dashboard),
          ProfileScreen(dashboard: dashboard),
        ];

        return Scaffold(
          extendBody: true,
          body: AnimatedSwitcher(
            duration: const Duration(milliseconds: 220),
            child: pages[_selectedIndex],
          ),
          bottomNavigationBar: SafeArea(
            minimum: const EdgeInsets.fromLTRB(18, 0, 18, 10),
            child: AquaBottomNav(
              selectedIndex: _selectedIndex,
              onChanged: (index) {
                setState(() {
                  _selectedIndex = index;
                });
              },
            ),
          ),
        );
      },
    );
  }
}
