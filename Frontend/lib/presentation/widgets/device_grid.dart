import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/sensor_reading.dart';
import '../view_models/device_catalog.dart';
import '../view_models/device_control.dart';
import 'soft_card.dart';

class DeviceGrid extends StatelessWidget {
  const DeviceGrid({required this.sensor, super.key});

  final SensorReading sensor;

  @override
  Widget build(BuildContext context) {
    final devices = [
      ...buildDeviceCatalog(context, sensor).take(7),
      DeviceControl(
        context.tr('更多功能'),
        '',
        CupertinoIcons.square_grid_2x2_fill,
        AppColors.primary,
      ),
      DeviceControl(
        context.tr('添加设备'),
        '',
        CupertinoIcons.plus_circle,
        AppColors.primary,
      ),
    ];

    return LayoutBuilder(
      builder: (context, constraints) {
        final isNarrow = constraints.maxWidth < 340;

        return GridView.builder(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          itemCount: devices.length,
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: isNarrow ? 2 : 3,
            mainAxisSpacing: 12,
            crossAxisSpacing: 12,
            childAspectRatio: isNarrow ? 1.45 : 1.28,
          ),
          itemBuilder: (context, index) =>
              DeviceTile(device: devices[index], compact: isNarrow),
        );
      },
    );
  }
}

class DeviceTile extends StatelessWidget {
  const DeviceTile({required this.device, this.compact = false, super.key});

  final DeviceControl device;
  final bool compact;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      padding: EdgeInsets.symmetric(
        horizontal: compact ? 12 : 10,
        vertical: compact ? 12 : 14,
      ),
      child: Row(
        children: [
          Expanded(
            flex: 5,
            child: FittedBox(
              fit: BoxFit.scaleDown,
              child: Icon(device.icon, color: device.color, size: 42),
            ),
          ),
          const SizedBox(width: 8),
          Expanded(
            flex: 6,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                FittedBox(
                  fit: BoxFit.scaleDown,
                  alignment: Alignment.centerLeft,
                  child: Text(device.title, style: AppText.deviceTitle),
                ),
                if (device.status.isNotEmpty) ...[
                  const SizedBox(height: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 9,
                      vertical: 3,
                    ),
                    decoration: BoxDecoration(
                      color: device.status == 'ON'
                          ? AppColors.primary.withAlpha(38)
                          : Colors.transparent,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      device.status,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: device.color == AppColors.orange
                            ? AppColors.orange
                            : device.id == 'camera' || device.id == 'feeder'
                            ? AppColors.success
                            : AppColors.primary,
                        fontSize: 14,
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }
}
