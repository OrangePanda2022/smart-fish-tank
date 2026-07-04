import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/tank_info.dart';
import '../view_models/device_control.dart';
import 'soft_card.dart';

class DeviceSummaryCard extends StatelessWidget {
  const DeviceSummaryCard({
    required this.tank,
    required this.devices,
    super.key,
  });

  final TankInfo tank;
  final List<DeviceControl> devices;

  @override
  Widget build(BuildContext context) {
    final onlineCount = devices.where((device) => device.isOnline).length;
    final enabledCount = devices.where((device) => device.isEnabled).length;
    final offlineCount = devices.length - onlineCount;

    return SoftCard(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  color: AppColors.primary.withAlpha(24),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: const Icon(
                  CupertinoIcons.cube_box_fill,
                  color: AppColors.primary,
                  size: 28,
                ),
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(tank.name, style: AppText.cardTitle),
                    const SizedBox(height: 4),
                    Text(
                      context.l10n.tankManaged(tank.size, tank.fishCount),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: AppText.caption,
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 18),
          Row(
            children: [
              Expanded(
                child: _SummaryStat(
                  value: '$onlineCount',
                  label: context.tr('在线设备'),
                  color: AppColors.success,
                ),
              ),
              Expanded(
                child: _SummaryStat(
                  value: '$enabledCount',
                  label: context.tr('运行中'),
                  color: AppColors.primary,
                ),
              ),
              Expanded(
                child: _SummaryStat(
                  value: '$offlineCount',
                  label: context.tr('需处理'),
                  color: offlineCount > 0 ? AppColors.orange : AppColors.muted,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class DeviceQuickControls extends StatelessWidget {
  const DeviceQuickControls({
    required this.devices,
    required this.enabledStates,
    required this.onToggle,
    super.key,
  });

  final List<DeviceControl> devices;
  final Map<String, bool> enabledStates;
  final void Function(DeviceControl device, bool value) onToggle;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 132,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        physics: const BouncingScrollPhysics(),
        itemCount: devices.length,
        separatorBuilder: (_, _) => const SizedBox(width: 12),
        itemBuilder: (context, index) {
          final device = devices[index];
          final enabled = enabledStates[device.id] ?? device.isEnabled;

          return SizedBox(
            width: 126,
            child: SoftCard(
              padding: const EdgeInsets.all(14),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(device.icon, color: device.color, size: 28),
                      const Spacer(),
                      Switch.adaptive(
                        value: enabled && device.isOnline,
                        onChanged: device.isOnline
                            ? (value) => onToggle(device, value)
                            : null,
                        materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                      ),
                    ],
                  ),
                  const Spacer(),
                  Text(
                    device.title,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: AppText.cardTitle,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    device.status,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: AppText.caption,
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class DeviceListTile extends StatelessWidget {
  const DeviceListTile({
    required this.device,
    required this.enabled,
    required this.onTap,
    required this.onToggle,
    super.key,
  });

  final DeviceControl device;
  final bool enabled;
  final VoidCallback onTap;
  final ValueChanged<bool> onToggle;

  @override
  Widget build(BuildContext context) {
    final active = enabled && device.isOnline;

    return GestureDetector(
      onTap: onTap,
      child: SoftCard(
        padding: const EdgeInsets.fromLTRB(16, 14, 12, 14),
        child: Row(
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: device.color.withAlpha(device.isOnline ? 30 : 14),
                borderRadius: BorderRadius.circular(16),
              ),
              child: Icon(device.icon, color: device.color, size: 27),
            ),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Flexible(
                        child: Text(
                          device.title,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                          style: AppText.cardTitle,
                        ),
                      ),
                      const SizedBox(width: 8),
                      _StatusPill(
                        text: device.isOnline
                            ? (active ? context.tr('运行中') : context.tr('待机'))
                            : context.tr('离线'),
                        color: device.isOnline
                            ? (active ? AppColors.success : AppColors.primary)
                            : AppColors.orange,
                      ),
                    ],
                  ),
                  const SizedBox(height: 5),
                  Text(
                    device.maintenance,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: AppText.caption,
                  ),
                ],
              ),
            ),
            const SizedBox(width: 8),
            Switch.adaptive(
              value: active,
              onChanged: device.isOnline ? onToggle : null,
              materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
            ),
            const Icon(
              CupertinoIcons.chevron_right,
              color: AppColors.muted,
              size: 18,
            ),
          ],
        ),
      ),
    );
  }
}

class DeviceDetailHeader extends StatelessWidget {
  const DeviceDetailHeader({
    required this.device,
    required this.enabled,
    required this.onToggle,
    super.key,
  });

  final DeviceControl device;
  final bool enabled;
  final ValueChanged<bool> onToggle;

  @override
  Widget build(BuildContext context) {
    final active = enabled && device.isOnline;

    return SoftCard(
      padding: const EdgeInsets.all(22),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                width: 68,
                height: 68,
                decoration: BoxDecoration(
                  color: device.color.withAlpha(32),
                  borderRadius: BorderRadius.circular(22),
                ),
                child: Icon(device.icon, color: device.color, size: 38),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(device.title, style: AppText.titleSmall),
                    const SizedBox(height: 8),
                    _StatusPill(
                      text: device.isOnline
                          ? (active ? context.tr('运行中') : context.tr('待机'))
                          : context.tr('离线'),
                      color: device.isOnline
                          ? (active ? AppColors.success : AppColors.primary)
                          : AppColors.orange,
                    ),
                  ],
                ),
              ),
              Switch.adaptive(
                value: active,
                onChanged: device.isOnline ? onToggle : null,
              ),
            ],
          ),
          const SizedBox(height: 18),
          Text(device.description, style: AppText.body),
        ],
      ),
    );
  }
}

class DeviceControlPanel extends StatelessWidget {
  const DeviceControlPanel({required this.device, super.key});

  final DeviceControl device;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      padding: const EdgeInsets.all(18),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(context.tr('控制与状态'), style: AppText.cardTitle),
          const SizedBox(height: 16),
          Row(
            children: device.metrics
                .map(
                  (metric) => Expanded(
                    child: _MetricBlock(
                      label: metric.label,
                      value: metric.value,
                    ),
                  ),
                )
                .toList(),
          ),
          const SizedBox(height: 18),
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
            decoration: BoxDecoration(
              color: AppColors.paleBlue,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Row(
              children: [
                Icon(device.icon, color: device.color, size: 22),
                const SizedBox(width: 10),
                Expanded(
                  child: Text(device.primaryAction, style: AppText.blueHeading),
                ),
                const Icon(
                  CupertinoIcons.chevron_right,
                  color: AppColors.primary,
                  size: 18,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _SummaryStat extends StatelessWidget {
  const _SummaryStat({
    required this.value,
    required this.label,
    required this.color,
  });

  final String value;
  final String label;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(value, style: AppText.statValue.copyWith(color: color)),
        const SizedBox(height: 4),
        Text(label, style: AppText.caption),
      ],
    );
  }
}

class _StatusPill extends StatelessWidget {
  const _StatusPill({required this.text, required this.color});

  final String text;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withAlpha(28),
        borderRadius: BorderRadius.circular(999),
      ),
      child: Text(
        text,
        maxLines: 1,
        style: TextStyle(
          color: color,
          fontSize: 12,
          fontWeight: FontWeight.w800,
        ),
      ),
    );
  }
}

class _MetricBlock extends StatelessWidget {
  const _MetricBlock({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          value,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          style: AppText.statValue,
        ),
        const SizedBox(height: 4),
        Text(
          label,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          style: AppText.caption,
        ),
      ],
    );
  }
}
