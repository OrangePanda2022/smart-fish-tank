import 'package:aqua/presentation/pages/settings/settings_detail_page.dart';
import 'package:flutter/material.dart';

class NotificationPage extends StatefulWidget {
  const NotificationPage({super.key});

  @override
  State<NotificationPage> createState() => _NotificationPageState();
}

class _NotificationPageState extends State<NotificationPage> {
  bool _alarmEnabled = true;
  bool _reminderEnabled = true;
  bool _systemEnabled = true;
  bool _temperatureAlarm = true;
  bool _phAlarm = true;
  bool _waterLevelAlarm = true;

  @override
  Widget build(BuildContext context) {
    return SettingsDetailPage(
      title: '消息通知',
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(20),
            ),
            clipBehavior: Clip.hardEdge,
            child: Column(
              children: [
                SwitchListTile(
                  title: const Text('报警通知'),
                  subtitle: const Text('鱼缸异常时推送通知'),
                  secondary: Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      color: Theme.of(context).colorScheme.primaryContainer,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      Icons.warning,
                      color: Theme.of(context).colorScheme.onPrimaryContainer,
                    ),
                  ),
                  value: _alarmEnabled,
                  onChanged: (value) => setState(() => _alarmEnabled = value),
                ),
                const Divider(height: 1),
                SwitchListTile(
                  title: const Text('定时提醒'),
                  subtitle: const Text('喂食、换水等日常提醒'),
                  secondary: Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      color: Theme.of(context).colorScheme.primaryContainer,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      Icons.alarm,
                      color: Theme.of(context).colorScheme.onPrimaryContainer,
                    ),
                  ),
                  value: _reminderEnabled,
                  onChanged: (value) =>
                      setState(() => _reminderEnabled = value),
                ),
                const Divider(height: 1),
                SwitchListTile(
                  title: const Text('系统通知'),
                  subtitle: const Text('版本更新、活动公告'),
                  secondary: Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      color: Theme.of(context).colorScheme.primaryContainer,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      Icons.info,
                      color: Theme.of(context).colorScheme.onPrimaryContainer,
                    ),
                  ),
                  value: _systemEnabled,
                  onChanged: (value) => setState(() => _systemEnabled = value),
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),
          if (_alarmEnabled) ...[
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 8),
              child: Text(
                '报警类型',
                style: Theme.of(context).textTheme.titleMedium,
              ),
            ),
            Card(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(20),
              ),
              clipBehavior: Clip.hardEdge,
              child: Column(
                children: [
                  SwitchListTile(
                    title: const Text('温度报警'),
                    subtitle: const Text('水温超出设定范围时通知'),
                    secondary: const Icon(Icons.thermostat),
                    value: _temperatureAlarm,
                    onChanged: (value) =>
                        setState(() => _temperatureAlarm = value),
                  ),
                  const Divider(height: 1),
                  SwitchListTile(
                    title: const Text('pH值报警'),
                    subtitle: const Text('水质酸碱度异常时通知'),
                    secondary: const Icon(Icons.science),
                    value: _phAlarm,
                    onChanged: (value) => setState(() => _phAlarm = value),
                  ),
                  const Divider(height: 1),
                  SwitchListTile(
                    title: const Text('水位报警'),
                    subtitle: const Text('水位过低或过高时通知'),
                    secondary: const Icon(Icons.water),
                    value: _waterLevelAlarm,
                    onChanged: (value) =>
                        setState(() => _waterLevelAlarm = value),
                  ),
                ],
              ),
            ),
          ],
        ],
      ),
    );
  }
}
