import 'package:aquaclaw/infrastructure/repositories/aqua_api_repository.dart';
import 'package:aquaclaw/domain/entities/analysis_report.dart';
import 'package:aquaclaw/domain/entities/sensor_reading.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test(
    'runAnalysis does not create a demo report when tank is disconnected',
    () {
      final repository = AquaApiRepository();

      expect(
        () => repository.runAnalysis(''),
        throwsA(
          isA<StateError>().having(
            (error) => error.message,
            'message',
            '未连接鱼缸，无法运行 AI 分析。',
          ),
        ),
      );
    },
  );

  test('videoStreamUri builds the tank MJPEG stream endpoint', () {
    final repository = AquaApiRepository(baseUrl: 'http://localhost:8080');

    expect(
      repository.videoStreamUri('tank-1').toString(),
      'http://localhost:8080/api/v1/tank/tank-1/stream',
    );
  });

  test('AnalysisReport parses bilingual analyse response', () {
    final report = AnalysisReport.fromJson({
      'data': {
        'report_id': 'REP_1',
        'tank_id': 'tank-1',
        'decision': {
          'status_score': 86,
          'zh': {
            'summary': '水质稳定',
            'actions': [
              {'device': '水泵', 'action': '保持运行'},
            ],
            'reasoning': '鱼群活跃，传感器指标正常。',
          },
          'en': {
            'summary': 'Water quality is stable',
            'actions': [
              {'device': 'pump', 'action': 'keep running'},
            ],
            'reasoning': 'Fish are active and sensor readings are normal.',
          },
        },
      },
    });

    expect(report.summary, '水质稳定');
    expect(
      report.contentForLanguageCode('en').summary,
      'Water quality is stable',
    );
    expect(report.contentForLanguageCode('en').actions.single.device, 'pump');
  });

  test('submitSensorReading refuses disconnected tank', () {
    final repository = AquaApiRepository();

    expect(
      () => repository.submitSensorReading(
        '',
        SensorReading(
          temperature: 26,
          ph: 7.2,
          oxygen: 6,
          ammonia: 0.1,
          waterLevel: 80,
          tds: 300,
          nitrate: 10,
          nitrite: 0.2,
          chloride: 30,
          timestamp: DateTime(2026, 7, 2),
        ),
      ),
      throwsA(
        isA<StateError>().having(
          (error) => error.message,
          'message',
          '未连接鱼缸，无法录入传感器数据。',
        ),
      ),
    );
  });

  test('SensorReading uses backend field names and Shanghai timestamp', () {
    final reading = SensorReading(
      temperature: 26,
      ph: 7.2,
      oxygen: 6,
      ammonia: 0.1,
      waterLevel: 80,
      tds: 300,
      nitrate: 10,
      nitrite: 0.2,
      chloride: 30,
      timestamp: DateTime.utc(2026, 7, 2, 12),
    );

    expect(reading.toSensorPayload(), {
      'temperature': 26,
      'ph': 7.2,
      'oxygen': 6,
      'ammonia': 0.1,
      'waterLevel': 80,
      'tds': 300,
      'nitrate': 10,
      'nitrite': 0.2,
      'chlorideIon': 30,
      'timestamp': '2026-07-02T20:00:00.000+08:00',
    });
  });
}
