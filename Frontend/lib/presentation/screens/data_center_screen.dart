import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/i18n/app_localizations.dart';
import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/entities/sensor_reading.dart';
import '../../domain/repositories/aquarium_repository.dart';
import '../widgets/app_background.dart';
import '../widgets/headers.dart';
import '../widgets/prediction_card.dart';
import '../widgets/section_header.dart';
import '../widgets/sensor_metric_grid.dart';
import '../widgets/wide_action_button.dart';
import 'history_screen.dart';

class DataCenterScreen extends StatelessWidget {
  const DataCenterScreen({
    required this.dashboard,
    required this.repository,
    required this.onRefresh,
    super.key,
  });

  final AquariumDashboard dashboard;
  final AquariumRepository repository;
  final VoidCallback onRefresh;

  @override
  Widget build(BuildContext context) {
    return AppBackground(
      child: SafeArea(
        bottom: false,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(12, 28, 12, 112),
          children: [
            const DataHeader(),
            const SizedBox(height: 22),
            PredictionCard(dashboard: dashboard),
            const SizedBox(height: 26),
            SectionHeader(
              title: context.tr('实时传感器数据'),
              actionText: context.tr('更新时间：09:41'),
              showChevron: false,
              compact: true,
              onAction: () {},
            ),
            const SizedBox(height: 14),
            SensorMetricGrid(sensor: dashboard.latest),
            const SizedBox(height: 18),
            WideActionButton(
              text: context.tr('手动录入数据'),
              icon: CupertinoIcons.plus_circle_fill,
              onTap: () => _openManualInput(context),
            ),
            const SizedBox(height: 12),
            WideActionButton(
              text: context.tr('全部数据'),
              icon: CupertinoIcons.chevron_right,
              onTap: () {
                Navigator.of(context).push(
                  MaterialPageRoute<void>(
                    builder: (_) => HistoryScreen(
                      dashboard: dashboard,
                      repository: repository,
                    ),
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _openManualInput(BuildContext context) async {
    final reading = await showDialog<SensorReading>(
      context: context,
      builder: (_) => _ManualSensorDialog(initial: dashboard.latest),
    );
    if (reading == null || !context.mounted) return;

    final messenger = ScaffoldMessenger.of(context);
    try {
      await repository.submitSensorReading(dashboard.tank.id, reading);
      if (!context.mounted) return;
      messenger.showSnackBar(SnackBar(content: Text(context.tr('传感器数据已录入'))));
      onRefresh();
    } catch (error) {
      if (!context.mounted) return;
      messenger.showSnackBar(
        SnackBar(content: Text('${context.tr('录入失败：')}$error')),
      );
    }
  }
}

class _ManualSensorDialog extends StatefulWidget {
  const _ManualSensorDialog({required this.initial});

  final SensorReading initial;

  @override
  State<_ManualSensorDialog> createState() => _ManualSensorDialogState();
}

class _ManualSensorDialogState extends State<_ManualSensorDialog> {
  final _formKey = GlobalKey<FormState>();
  late final Map<String, TextEditingController> _controllers;

  @override
  void initState() {
    super.initState();
    final initial = widget.initial;
    _controllers = {
      'temperature': _controller(initial.temperature),
      'ph': _controller(initial.ph),
      'oxygen': _controller(initial.oxygen),
      'ammonia': _controller(initial.ammonia),
      'waterLevel': _controller(initial.waterLevel),
      'tds': _controller(initial.tds),
      'nitrate': _controller(initial.nitrate),
      'nitrite': _controller(initial.nitrite),
      'chloride': _controller(initial.chloride),
    };
  }

  @override
  void dispose() {
    for (final controller in _controllers.values) {
      controller.dispose();
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(context.tr('手动录入数据')),
      contentPadding: const EdgeInsets.fromLTRB(24, 12, 24, 0),
      content: SizedBox(
        width: 420,
        child: Form(
          key: _formKey,
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(context.tr('将以当前时间写入后端传感器接口。'), style: AppText.caption),
                const SizedBox(height: 14),
                _NumberField(
                  controller: _controllers['temperature']!,
                  label: context.tr('水温'),
                  suffix: '°C',
                  min: -50,
                  max: 100,
                ),
                _NumberField(
                  controller: _controllers['ph']!,
                  label: 'pH',
                  min: 0,
                  max: 14,
                ),
                _NumberField(
                  controller: _controllers['oxygen']!,
                  label: context.tr('溶解氧'),
                  suffix: 'mg/L',
                  min: 0,
                  max: 20,
                ),
                _NumberField(
                  controller: _controllers['ammonia']!,
                  label: context.tr('氨氮'),
                  suffix: 'mg/L',
                  min: 0,
                  max: 10,
                ),
                _NumberField(
                  controller: _controllers['waterLevel']!,
                  label: context.tr('水位'),
                  suffix: '%',
                  min: 0,
                  max: 100,
                ),
                _NumberField(
                  controller: _controllers['tds']!,
                  label: 'TDS',
                  suffix: 'ppm',
                  min: 0,
                  max: 2000,
                ),
                _NumberField(
                  controller: _controllers['nitrate']!,
                  label: context.tr('硝酸根'),
                  suffix: 'mg/L',
                  min: 0,
                  max: 200,
                ),
                _NumberField(
                  controller: _controllers['nitrite']!,
                  label: context.tr('亚硝酸根'),
                  suffix: 'mg/L',
                  min: 0,
                  max: 10,
                ),
                _NumberField(
                  controller: _controllers['chloride']!,
                  label: context.tr('氯离子'),
                  suffix: 'mg/L',
                  min: 0,
                  max: 1000,
                ),
              ],
            ),
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(context.tr('取消')),
        ),
        FilledButton(onPressed: _submit, child: Text(context.tr('提交'))),
      ],
    );
  }

  TextEditingController _controller(double value) {
    return TextEditingController(text: value.toStringAsFixed(2));
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;

    Navigator.of(context).pop(
      SensorReading(
        temperature: _value('temperature'),
        ph: _value('ph'),
        oxygen: _value('oxygen'),
        ammonia: _value('ammonia'),
        waterLevel: _value('waterLevel'),
        tds: _value('tds'),
        nitrate: _value('nitrate'),
        nitrite: _value('nitrite'),
        chloride: _value('chloride'),
        timestamp: DateTime.now(),
      ),
    );
  }

  double _value(String key) {
    return double.parse(_controllers[key]!.text.trim());
  }
}

class _NumberField extends StatelessWidget {
  const _NumberField({
    required this.controller,
    required this.label,
    required this.min,
    required this.max,
    this.suffix = '',
  });

  final TextEditingController controller;
  final String label;
  final double min;
  final double max;
  final String suffix;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: TextFormField(
        controller: controller,
        keyboardType: const TextInputType.numberWithOptions(decimal: true),
        decoration: InputDecoration(
          labelText: label,
          suffixText: suffix.isEmpty ? null : suffix,
          filled: true,
          fillColor: AppColors.background,
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: BorderSide.none,
          ),
        ),
        validator: (value) {
          final number = double.tryParse(value?.trim() ?? '');
          if (number == null) return context.tr('请输入有效数字');
          if (number < min || number > max) {
            return context.l10n.validationRange(min, max);
          }
          return null;
        },
      ),
    );
  }
}
