import '../entities/analysis_report.dart';
import '../entities/aquarium_dashboard.dart';
import '../entities/sensor_reading.dart';

abstract class AquariumRepository {
  Future<AquariumDashboard> loadDashboard();

  Future<List<SensorReading>> loadHistory(String tankId);

  Future<SensorReading> submitSensorReading(
    String tankId,
    SensorReading reading,
  );

  Future<AnalysisReport> runAnalysis(String tankId);

  Uri videoStreamUri(String tankId);
}
