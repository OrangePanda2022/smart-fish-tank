import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/domain/models/SettingModel.dart';
import 'package:aqua/providers/settings_provider.dart';
import 'settings/personal_info_page.dart';
import 'settings/account_security_page.dart';
import 'settings/notification_page.dart';
import 'settings/tank_profile_page.dart';
import 'settings/water_goals_page.dart';
import 'settings/sensor_calibration_page.dart';
import 'settings/appearance_page.dart';
import 'settings/about_page.dart';

class SettingsPage extends ConsumerWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final settingGroups = ref.watch(settingGroupsProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('设置')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          ...settingGroups.map((group) {
            return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (group.title.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 8),
                    child: Text(
                      group.title,
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                  ),
                Card(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(20),
                  ),
                  clipBehavior: Clip.hardEdge,
                  elevation: 0,
                  child: Column(
                    children: List.generate(group.items.length, (index) {
                      final item = group.items[index];
                      return Column(
                        children: [
                          _buildItem(context, ref, item),
                          if (index != group.items.length - 1)
                            const Divider(height: 1, color: Colors.transparent),
                        ],
                      );
                    }),
                  ),
                ),
                const SizedBox(height: 16),
              ],
            );
          }),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8.0),
            child: ElevatedButton.icon(
              style: ElevatedButton.styleFrom(
                backgroundColor: Theme.of(context).colorScheme.errorContainer,
                foregroundColor: Theme.of(context).colorScheme.onErrorContainer,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
                minimumSize: const Size.fromHeight(50),
                elevation: 0,
              ),
              icon: const Icon(Icons.logout),
              label: const Text('退出登录'),
              onPressed: () {},
            ),
          ),
          const SizedBox(height: 20),
        ],
      ),
    );
  }

  Widget _buildItem(BuildContext context, WidgetRef ref, SettingItem item) {
    switch (item.type) {
      case SettingItemType.switchTile:
        return SwitchListTile(
          title: Text(item.title),
          subtitle: item.subtitle != null ? Text(item.subtitle!) : null,
          secondary: Icon(item.icon),
          value: item.value ?? false,
          onChanged: (val) {
            item.value = val;
            ref.read(settingGroupsProvider.notifier);
          },
        );
      case SettingItemType.normal:
        return ListTile(
          contentPadding: const EdgeInsets.symmetric(
            horizontal: 16.0,
            vertical: 8.0,
          ),
          tileColor: Theme.of(context).colorScheme.primaryContainer,
          leading: Icon(
            item.icon,
            color: Theme.of(context).colorScheme.onPrimaryContainer,
          ),
          title: Text(item.title),
          subtitle: item.subtitle != null ? Text(item.subtitle!) : null,
          trailing: const Icon(Icons.chevron_right),
          onTap: () => _navigateToSetting(context, item.title),
        );
    }
  }

  void _navigateToSetting(BuildContext context, String title) {
    switch (title) {
      case '个人信息':
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => const PersonalInfoPage()),
        );
        break;
      case '账号安全':
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => const AccountSecurityPage()),
        );
        break;
      case '消息通知':
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => const NotificationPage()),
        );
        break;
      case '鱼缸档案':
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => const TankProfilePage()),
        );
        break;
      case '水质目标':
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => const WaterGoalsPage()),
        );
        break;
      case '传感器校准':
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => const SensorCalibrationPage(),
          ),
        );
        break;
      case '外观设置':
        showModalBottomSheet(
          context: context,
          backgroundColor: Colors.transparent,
          builder: (context) => const AppearancePage(),
        );
        break;
      case '关于 AquaSmart':
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => const AboutPage()),
        );
        break;
    }
  }
}
