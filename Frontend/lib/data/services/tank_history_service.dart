import 'dart:math';
import 'package:aqua/domain/models/tank_history.dart';

class TankHistoryService {
  Future<TankHistoryData> getHistory(StatType type) async {
    await Future.delayed(const Duration(milliseconds: 500));
    return _generateMockData(type);
  }

  TankHistoryData _generateMockData(StatType type) {
    final now = DateTime.now();
    final random = Random();
    final dataPoints = <HistoryDataPoint>[];

    double baseValue;
    switch (type) {
      case StatType.temperature:
        baseValue = 26.0;
        break;
      case StatType.ph:
        baseValue = 7.0;
        break;
      case StatType.waterLevel:
        baseValue = 90.0;
        break;
    }

    for (int i = 0; i < 30; i++) {
      final timestamp = now.subtract(Duration(days: 30 - i));
      double variation;
      switch (type) {
        case StatType.temperature:
          variation = (random.nextDouble() - 0.5) * 3;
          break;
        case StatType.ph:
          variation = (random.nextDouble() - 0.5) * 0.8;
          break;
        case StatType.waterLevel:
          variation = (random.nextDouble() - 0.5) * 10;
          break;
      }
      dataPoints.add(
        HistoryDataPoint(timestamp: timestamp, value: baseValue + variation),
      );
    }

    return TankHistoryData(type: type, unit: type.unit, dataPoints: dataPoints);
  }
}
