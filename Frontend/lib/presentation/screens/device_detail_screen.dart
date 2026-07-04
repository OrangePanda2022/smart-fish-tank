import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/repositories/aquarium_repository.dart';
import '../view_models/device_control.dart';
import '../widgets/app_background.dart';
import '../widgets/device_cards.dart';
import '../widgets/mjpeg_stream_view.dart';
import '../widgets/soft_card.dart';

class DeviceDetailScreen extends StatefulWidget {
  const DeviceDetailScreen({
    required this.device,
    required this.tankId,
    required this.initialEnabled,
    required this.repository,
    super.key,
  });

  final DeviceControl device;
  final String tankId;
  final bool initialEnabled;
  final AquariumRepository repository;

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
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: AppBackground(
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
              if (widget.device.id == 'camera') ...[
                const SizedBox(height: 14),
                _CameraStreamCard(
                  streamUri: widget.repository.videoStreamUri(widget.tankId),
                ),
              ],
              const SizedBox(height: 14),
              _MaintenanceCard(device: widget.device),
            ],
          ),
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

class _CameraStreamCard extends StatelessWidget {
  const _CameraStreamCard({required this.streamUri});

  final Uri streamUri;

  @override
  Widget build(BuildContext context) {
    final hasStream = streamUri.toString().isNotEmpty;

    return SoftCard(
      padding: const EdgeInsets.all(18),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(context.tr('实时画面'), style: AppText.cardTitle),
          const SizedBox(height: 14),
          ClipRRect(
            borderRadius: BorderRadius.circular(16),
            child: AspectRatio(
              aspectRatio: 16 / 9,
              child: DecoratedBox(
                decoration: const BoxDecoration(color: AppColors.navy),
                child: hasStream
                    ? MjpegStreamView(
                        url: streamUri.toString(),
                        fit: BoxFit.cover,
                      )
                    : _StreamPlaceholder(
                        icon: CupertinoIcons.video_camera_solid,
                        title: context.tr('未连接鱼缸'),
                        subtitle: context.tr('获取鱼缸信息后可查看实时画面'),
                      ),
              ),
            ),
          ),
          const SizedBox(height: 12),
          Text(
            context.tr('直播源来自后端 MJPEG 接口，设备控制仍保持本地预览状态。'),
            style: AppText.caption,
          ),
        ],
      ),
    );
  }
}

class _StreamPlaceholder extends StatelessWidget {
  const _StreamPlaceholder({
    required this.icon,
    required this.title,
    required this.subtitle,
  });

  final IconData icon;
  final String title;
  final String subtitle;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(18),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, color: Colors.white.withAlpha(210), size: 34),
            const SizedBox(height: 10),
            Text(
              title,
              textAlign: TextAlign.center,
              style: AppText.cardTitle.copyWith(color: Colors.white),
            ),
            const SizedBox(height: 6),
            Text(
              subtitle,
              textAlign: TextAlign.center,
              style: AppText.caption.copyWith(color: Colors.white70),
            ),
          ],
        ),
      ),
    );
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
          Text(context.tr('维护提醒'), style: AppText.cardTitle),
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
                    Text(
                      context.tr('设备控制当前为本地预览状态，后续可接入真实控制接口。'),
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
