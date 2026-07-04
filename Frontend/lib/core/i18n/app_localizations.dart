import 'package:flutter/material.dart';

class AppLocalizations {
  const AppLocalizations(this.locale);

  final Locale locale;

  static const supportedLocales = [Locale('zh'), Locale('en')];

  static const delegate = _AppLocalizationsDelegate();

  static AppLocalizations of(BuildContext context) {
    final localizations = Localizations.of<AppLocalizations>(
      context,
      AppLocalizations,
    );
    assert(localizations != null, 'No AppLocalizations found in context');
    return localizations!;
  }

  bool get isEnglish => locale.languageCode == 'en';

  String t(String zh, {String? en}) {
    if (!isEnglish) return zh;
    return en ?? _en[zh] ?? zh;
  }

  String tankLitersFish(int size, int fishCount) {
    return isEnglish
        ? '${size}L · $fishCount fish'
        : '${size}L · $fishCount 条鱼';
  }

  String tankManaged(int size, int fishCount) {
    return isEnglish
        ? '${size}L · $fishCount fish · Smart care'
        : '${size}L · $fishCount 条鱼 · 智能托管';
  }

  String deviceCount(int count) {
    return isEnglish ? '$count devices' : '$count 个设备';
  }

  String showingRows(int visible, int total) {
    return isEnglish ? 'Showing $visible/$total rows' : '显示 $visible/$total 条';
  }

  String expandRows(int count) {
    return isEnglish ? 'Expand all $count rows' : '展开全部 $count 条';
  }

  String validationRange(double min, double max) {
    return isEnglish
        ? 'Range: ${_format(min)} - ${_format(max)}'
        : '范围：${_format(min)} - ${_format(max)}';
  }

  String _format(double value) {
    if (value == value.roundToDouble()) return value.toStringAsFixed(0);
    return value.toString();
  }

  static const Map<String, String> _en = {
    '我的鱼缸': 'My Aquarium',
    '数据中心': 'Data Center',
    '历史数据': 'History',
    '首页': 'Home',
    '数据': 'Data',
    '设备': 'Devices',
    '我的': 'Me',
    '设备控制': 'Device Control',
    '全部设备': 'All Devices',
    '实时传感器数据': 'Live Sensor Data',
    '更新时间：09:41': 'Updated: 09:41',
    '手动录入数据': 'Manual Entry',
    '全部数据': 'All Data',
    '传感器数据已录入': 'Sensor data submitted',
    '录入失败：': 'Submit failed: ',
    '将以当前时间写入后端传感器接口。':
        'The reading will be written to the backend with the current time.',
    '水温': 'Water Temp',
    '温度': 'Temperature',
    'pH值': 'pH',
    '溶氧': 'Dissolved O2',
    '溶解氧': 'Dissolved Oxygen',
    '氨氮': 'Ammonia',
    '水位': 'Water Level',
    '硝酸根': 'Nitrate',
    '硝酸盐': 'Nitrate',
    '亚硝酸根': 'Nitrite',
    '亚硝酸盐': 'Nitrite',
    '氯离子': 'Chloride',
    '取消': 'Cancel',
    '提交': 'Submit',
    '请输入有效数字': 'Enter a valid number',
    '正常': 'Normal',
    '鱼缸状态  未分析': 'Aquarium Status  Not analyzed',
    '鱼缸状态  良好': 'Aquarium Status  Good',
    '设备在线': 'Devices online',
    '未分析': 'Not analyzed',
    '分': 'pts',
    '点击立即分析生成健康评分': 'Run analysis to generate a health score',
    '已喂食': 'Fed',
    '今天 08:00': 'Today 08:00',
    '设备管理': 'Device Management',
    '在线设备': 'Online',
    '运行中': 'Running',
    '需处理': 'Needs care',
    '待机': 'Standby',
    '离线': 'Offline',
    '控制与状态': 'Controls & Status',
    '实时画面': 'Live View',
    '未连接鱼缸': 'No aquarium connected',
    '获取鱼缸信息后可查看实时画面': 'Live view is available after aquarium info loads',
    '直播源来自后端 MJPEG 接口，设备控制仍保持本地预览状态。':
        'The stream uses the backend MJPEG endpoint. Device control remains local preview.',
    '维护提醒': 'Maintenance',
    '设备控制当前为本地预览状态，后续可接入真实控制接口。':
        'Device controls are local previews and can be connected to real control APIs later.',
    '鱼缸档案': 'Aquarium Profile',
    '基础信息': 'Basic Info',
    '偏好设置': 'Preferences',
    '本地设置': 'Local',
    '通知提醒': 'Notifications',
    '接收设备和鱼缸状态通知': 'Receive device and aquarium status alerts',
    '水质异常提醒': 'Water Quality Alerts',
    '水质指标异常时及时提醒': 'Alert when water metrics are abnormal',
    '自动分析提醒': 'Analysis Reminder',
    '提醒你定期运行 AI 分析': 'Remind you to run AI analysis regularly',
    '夜间免打扰': 'Quiet Hours',
    '夜间仅保留重要异常通知': 'Only important alerts at night',
    '语言': 'Language',
    '中文': 'Chinese',
    '英文': 'English',
    '维护与支持': 'Maintenance & Support',
    '关于 AquaClaw': 'About AquaClaw',
    '智能鱼缸状态管理工具': 'Smart aquarium status manager',
    'AquaClaw 用于查看鱼缸状态、设备信息和水质数据。':
        'AquaClaw helps you view aquarium status, devices, and water quality data.',
    '应用版本': 'App Version',
    '当前版本：1.0.0': 'Current version: 1.0.0',
    '使用帮助': 'Help',
    '查看常见操作说明': 'View common usage tips',
    '你可以在首页查看鱼缸状态，在数据中心查看传感器和历史数据。':
        'Use Home for aquarium status and Data Center for sensor and history data.',
    '问题反馈': 'Feedback',
    '记录使用中的问题和建议': 'Record issues and suggestions',
    '反馈入口已预留，后续可以接入表单或客服渠道。':
        'Feedback entry is reserved for a future form or support channel.',
    '知道了': 'OK',
    '鱼缸 ID': 'Aquarium ID',
    '未设置': 'Not set',
    '鱼缸名称': 'Aquarium Name',
    '鱼缸容量': 'Capacity',
    '鱼只数量': 'Fish Count',
    'AI 趋势预测': 'AI Trend Forecast',
    '未来24小时': 'Next 24 hours',
    'AI 建议': 'AI Advice',
    '未来趋势稳定，保持当前设备策略': 'Trend is stable. Keep the current device strategy.',
    'AI 智能分析': 'AI Analysis',
    '上次分析：今天 09:30': 'Last analysis: Today 09:30',
    '水质正常': 'Water normal',
    '温度适宜': 'Temperature ideal',
    '鱼儿活跃': 'Fish active',
    '分析中...': 'Analyzing...',
    '立即分析': 'Analyze Now',
    '状态概要': 'Status Summary',
    '建议操作': 'Recommended Actions',
    '暂无建议操作': 'No recommended actions',
    '推理过程': 'Reasoning',
    'AI 分析失败': 'AI Analysis Failed',
    '分析请求失败或超时，请确认后端服务可用后重试。':
        'Analysis failed or timed out. Check the backend service and try again.',
    '打印报告': 'Print Report',
    '生成中...': 'Generating...',
    '选择打印语言': 'Report Language',
    '打印中文版': 'Print Chinese',
    '打印英文版': 'Print English',
    'AI 分析结果': 'AI Analysis Result',
    '暂无概要': 'No summary',
    '更多功能': 'More',
    '添加设备': 'Add Device',
    '灯光': 'Light',
    '水泵': 'Pump',
    '氧气': 'Aerator',
    '加热': 'Heater',
    '喂食': 'Feeder',
    '换水': 'Water Change',
    '摄像头': 'Camera',
    '过滤': 'Filter',
    '备用插座': 'Spare Outlet',
    '模拟自然昼夜节律，当前处于日间补光。':
        'Simulates a natural day-night rhythm. Daylight fill is active.',
    '调节亮度': 'Adjust Brightness',
    '建议每天 08:00-20:00 开启': 'Recommended on from 08:00 to 20:00',
    '亮度': 'Brightness',
    '模式': 'Mode',
    '日光': 'Daylight',
    '定时': 'Schedule',
    '20:00 关闭': 'Off at 20:00',
    '保持水体循环，辅助过滤系统稳定工作。':
        'Keeps water circulating and helps the filter stay stable.',
    '调整流量': 'Adjust Flow',
    '滤芯建议 14 天后检查': 'Check filter cartridge in 14 days',
    '流量': 'Flow',
    '中档': 'Medium',
    '运行': 'Run',
    '连续': 'Continuous',
    '根据溶氧水平维持鱼缸含氧量。':
        'Maintains oxygen level based on dissolved oxygen readings.',
    '增强供氧': 'Boost Oxygen',
    '气石状态良好': 'Air stone is healthy',
    '强度': 'Intensity',
    '智能': 'Smart',
    '自动维持适合观赏鱼的水温区间。':
        'Maintains a suitable water temperature for ornamental fish.',
    '设置温度': 'Set Temperature',
    '目标温度 26.0°C': 'Target temperature 26.0°C',
    '当前水温': 'Current Temp',
    '目标': 'Target',
    '功率': 'Power',
    '08:00 已喂食': 'Fed at 08:00',
    '自动投喂已按计划完成。': 'Automatic feeding completed as scheduled.',
    '立即投喂': 'Feed Now',
    '饲料余量充足': 'Food supply is sufficient',
    '今日次数': 'Today',
    '下次': 'Next',
    '份量': 'Portion',
    '标准': 'Standard',
    '剩余 7 天': '7 days left',
    '记录换水周期，提醒维持稳定水质。':
        'Tracks water-change cycles to help maintain stable water quality.',
    '记录换水': 'Log Water Change',
    '下次换水 7 天后': 'Next water change in 7 days',
    '周期': 'Cycle',
    '14 天': '14 days',
    '上次': 'Last',
    '7 天前': '7 days ago',
    '建议量': 'Suggested',
    '在线': 'Online',
    '远程查看鱼缸画面和设备状态。': 'View aquarium video and device status remotely.',
    '查看画面': 'View Camera',
    '网络连接稳定': 'Network is stable',
    '清晰度': 'Quality',
    '延迟': 'Latency',
    '低': 'Low',
    '存储': 'Storage',
    '本地': 'Local',
    '过滤杂质并辅助稳定氨氮、硝酸盐指标。':
        'Filters impurities and helps stabilize ammonia and nitrate metrics.',
    '切换模式': 'Switch Mode',
    '建议 3 天后清洗滤棉': 'Clean filter cotton in 3 days',
    '自动': 'Auto',
    '稳定': 'Stable',
    '滤棉': 'Filter Cotton',
    '良好': 'Good',
    '可接入扩展设备，当前未连接。': 'Reserved for expansion devices. Currently disconnected.',
    '重新连接': 'Reconnect',
    '检查电源与网络': 'Check power and network',
    '状态': 'Status',
    '信号': 'Signal',
    '无': 'None',
    '负载': 'Load',
    '加热设备': 'heater',
    '增氧设备': 'aerator',
    '喂食器': 'feeder',
    '降温': ' cool down',
    '升温': ' heat up',
    '开启': ' turn on',
    '关闭': ' turn off',
    '调暗': ' dim',
    '调亮': ' brighten',
    '保持': ' hold',
    '使用中文标题、标签与页脚生成报告':
        'Use Chinese headings, labels, and footer for the report',
  };
}

extension AppLocalizationsX on BuildContext {
  AppLocalizations get l10n => AppLocalizations.of(this);

  String tr(String zh, {String? en}) => l10n.t(zh, en: en);
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) {
    return AppLocalizations.supportedLocales.any(
      (supported) => supported.languageCode == locale.languageCode,
    );
  }

  @override
  Future<AppLocalizations> load(Locale locale) async {
    final languageCode = locale.languageCode == 'en' ? 'en' : 'zh';
    return AppLocalizations(Locale(languageCode));
  }

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}
