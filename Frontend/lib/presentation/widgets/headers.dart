import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/tank_info.dart';

class HomeHeader extends StatelessWidget {
  const HomeHeader({required this.tank, super.key});

  final TankInfo tank;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Expanded(
              child: Row(
                children: [
                  Flexible(
                    child: Text(context.tr('我的鱼缸'), style: AppText.title),
                  ),
                  const SizedBox(width: 7),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 22),
        Row(
          children: [
            Flexible(
              child: Text(
                tank.name,
                style: const TextStyle(
                  color: AppColors.secondaryText,
                  fontWeight: FontWeight.w600,
                  fontSize: 20,
                ),
              ),
            ),
            const SizedBox(width: 8),
            const Icon(CupertinoIcons.wifi, color: AppColors.primary, size: 22),
            const Spacer(),
          ],
        ),
      ],
    );
  }
}

class DataHeader extends StatelessWidget {
  const DataHeader({super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(context.tr('数据中心'), style: AppText.compactTitle),
        const SizedBox(width: 8),
        const Icon(
          CupertinoIcons.question_circle,
          color: AppColors.secondaryText,
        ),
      ],
    );
  }
}

class HistoryTopBar extends StatelessWidget {
  const HistoryTopBar({required this.onRefresh, super.key});

  final VoidCallback onRefresh;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 48,
      child: Stack(
        alignment: Alignment.center,
        children: [
          Align(
            alignment: Alignment.centerLeft,
            child: IconButton(
              onPressed: () => Navigator.of(context).pop(),
              icon: const Icon(CupertinoIcons.chevron_left, size: 28),
            ),
          ),
          Text(context.tr('历史数据'), style: AppText.navTitle),
          Align(
            alignment: Alignment.centerRight,
            child: IconButton(
              onPressed: onRefresh,
              icon: const Icon(CupertinoIcons.refresh, size: 27),
            ),
          ),
        ],
      ),
    );
  }
}

class NotificationBell extends StatelessWidget {
  const NotificationBell({super.key});

  @override
  Widget build(BuildContext context) {
    return Stack(
      clipBehavior: Clip.none,
      children: [
        const Icon(CupertinoIcons.bell, color: AppColors.ink, size: 30),
        Positioned(
          top: -8,
          right: -7,
          child: Container(
            width: 23,
            height: 23,
            alignment: Alignment.center,
            decoration: const BoxDecoration(
              color: AppColors.alert,
              shape: BoxShape.circle,
            ),
            child: const Text(
              '3',
              style: TextStyle(
                color: Colors.white,
                fontSize: 12,
                fontWeight: FontWeight.w800,
              ),
            ),
          ),
        ),
      ],
    );
  }
}
