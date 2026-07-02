import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../view_models/device_catalog.dart';
import '../view_models/device_control.dart';
import '../widgets/app_background.dart';
import '../widgets/device_cards.dart';
import '../widgets/section_header.dart';
import 'device_detail_screen.dart';

class DeviceScreen extends StatefulWidget {
  const DeviceScreen({required this.dashboard, super.key});

  final AquariumDashboard dashboard;

  @override
  State<DeviceScreen> createState() => _DeviceScreenState();
}

class _DeviceScreenState extends State<DeviceScreen> {
  final Map<String, bool> _enabledStates = {};

  @override
  Widget build(BuildContext context) {
    final devices = buildDeviceCatalog(widget.dashboard.latest);

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
              title: '全部设备',
              actionText: '${devices.length} 个设备',
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
          initialEnabled: _isEnabled(device),
        ),
      ),
    );
  }
}

class _DeviceHeader extends StatelessWidget {
  const _DeviceHeader();

  @override
  Widget build(BuildContext context) {
    return const Row(
      children: [
        Text('设备管理', style: AppText.compactTitle),
        SizedBox(width: 8),
        Icon(CupertinoIcons.slider_horizontal_3, color: AppColors.primary),
      ],
    );
  }
}
