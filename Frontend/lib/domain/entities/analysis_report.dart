import '../../core/utils/json_parsing.dart';

class AnalysisReport {
  const AnalysisReport({
    required this.reportId,
    required this.tankId,
    required this.statusScore,
    required this.zh,
    required this.en,
  });

  final String reportId;
  final String tankId;
  final int statusScore;
  final AnalysisContent zh;
  final AnalysisContent en;

  String get summary => zh.summary;
  List<AnalysisAction> get actions => zh.actions;
  String get reasoning => zh.reasoning;

  AnalysisContent contentForLanguageCode(String languageCode) {
    return languageCode == 'en' ? en : zh;
  }

  factory AnalysisReport.fromJson(Map<String, dynamic> json) {
    final data = json['data'] is Map<String, dynamic>
        ? json['data'] as Map<String, dynamic>
        : const <String, dynamic>{};
    final decision = data['decision'] is Map<String, dynamic>
        ? data['decision'] as Map<String, dynamic>
        : const <String, dynamic>{};
    final legacy = AnalysisContent.fromJson(decision);
    final zh = decision['zh'] is Map<String, dynamic>
        ? AnalysisContent.fromJson(decision['zh'] as Map<String, dynamic>)
        : legacy;
    final en = decision['en'] is Map<String, dynamic>
        ? AnalysisContent.fromJson(decision['en'] as Map<String, dynamic>)
        : zh;

    return AnalysisReport(
      reportId: asString(data['report_id']),
      tankId: asString(data['tank_id']),
      statusScore: asInt(decision['status_score']),
      zh: zh,
      en: en,
    );
  }
}

class AnalysisContent {
  const AnalysisContent({
    required this.summary,
    required this.actions,
    required this.reasoning,
  });

  final String summary;
  final List<AnalysisAction> actions;
  final String reasoning;

  factory AnalysisContent.fromJson(Map<String, dynamic> json) {
    final rawActions = json['actions'];
    return AnalysisContent(
      summary: asString(json['summary']),
      actions: rawActions is List
          ? rawActions
                .whereType<Map<String, dynamic>>()
                .map(AnalysisAction.fromJson)
                .toList()
          : const [],
      reasoning: asString(json['reasoning']),
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
