import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/domain/models/SettingModel.dart';

final settingGroupsProvider =
    NotifierProvider<SettingGroupsNotifier, List<SettingGroup>>(
      SettingGroupsNotifier.new,
    );

class SettingGroupsNotifier extends Notifier<List<SettingGroup>> {
  @override
  List<SettingGroup> build() {
    return [
      SettingGroup(
        title: '',
        items: [
          SettingItem(
            title: 'Admin User',
            subtitle: 'admin@aquasmart.com',
            icon: Icons.account_circle,
            type: SettingItemType.normal,
          ),
        ],
      ),
      SettingGroup(
        title: '账户设置',
        items: [
          SettingItem(
            title: '个人信息',
            subtitle: '修改头像、昵称、联系方式',
            icon: Icons.person,
            type: SettingItemType.normal,
          ),
          SettingItem(
            title: '账号安全',
            subtitle: '密码修改、设备管理',
            icon: Icons.shield,
            type: SettingItemType.normal,
          ),
          SettingItem(
            title: '消息通知',
            subtitle: '报警与提醒设置',
            icon: Icons.notifications,
            type: SettingItemType.normal,
          ),
        ],
      ),
      SettingGroup(
        title: '鱼缸设置',
        items: [
          SettingItem(
            title: '鱼缸档案',
            subtitle: '名称：主缸 / 容量：120L',
            icon: Icons.set_meal,
            type: SettingItemType.normal,
          ),
          SettingItem(
            title: '水质目标',
            subtitle: '设置温度、pH报警阈值',
            icon: Icons.water_drop,
            type: SettingItemType.normal,
          ),
          SettingItem(
            title: '传感器校准',
            subtitle: '温度计、pH计校准',
            icon: Icons.straighten,
            type: SettingItemType.normal,
          ),
        ],
      ),
      SettingGroup(
        title: '通用',
        items: [
          SettingItem(
            title: '外观设置',
            subtitle: '深色模式、主题色',
            icon: Icons.palette,
            type: SettingItemType.normal,
          ),
          SettingItem(
            title: '关于 AquaSmart',
            subtitle: '版本 1.0.0',
            icon: Icons.info,
            type: SettingItemType.normal,
          ),
        ],
      ),
    ];
  }
}
