import 'package:aqua/presentation/pages/settings/settings_detail_page.dart';
import 'package:flutter/material.dart';

class WaterGoalsPage extends StatefulWidget {
  const WaterGoalsPage({super.key});

  @override
  State<WaterGoalsPage> createState() => _WaterGoalsPageState();
}

class _WaterGoalsPageState extends State<WaterGoalsPage> {
  double _temperatureMin = 24.0;
  double _temperatureMax = 28.0;
  double _phMin = 6.5;
  double _phMax = 7.5;

  @override
  Widget build(BuildContext context) {
    return SettingsDetailPage(
      title: '水质目标',
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
                      Text(
                        '温度范围',
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        '${_temperatureMin.toStringAsFixed(1)}°C',
                        style: Theme.of(context).textTheme.bodyLarge,
                      ),
                      Text(
                        '${_temperatureMax.toStringAsFixed(1)}°C',
                        style: Theme.of(context).textTheme.bodyLarge,
                      ),
                    ],
                  ),
                  RangeSlider(
                    values: RangeValues(_temperatureMin, _temperatureMax),
                    min: 18.0,
                    max: 32.0,
                    divisions: 28,
                    labels: RangeLabels(
                      '${_temperatureMin.toStringAsFixed(1)}°C',
                      '${_temperatureMax.toStringAsFixed(1)}°C',
                    ),
                    onChanged: (values) {
                      setState(() {
                        _temperatureMin = values.start;
                        _temperatureMax = values.end;
                      });
                    },
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
                      Text(
                        'pH值范围',
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        _phMin.toStringAsFixed(1),
                        style: Theme.of(context).textTheme.bodyLarge,
                      ),
                      Text(
                        _phMax.toStringAsFixed(1),
                        style: Theme.of(context).textTheme.bodyLarge,
                      ),
                    ],
                  ),
                  RangeSlider(
                    values: RangeValues(_phMin, _phMax),
                    min: 5.0,
                    max: 9.0,
                    divisions: 40,
                    labels: RangeLabels(
                      _phMin.toStringAsFixed(1),
                      _phMax.toStringAsFixed(1),
                    ),
                    onChanged: (values) {
                      setState(() {
                        _phMin = values.start;
                        _phMax = values.end;
                      });
                    },
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
                      Text(
                        '水位警戒线',
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        '低水位：20%',
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                      Text(
                        '高水位：90%',
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                    ],
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
