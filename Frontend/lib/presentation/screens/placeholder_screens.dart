import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../widgets/app_background.dart';
import '../widgets/device_grid.dart';
import '../widgets/soft_card.dart';

class DevicesPlaceholder extends StatelessWidget {
  const DevicesPlaceholder({required this.dashboard, super.key});

  final AquariumDashboard dashboard;

  @override
  Widget build(BuildContext context) {
    return AppBackground(
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(22),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('设备', style: AppText.title),
              const SizedBox(height: 18),
              DeviceGrid(sensor: dashboard.latest),
            ],
          ),
        ),
      ),
    );
  }
}

class ProfilePlaceholder extends StatelessWidget {
  const ProfilePlaceholder({required this.dashboard, super.key});

  final AquariumDashboard dashboard;

  @override
  Widget build(BuildContext context) {
    return AppBackground(
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('我的', style: AppText.title),
              const SizedBox(height: 24),
              SoftCard(
                padding: const EdgeInsets.all(22),
                child: Row(
                  children: [
                    const CircleAvatar(
                      radius: 30,
                      backgroundColor: AppColors.paleBlue,
                      child: Icon(
                        CupertinoIcons.person_fill,
                        color: AppColors.primary,
                        size: 30,
                      ),
                    ),
                    const SizedBox(width: 16),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(dashboard.tank.name, style: AppText.cardTitle),
                        const SizedBox(height: 6),
                        Text(
                          '${dashboard.tank.size}L · ${dashboard.tank.fishCount} 条鱼',
                          style: AppText.caption,
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
