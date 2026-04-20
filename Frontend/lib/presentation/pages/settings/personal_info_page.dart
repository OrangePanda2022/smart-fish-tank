import 'package:aqua/presentation/pages/settings/settings_detail_page.dart';
import 'package:aqua/presentation/pages/settings/settings_list_tile.dart';
import 'package:flutter/material.dart';

class PersonalInfoPage extends StatelessWidget {
  const PersonalInfoPage({super.key});

  @override
  Widget build(BuildContext context) {
    return SettingsDetailPage(
      title: '个人信息',
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(20),
            ),
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                  Container(
                    width: 80,
                    height: 80,
                    decoration: BoxDecoration(
                      color: Theme.of(context).colorScheme.primaryContainer,
                      shape: BoxShape.circle,
                    ),
                    child: Icon(
                      Icons.account_circle,
                      size: 48,
                      color: Theme.of(context).colorScheme.onPrimaryContainer,
                    ),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    'Admin User',
                    style: Theme.of(context).textTheme.titleLarge,
                  ),
                  Text(
                    'admin@aquasmart.com',
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Theme.of(context).colorScheme.outline,
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
            clipBehavior: Clip.hardEdge,
            child: Column(
              children: [
                SettingsListTile(
                  title: '头像',
                  subtitle: '点击更换头像',
                  icon: Icons.face,
                  onTap: () {},
                ),
                const Divider(height: 1),
                SettingsListTile(
                  title: '昵称',
                  subtitle: 'Admin User',
                  icon: Icons.badge,
                  onTap: () {},
                ),
                const Divider(height: 1),
                SettingsListTile(
                  title: '联系方式',
                  subtitle: '138****8888',
                  icon: Icons.phone,
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
