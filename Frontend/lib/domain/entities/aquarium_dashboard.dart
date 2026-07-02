import 'prediction_summary.dart';
import 'sensor_reading.dart';
import 'tank_info.dart';

class AquariumDashboard {
  const AquariumDashboard({
    required this.tank,
    required this.latest,
    required this.prediction,
    this.score,
    required this.analysisSummary,
  });

  final TankInfo tank;
  final SensorReading latest;
  final PredictionSummary prediction;
  final int? score;
  final String analysisSummary;

  factory AquariumDashboard.empty() {
    final emptyReading = SensorReading(
      temperature: 0,
      ph: 0,
      oxygen: 0,
      ammonia: 0,
      waterLevel: 0,
      tds: 0,
      nitrate: 0,
      nitrite: 0,
      chloride: 0,
      timestamp: DateTime.fromMillisecondsSinceEpoch(0),
    );

    return AquariumDashboard(
      tank: const TankInfo(id: '', name: '', size: 0, fishCount: 0),
      latest: emptyReading,
      prediction: PredictionSummary(
        forecast: emptyReading,
        confidence: 0,
        warning: '',
        generatedAt: DateTime.fromMillisecondsSinceEpoch(0),
        actions: const [],
      ),
      score: null,
      analysisSummary: '',
    );
  }
}
