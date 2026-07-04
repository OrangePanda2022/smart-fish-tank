import 'package:flutter/cupertino.dart';

import '../../core/i18n/app_language_scope.dart';
import '../../core/i18n/app_localizations.dart';
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
              title: context.tr('鱼缸档案'),
              actionText: context.tr('基础信息'),
              showChevron: false,
              compact: true,
              onAction: () {},
            ),
            const SizedBox(height: 14),
            TankProfileCard(tank: widget.dashboard.tank),
            const SizedBox(height: 24),
            SectionHeader(
              title: context.tr('偏好设置'),
              actionText: context.tr('本地设置'),
              showChevron: false,
              compact: true,
              onAction: () {},
            ),
            const SizedBox(height: 14),
            ProfileTileGroup(
              children: [
                SettingsToggleTile(
                  icon: CupertinoIcons.bell_fill,
                  title: context.tr('通知提醒'),
                  subtitle: context.tr('接收设备和鱼缸状态通知'),
                  value: _notificationEnabled,
                  onChanged: (value) {
                    setState(() => _notificationEnabled = value);
                  },
                ),
                SettingsToggleTile(
                  icon: CupertinoIcons.exclamationmark_triangle_fill,
                  title: context.tr('水质异常提醒'),
                  subtitle: context.tr('水质指标异常时及时提醒'),
                  value: _waterAlertEnabled,
                  onChanged: (value) {
                    setState(() => _waterAlertEnabled = value);
                  },
                ),
                SettingsToggleTile(
                  icon: CupertinoIcons.sparkles,
                  title: context.tr('自动分析提醒'),
                  subtitle: context.tr('提醒你定期运行 AI 分析'),
                  value: _analysisReminderEnabled,
                  onChanged: (value) {
                    setState(() => _analysisReminderEnabled = value);
                  },
                ),
                SettingsToggleTile(
                  icon: CupertinoIcons.moon_stars_fill,
                  title: context.tr('夜间免打扰'),
                  subtitle: context.tr('夜间仅保留重要异常通知'),
                  value: _quietHoursEnabled,
                  onChanged: (value) {
                    setState(() => _quietHoursEnabled = value);
                  },
                ),
              ],
            ),
            const SizedBox(height: 24),
            ProfileTileGroup(
              children: [
                ProfileMenuTile(
                  icon: CupertinoIcons.globe,
                  title: context.tr('语言'),
                  subtitle: context.l10n.isEnglish
                      ? context.tr('英文')
                      : context.tr('中文'),
                  color: AppColors.teal,
                  onTap: _showLanguagePicker,
                ),
              ],
            ),
            const SizedBox(height: 24),
            SectionHeader(
              title: context.tr('维护与支持'),
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
                  title: context.tr('关于 AquaClaw'),
                  subtitle: context.tr('智能鱼缸状态管理工具'),
                  onTap: () => _showInfo(
                    title: context.tr('关于 AquaClaw'),
                    message: context.tr('AquaClaw 用于查看鱼缸状态、设备信息和水质数据。'),
                  ),
                ),
                ProfileMenuTile(
                  icon: CupertinoIcons.doc_text_fill,
                  title: context.tr('应用版本'),
                  subtitle: '1.0.0',
                  color: AppColors.teal,
                  onTap: () => _showInfo(
                    title: context.tr('应用版本'),
                    message: context.tr('当前版本：1.0.0'),
                  ),
                ),
                ProfileMenuTile(
                  icon: CupertinoIcons.question_circle_fill,
                  title: context.tr('使用帮助'),
                  subtitle: context.tr('查看常见操作说明'),
                  color: AppColors.orange,
                  onTap: () => _showInfo(
                    title: context.tr('使用帮助'),
                    message: context.tr('你可以在首页查看鱼缸状态，在数据中心查看传感器和历史数据。'),
                  ),
                ),
                ProfileMenuTile(
                  icon: CupertinoIcons.chat_bubble_2_fill,
                  title: context.tr('问题反馈'),
                  subtitle: context.tr('记录使用中的问题和建议'),
                  color: AppColors.purple,
                  onTap: () => _showInfo(
                    title: context.tr('问题反馈'),
                    message: context.tr('反馈入口已预留，后续可以接入表单或客服渠道。'),
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
              child: Text(context.tr('知道了')),
            ),
          ],
        );
      },
    );
  }

  void _showLanguagePicker() {
    showCupertinoModalPopup<void>(
      context: context,
      builder: (context) {
        return CupertinoActionSheet(
          title: Text(context.tr('语言')),
          actions: [
            CupertinoActionSheetAction(
              onPressed: () {
                AppLanguageScope.of(context).setLocale(const Locale('zh'));
                Navigator.of(context).pop();
              },
              child: Text(context.tr('中文')),
            ),
            CupertinoActionSheetAction(
              onPressed: () {
                AppLanguageScope.of(context).setLocale(const Locale('en'));
                Navigator.of(context).pop();
              },
              child: Text(context.tr('英文')),
            ),
          ],
          cancelButton: CupertinoActionSheetAction(
            onPressed: () => Navigator.of(context).pop(),
            child: Text(context.tr('取消')),
          ),
        );
      },
    );
  }
}

class _ProfileHeader extends StatelessWidget {
  const _ProfileHeader();

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(context.tr('我的'), style: AppText.compactTitle),
        const SizedBox(width: 8),
        const Icon(CupertinoIcons.person_crop_circle, color: AppColors.primary),
      ],
    );
  }
}
