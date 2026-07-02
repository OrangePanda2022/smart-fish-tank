import '../../core/utils/json_parsing.dart';
import 'sensor_reading.dart';

class PredictionSummary {
  const PredictionSummary({
    required this.forecast,
    required this.confidence,
    required this.warning,
    required this.generatedAt,
    required this.actions,
  });

  final SensorReading forecast;
  final double confidence;
  final String warning;
  final DateTime generatedAt;
  final List<ControlAction> actions;

  bool get hasWarning => warning.isNotEmpty;

  factory PredictionSummary.fromJson(
    Map<String, dynamic> json, {
    required SensorReading fallback,
  }) {
    final data = json['data'] is Map<String, dynamic>
        ? json['data'] as Map<String, dynamic>
        : json;
    final predictions =
        data['mpc_predicted'] is List &&
            (data['mpc_predicted'] as List).isNotEmpty
        ? data['mpc_predicted'] as List
        : data['predictions'] is List
        ? data['predictions'] as List
        : const [];
    final forecast =
        predictions.isNotEmpty && predictions.last is Map<String, dynamic>
        ? SensorReading.fromPredictionPoint(
            predictions.last as Map<String, dynamic>,
            fallback: fallback,
          )
        : fallback;
    final rawActions = data['mpc_actions'];

    return PredictionSummary(
      forecast: forecast,
      confidence: asDouble(data['confidence']),
      warning: asString(data['anomaly_warning']),
      generatedAt:
          DateTime.tryParse(asString(data['generated_at'])) ?? DateTime.now(),
      actions: rawActions is List
          ? rawActions
                .whereType<Map<String, dynamic>>()
                .map(ControlAction.fromJson)
                .toList()
          : const [],
    );
  }
}

class ControlAction {
  const ControlAction({
    required this.device,
    required this.action,
    required this.value,
  });

  final String device;
  final String action;
  final double value;

  factory ControlAction.fromJson(Map<String, dynamic> json) {
    return ControlAction(
      device: asString(json['device']),
      action: asString(json['action']),
      value: asDouble(json['value']),
    );
  }
}
