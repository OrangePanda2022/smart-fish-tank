import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../view_models/device_control.dart';
import '../widgets/app_background.dart';
import '../widgets/device_cards.dart';
import '../widgets/soft_card.dart';

class DeviceDetailScreen extends StatefulWidget {
  const DeviceDetailScreen({
    required this.device,
    required this.initialEnabled,
    super.key,
  });

  final DeviceControl device;
  final bool initialEnabled;

  @override
  State<DeviceDetailScreen> createState() => _DeviceDetailScreenState();
}

class _DeviceDetailScreenState extends State<DeviceDetailScreen> {
  late bool _enabled;

  @override
  void initState() {
    super.initState();
    _enabled = widget.initialEnabled;
  }

  @override
  Widget build(BuildContext context) {
    return AppBackground(
      child: SafeArea(
        bottom: false,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 112),
          children: [
            _DeviceDetailTopBar(title: widget.device.title),
            const SizedBox(height: 16),
            DeviceDetailHeader(
              device: widget.device,
              enabled: _enabled,
              onToggle: _setEnabled,
            ),
            const SizedBox(height: 14),
            DeviceControlPanel(device: widget.device),
            const SizedBox(height: 14),
            _MaintenanceCard(device: widget.device),
          ],
        ),
      ),
    );
  }

  void _setEnabled(bool value) {
    setState(() {
      _enabled = value;
    });
  }
}

class _DeviceDetailTopBar extends StatelessWidget {
  const _DeviceDetailTopBar({required this.title});

  final String title;

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
          Text(title, style: AppText.navTitle),
        ],
      ),
    );
  }
}

class _MaintenanceCard extends StatelessWidget {
  const _MaintenanceCard({required this.device});

  final DeviceControl device;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      padding: const EdgeInsets.all(18),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('维护提醒', style: AppText.cardTitle),
          const SizedBox(height: 14),
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 38,
                height: 38,
                decoration: BoxDecoration(
                  color: AppColors.success.withAlpha(28),
                  borderRadius: BorderRadius.circular(14),
                ),
                child: const Icon(
                  CupertinoIcons.checkmark_seal_fill,
                  color: AppColors.success,
                  size: 22,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(device.maintenance, style: AppText.body),
                    const SizedBox(height: 6),
                    const Text(
                      '设备控制当前为本地预览状态，后续可接入真实控制接口。',
                      style: AppText.caption,
                    ),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
