import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/repositories/aquarium_repository.dart';
import '../view_models/device_catalog.dart';
import '../view_models/device_control.dart';
import '../widgets/app_background.dart';
import '../widgets/device_cards.dart';
import '../widgets/section_header.dart';
import 'device_detail_screen.dart';

class DeviceScreen extends StatefulWidget {
  const DeviceScreen({
    required this.dashboard,
    required this.repository,
    super.key,
  });

  final AquariumDashboard dashboard;
  final AquariumRepository repository;

  @override
  State<DeviceScreen> createState() => _DeviceScreenState();
}

class _DeviceScreenState extends State<DeviceScreen> {
  final Map<String, bool> _enabledStates = {};

  @override
  Widget build(BuildContext context) {
    final devices = buildDeviceCatalog(context, widget.dashboard.latest);

    return AppBackground(
      child: SafeArea(
        bottom: false,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 28, 16, 112),
          children: [
            const _DeviceHeader(),
            const SizedBox(height: 20),
            DeviceSummaryCard(tank: widget.dashboard.tank, devices: devices),
            const SizedBox(height: 24),
            SectionHeader(
              title: context.tr('全部设备'),
              actionText: context.l10n.deviceCount(devices.length),
              showChevron: false,
              compact: true,
              onAction: () {},
            ),
            const SizedBox(height: 14),
            ...devices.map(
              (device) => Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: DeviceListTile(
                  device: device,
                  enabled: _isEnabled(device),
                  onTap: () => _openDevice(device),
                  onToggle: (value) => _setEnabled(device, value),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  bool _isEnabled(DeviceControl device) {
    return _enabledStates[device.id] ?? device.isEnabled;
  }

  void _setEnabled(DeviceControl device, bool value) {
    setState(() {
      _enabledStates[device.id] = value;
    });
  }

  void _openDevice(DeviceControl device) {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => DeviceDetailScreen(
          device: device,
          tankId: widget.dashboard.tank.id,
          initialEnabled: _isEnabled(device),
          repository: widget.repository,
        ),
      ),
    );
  }
}

class _DeviceHeader extends StatelessWidget {
  const _DeviceHeader();

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(context.tr('设备管理'), style: AppText.compactTitle),
        const SizedBox(width: 8),
        const Icon(
          CupertinoIcons.slider_horizontal_3,
          color: AppColors.primary,
        ),
      ],
    );
  }
}
