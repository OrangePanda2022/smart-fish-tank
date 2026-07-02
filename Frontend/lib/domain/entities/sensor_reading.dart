import '../../core/utils/json_parsing.dart';

class SensorReading {
  const SensorReading({
    required this.temperature,
    required this.ph,
    required this.oxygen,
    required this.ammonia,
    required this.waterLevel,
    required this.tds,
    required this.nitrate,
    required this.nitrite,
    required this.chloride,
    required this.timestamp,
  });

  final double temperature;
  final double ph;
  final double oxygen;
  final double ammonia;
  final double waterLevel;
  final double tds;
  final double nitrate;
  final double nitrite;
  final double chloride;
  final DateTime timestamp;

  bool get isUseful => temperature > 0 && ph > 0;

  factory SensorReading.fromJson(Map<String, dynamic> json) {
    return SensorReading(
      temperature: asDouble(json['temperature']),
      ph: asDouble(json['ph']),
      oxygen: asDouble(json['oxygen']),
      ammonia: asDouble(json['ammonia']),
      waterLevel: asDouble(json['waterLevel'] ?? json['water_level']),
      tds: asDouble(json['tds']),
      nitrate: asDouble(json['nitrate']),
      nitrite: asDouble(json['nitrite']),
      chloride: asDouble(json['chlorideIon'] ?? json['chloride']),
      timestamp:
          DateTime.tryParse(asString(json['timestamp'])) ?? DateTime.now(),
    );
  }

  factory SensorReading.fromPredictionPoint(
    Map<String, dynamic> json, {
    required SensorReading fallback,
  }) {
    final state = json['state'];
    if (state is! List) return fallback;

    double stateAt(int index, double fallbackValue) {
      if (index >= state.length) return fallbackValue;
      return asDouble(state[index], fallback: fallbackValue);
    }

    return SensorReading(
      temperature: stateAt(0, fallback.temperature),
      ph: stateAt(1, fallback.ph),
      oxygen: stateAt(2, fallback.oxygen),
      ammonia: stateAt(3, fallback.ammonia),
      waterLevel: stateAt(4, fallback.waterLevel),
      tds: stateAt(5, fallback.tds),
      nitrate: stateAt(6, fallback.nitrate),
      nitrite: stateAt(7, fallback.nitrite),
      chloride: stateAt(8, fallback.chloride),
      timestamp:
          DateTime.tryParse(asString(json['time'])) ?? fallback.timestamp,
    );
  }
}
