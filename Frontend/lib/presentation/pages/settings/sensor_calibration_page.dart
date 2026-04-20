import 'package:aqua/presentation/pages/settings/settings_detail_page.dart';
import 'package:flutter/material.dart';

class SensorCalibrationPage extends StatefulWidget {
  const SensorCalibrationPage({super.key});

  @override
  State<SensorCalibrationPage> createState() => _SensorCalibrationPageState();
}

class _SensorCalibrationPageState extends State<SensorCalibrationPage> {
  double _temperatureOffset = 0.0;
  double _phOffset = 0.0;

  @override
  Widget build(BuildContext context) {
    return SettingsDetailPage(
      title: '传感器校准',
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(20),
            ),
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.all(10),
                        decoration: BoxDecoration(
                          color: Theme.of(context).colorScheme.primaryContainer,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Icon(
                          Icons.thermostat,
                          color: Theme.of(
                            context,
                          ).colorScheme.onPrimaryContainer,
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              '温度传感器',
                              style: Theme.of(context).textTheme.titleMedium,
                            ),
                            Text(
                              '当前偏移值：${_temperatureOffset >= 0 ? '+' : ''}${_temperatureOffset.toStringAsFixed(1)}°C',
                              style: Theme.of(context).textTheme.bodySmall
                                  ?.copyWith(
                                    color: Theme.of(
                                      context,
                                    ).colorScheme.outline,
                                  ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      const Text('-2.0'),
                      Expanded(
                        child: Slider(
                          value: _temperatureOffset,
                          min: -2.0,
                          max: 2.0,
                          divisions: 40,
                          label:
                              '${_temperatureOffset >= 0 ? '+' : ''}${_temperatureOffset.toStringAsFixed(1)}°C',
                          onChanged: (value) {
                            setState(() => _temperatureOffset = value);
                          },
                        ),
                      ),
                      const Text('+2.0'),
                    ],
                  ),
                  Center(
                    child: TextButton(
                      onPressed: () {
                        setState(() => _temperatureOffset = 0.0);
                      },
                      child: const Text('重置'),
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          Card(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(20),
            ),
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.all(10),
                        decoration: BoxDecoration(
                          color: Theme.of(context).colorScheme.primaryContainer,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Icon(
                          Icons.science,
                          color: Theme.of(
                            context,
                          ).colorScheme.onPrimaryContainer,
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'pH传感器',
                              style: Theme.of(context).textTheme.titleMedium,
                            ),
                            Text(
                              '当前偏移值：${_phOffset >= 0 ? '+' : ''}${_phOffset.toStringAsFixed(1)}',
                              style: Theme.of(context).textTheme.bodySmall
                                  ?.copyWith(
                                    color: Theme.of(
                                      context,
                                    ).colorScheme.outline,
                                  ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      const Text('-1.0'),
                      Expanded(
                        child: Slider(
                          value: _phOffset,
                          min: -1.0,
                          max: 1.0,
                          divisions: 40,
                          label:
                              '${_phOffset >= 0 ? '+' : ''}${_phOffset.toStringAsFixed(2)}',
                          onChanged: (value) {
                            setState(() => _phOffset = value);
                          },
                        ),
                      ),
                      const Text('+1.0'),
                    ],
                  ),
                  Center(
                    child: TextButton(
                      onPressed: () {
                        setState(() => _phOffset = 0.0);
                      },
                      child: const Text('重置'),
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          Card(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(20),
            ),
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.all(10),
                        decoration: BoxDecoration(
                          color: Theme.of(context).colorScheme.primaryContainer,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Icon(
                          Icons.water,
                          color: Theme.of(
                            context,
                          ).colorScheme.onPrimaryContainer,
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              '水位传感器',
                              style: Theme.of(context).textTheme.titleMedium,
                            ),
                            Text(
                              '状态：已校准',
                              style: Theme.of(context).textTheme.bodySmall
                                  ?.copyWith(color: Colors.green),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  SizedBox(
                    width: double.infinity,
                    child: OutlinedButton(
                      onPressed: () {},
                      child: const Text('开始校准'),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
