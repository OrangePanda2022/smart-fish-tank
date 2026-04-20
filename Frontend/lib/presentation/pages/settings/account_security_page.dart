import 'package:aqua/presentation/pages/settings/settings_detail_page.dart';
import 'package:aqua/presentation/pages/settings/settings_list_tile.dart';
import 'package:flutter/material.dart';

class AccountSecurityPage extends StatelessWidget {
  const AccountSecurityPage({super.key});

  @override
  Widget build(BuildContext context) {
    return SettingsDetailPage(
      title: '账号安全',
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
                SettingsListTile(
                  title: '登录密码',
                  subtitle: '上次修改：2024-01-15',
                  icon: Icons.lock,
                  onTap: () {},
                ),
                const Divider(height: 1),
                SettingsListTile(
                  title: '修改密码',
                  subtitle: '定期修改密码可保护账号安全',
                  icon: Icons.key,
                  onTap: () {},
                ),
                const Divider(height: 1),
                SettingsListTile(
                  title: '设备管理',
                  subtitle: '已连接 3 台设备',
                  icon: Icons.devices,
                  onTap: () {},
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),
          Card(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(20),
            ),
            clipBehavior: Clip.hardEdge,
            child: Column(
              children: [
                SettingsListTile(
                  title: '两步验证',
                  subtitle: '未开启',
                  icon: Icons.verified_user,
                  trailing: Switch(value: false, onChanged: (value) {}),
                  onTap: () {},
                ),
                const Divider(height: 1),
                SettingsListTile(
                  title: '登录通知',
                  subtitle: '有新设备登录时通知',
                  icon: Icons.notifications,
                  trailing: Switch(value: true, onChanged: (value) {}),
                  onTap: () {},
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
