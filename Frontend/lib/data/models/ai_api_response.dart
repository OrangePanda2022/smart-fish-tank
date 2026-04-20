class AiApiResponse {
  final int code;
  final String message;
  final AiApiData? data;

  AiApiResponse({required this.code, required this.message, this.data});

  factory AiApiResponse.fromJson(Map<String, dynamic> json) {
    return AiApiResponse(
      code: json['code'] as int,
      message: json['message'] as String,
      data: json['data'] != null
          ? AiApiData.fromJson(json['data'] as Map<String, dynamic>)
          : null,
    );
  }
}

class AiApiData {
  final String reportId;
  final String tankId;
  final AiApiDecision decision;

  AiApiData({
    required this.reportId,
    required this.tankId,
    required this.decision,
  });

  factory AiApiData.fromJson(Map<String, dynamic> json) {
    return AiApiData(
      reportId: json['report_id'] as String,
      tankId: json['tank_id'] as String,
      decision: AiApiDecision.fromJson(
        json['decision'] as Map<String, dynamic>,
      ),
    );
  }
}

class AiApiDecision {
  final int statusScore;
  final String summary;
  final List<AiApiAction> actions;
  final String reasoning;

  AiApiDecision({
    required this.statusScore,
    required this.summary,
    required this.actions,
    required this.reasoning,
  });

  factory AiApiDecision.fromJson(Map<String, dynamic> json) {
    return AiApiDecision(
      statusScore: json['status_score'] as int,
      summary: json['summary'] as String,
      actions: (json['actions'] as List<dynamic>)
          .map((e) => AiApiAction.fromJson(e as Map<String, dynamic>))
          .toList(),
      reasoning: json['reasoning'] as String,
    );
  }
}

class AiApiAction {
  final String device;
  final String action;

  AiApiAction({required this.device, required this.action});

  factory AiApiAction.fromJson(Map<String, dynamic> json) {
    return AiApiAction(
      device: json['device'] as String,
      action: json['action'] as String,
    );
  }
}
