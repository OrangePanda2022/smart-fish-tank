import 'package:flutter/cupertino.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../widgets/app_background.dart';
import '../widgets/profile_cards.dart';
import '../widgets/section_header.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({required this.dashboard, super.key});

  final AquariumDashboard dashboard;

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  bool _notificationEnabled = true;
  bool _waterAlertEnabled = true;
  bool _analysisReminderEnabled = false;
  bool _quietHoursEnabled = true;

  @override
  Widget build(BuildContext context) {
    return AppBackground(
      child: SafeArea(
        bottom: false,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 28, 16, 112),
          children: [
            const _ProfileHeader(),
            const SizedBox(height: 20),
            ProfileHeaderCard(tank: widget.dashboard.tank),
            const SizedBox(height: 24),
            SectionHeader(
              title: '鱼缸档案',
              actionText: '基础信息',
              showChevron: false,
              compact: true,
              onAction: () {},
            ),
            const SizedBox(height: 14),
            TankProfileCard(tank: widget.dashboard.tank),
            const SizedBox(height: 24),
            SectionHeader(
              title: '偏好设置',
              actionText: '本地设置',
              showChevron: false,
              compact: true,
              onAction: () {},
            ),
            const SizedBox(height: 14),
            ProfileTileGroup(
              children: [
                SettingsToggleTile(
                  icon: CupertinoIcons.bell_fill,
                  title: '通知提醒',
                  subtitle: '接收设备和鱼缸状态通知',
                  value: _notificationEnabled,
                  onChanged: (value) {
                    setState(() => _notificationEnabled = value);
                  },
                ),
                SettingsToggleTile(
                  icon: CupertinoIcons.exclamationmark_triangle_fill,
                  title: '水质异常提醒',
                  subtitle: '水质指标异常时及时提醒',
                  value: _waterAlertEnabled,
                  onChanged: (value) {
                    setState(() => _waterAlertEnabled = value);
                  },
                ),
                SettingsToggleTile(
                  icon: CupertinoIcons.sparkles,
                  title: '自动分析提醒',
                  subtitle: '提醒你定期运行 AI 分析',
                  value: _analysisReminderEnabled,
                  onChanged: (value) {
                    setState(() => _analysisReminderEnabled = value);
                  },
                ),
                SettingsToggleTile(
                  icon: CupertinoIcons.moon_stars_fill,
                  title: '夜间免打扰',
                  subtitle: '夜间仅保留重要异常通知',
                  value: _quietHoursEnabled,
                  onChanged: (value) {
                    setState(() => _quietHoursEnabled = value);
                  },
                ),
              ],
            ),
            const SizedBox(height: 24),
            SectionHeader(
              title: '维护与支持',
              actionText: 'AquaClaw',
              showChevron: false,
              compact: true,
              onAction: () {},
            ),
            const SizedBox(height: 14),
            ProfileTileGroup(
              children: [
                ProfileMenuTile(
                  icon: CupertinoIcons.info_circle_fill,
                  title: '关于 AquaClaw',
                  subtitle: '智能鱼缸状态管理工具',
                  onTap: () => _showInfo(
                    title: '关于 AquaClaw',
                    message: 'AquaClaw 用于查看鱼缸状态、设备信息和水质数据。',
                  ),
                ),
                ProfileMenuTile(
                  icon: CupertinoIcons.doc_text_fill,
                  title: '应用版本',
                  subtitle: '1.0.0',
                  color: AppColors.teal,
                  onTap: () => _showInfo(title: '应用版本', message: '当前版本：1.0.0'),
                ),
                ProfileMenuTile(
                  icon: CupertinoIcons.question_circle_fill,
                  title: '使用帮助',
                  subtitle: '查看常见操作说明',
                  color: AppColors.orange,
                  onTap: () => _showInfo(
                    title: '使用帮助',
                    message: '你可以在首页查看鱼缸状态，在数据中心查看传感器和历史数据。',
                  ),
                ),
                ProfileMenuTile(
                  icon: CupertinoIcons.chat_bubble_2_fill,
                  title: '问题反馈',
                  subtitle: '记录使用中的问题和建议',
                  color: AppColors.purple,
                  onTap: () => _showInfo(
                    title: '问题反馈',
                    message: '反馈入口已预留，后续可以接入表单或客服渠道。',
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _showInfo({required String title, required String message}) {
    showCupertinoDialog<void>(
      context: context,
      builder: (context) {
        return CupertinoAlertDialog(
          title: Text(title),
          content: Padding(
            padding: const EdgeInsets.only(top: 8),
            child: Text(message),
          ),
          actions: [
            CupertinoDialogAction(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('知道了'),
            ),
          ],
        );
      },
    );
  }
}

class _ProfileHeader extends StatelessWidget {
  const _ProfileHeader();

  @override
  Widget build(BuildContext context) {
    return const Row(
      children: [
        Text('我的', style: AppText.compactTitle),
        SizedBox(width: 8),
        Icon(CupertinoIcons.person_crop_circle, color: AppColors.primary),
      ],
    );
  }
}
