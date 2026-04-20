import 'package:aqua/domain/models/DeviceInfo.dart';
import 'package:flutter/material.dart';

class DeviceCard extends StatelessWidget {
  final Device device;
  final ValueChanged<bool>? onToggle;
  final ValueChanged<double>? onIntensityChange;
  final ValueChanged<double>? onColorTempChange;

  const DeviceCard({
    super.key,
    required this.device,
    this.onToggle,
    this.onIntensityChange,
    this.onColorTempChange
  });

  IconData getIcon() {
    switch (device.type) {
      case 'light':
        return Icons.lightbulb;
      case 'pump':
        return Icons.water;
      case 'heater':
        return Icons.thermostat;
      case 'air':
        return Icons.air;
      case 'feeder':
        return Icons.restaurant;
      default:
        return Icons.device_unknown;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Theme.of(context).colorScheme.primaryContainer,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(32)),
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 图标 + 开关
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Container(
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.onPrimaryContainer,
                    borderRadius: BorderRadius.circular(32),
                  ),
                  padding: EdgeInsets.all(12),
                  child: Icon(getIcon(), size: 24, color: Colors.white),
                ),
                Switch(
                  value: device.isOn,
                  onChanged: onToggle,
                ),
              ],
            ),
            SizedBox(height: 12),
            Text(
              device.name,
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            SizedBox(height: 8),
            if (device.intensity != null)
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text("强度: ${device.intensity!.toInt()}%"),
                  Slider(
                    year2023: false,
                    value: device.intensity!,
                    min: 0,
                    max: 100,
                    onChanged: onIntensityChange,
                  ),
                ],
              ),
            if (device.colorTemp != null)
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text("色温: ${device.colorTemp!.toInt()}K"),
                  Slider(
                    year2023: false,
                    value: device.colorTemp!,
                    min: 2000,
                    max: 10000,
                    onChanged: onColorTempChange,
                  ),
                ],
              ),
            if (device.flow != null) Text("水流: ${device.flow}"),
            if (device.setTemp != null) Text("设置: ${device.setTemp}°C"),
            if (device.schedule != null) Text("计划: ${device.schedule}"),
            if (device.nextFeed != null) Text("下一次: ${device.nextFeed}"),
          ],
        ),
      ),
    );
  }
}