import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/presentation/widgets/DeviceCard.dart';
import 'package:aqua/providers/device_provider.dart';

class DevicePage extends ConsumerWidget {
  const DevicePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final devices = ref.watch(devicesProvider);

    return Scaffold(
      appBar: AppBar(title: const Text("设备控制")),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: _buildDeviceCardLayout(context, ref, devices),
      ),
    );
  }

  Widget _buildDeviceCardLayout(
    BuildContext context,
    WidgetRef ref,
    List devices,
  ) {
    if (devices.isEmpty) return const SizedBox();

    List<Widget> widgets = [];

    widgets.add(
      DeviceCard(
        device: devices[0],
        onToggle: (val) =>
            ref.read(devicesProvider.notifier).toggleDevice(0, val),
        onIntensityChange: (val) =>
            ref.read(devicesProvider.notifier).updateIntensity(0, val),
        onColorTempChange: (val) =>
            ref.read(devicesProvider.notifier).updateColorTemp(0, val),
      ),
    );
    widgets.add(const SizedBox(height: 16));

    for (int i = 1; i < devices.length; i += 2) {
      final left = devices[i];
      final right = (i + 1 < devices.length) ? devices[i + 1] : null;

      widgets.add(
        Row(
          children: [
            Expanded(
              child: DeviceCard(
                device: left,
                onToggle: (val) =>
                    ref.read(devicesProvider.notifier).toggleDevice(i, val),
              ),
            ),
            const SizedBox(width: 12),
            if (right != null)
              Expanded(
                child: DeviceCard(
                  device: right,
                  onToggle: (val) => ref
                      .read(devicesProvider.notifier)
                      .toggleDevice(i + 1, val),
                ),
              )
            else
              const Expanded(child: SizedBox()),
          ],
        ),
      );
      widgets.add(const SizedBox(height: 12));
    }

    return Column(children: widgets);
  }
}
