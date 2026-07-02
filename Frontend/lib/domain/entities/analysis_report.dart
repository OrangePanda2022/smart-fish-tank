import '../../core/utils/json_parsing.dart';

class AnalysisReport {
  const AnalysisReport({
    required this.reportId,
    required this.tankId,
    required this.statusScore,
    required this.summary,
    required this.actions,
    required this.reasoning,
  });

  final String reportId;
  final String tankId;
  final int statusScore;
  final String summary;
  final List<AnalysisAction> actions;
  final String reasoning;

  factory AnalysisReport.fromJson(Map<String, dynamic> json) {
    final data = json['data'] is Map<String, dynamic>
        ? json['data'] as Map<String, dynamic>
        : const <String, dynamic>{};
    final decision = data['decision'] is Map<String, dynamic>
        ? data['decision'] as Map<String, dynamic>
        : const <String, dynamic>{};
    final rawActions = decision['actions'];

    return AnalysisReport(
      reportId: asString(data['report_id']),
      tankId: asString(data['tank_id']),
      statusScore: asInt(decision['status_score']),
      summary: asString(decision['summary']),
      actions: rawActions is List
          ? rawActions
                .whereType<Map<String, dynamic>>()
                .map(AnalysisAction.fromJson)
                .toList()
          : const [],
      reasoning: asString(decision['reasoning']),
    );
  }

  factory AnalysisReport.fallback({required String tankId, String? summary}) {
    return AnalysisReport(
      reportId: 'demo-report',
      tankId: tankId,
      statusScore: 95,
      summary: summary ?? '水质状况良好，各项指标均在正常范围内。',
      actions: const [
        AnalysisAction(device: 'heater', action: '维持当前温度设置'),
        AnalysisAction(device: 'filter', action: '保持当前过滤频率'),
      ],
      reasoning: summary ?? '当前传感器指标处于安全范围，建议继续观察水温、pH 和溶解氧变化。',
    );
  }
}

class AnalysisAction {
  const AnalysisAction({required this.device, required this.action});

  final String device;
  final String action;

  factory AnalysisAction.fromJson(Map<String, dynamic> json) {
    return AnalysisAction(
      device: asString(json['device']),
      action: asString(json['action']),
    );
  }
}
