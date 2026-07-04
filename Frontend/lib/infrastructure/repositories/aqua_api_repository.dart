import 'dart:async';
import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

import '../../domain/entities/analysis_report.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/entities/prediction_summary.dart';
import '../../domain/entities/sensor_reading.dart';
import '../../domain/entities/tank_info.dart';
import '../../domain/repositories/aquarium_repository.dart';

class AquaApiRepository implements AquariumRepository {
  AquaApiRepository({
    this.baseUrl = const String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: 'https://api-aqua.hubu.wiki',
    ),
  });

  final String baseUrl;

  static const _analysisScoreKeyPrefix = 'analysis_score_';
  static const _analysisSummaryKeyPrefix = 'analysis_summary_';
  static final Map<String, AquariumDashboard> _dashboardCache = {};
  static final Map<String, List<SensorReading>> _historyCache = {};
  static AquariumDashboard? _lastDashboard;

  @override
  Future<AquariumDashboard> loadDashboard() async {
    try {
      final tank = await _loadTank();
      final latest = await _loadLatest(tank.id);
      final analysis = await _loadCachedAnalysisSummary(tank.id);
      final prediction = await _loadPrediction(tank.id, latest);
      final cachedScore =
          await _loadCachedAnalysisScore(tank.id) ??
          _dashboardCache[tank.id]?.score;
      final dashboard = AquariumDashboard(
        tank: tank,
        latest: latest,
        prediction: prediction,
        score: cachedScore,
        analysisSummary: analysis ?? '',
      );
      _cacheDashboard(dashboard);
      return dashboard;
    } catch (_) {
      return _cachedDashboardOrEmpty();
    }
  }

  @override
  Future<List<SensorReading>> loadHistory(String tankId) async {
    try {
      if (tankId.isEmpty) return const [];
      final data = await _getJson(
        '/api/v1/tanks/$tankId/sensors/history?limit=100',
      );
      final raw = data['data'];
      if (raw is! List || raw.isEmpty) return _cachedHistoryOrEmpty(tankId);
      final history = raw
          .whereType<Map<String, dynamic>>()
          .map(SensorReading.fromJson)
          .where((reading) => reading.isUseful)
          .toList();
      if (history.isEmpty) return _cachedHistoryOrEmpty(tankId);
      _historyCache[tankId] = List.unmodifiable(history);
      return history;
    } catch (_) {
      return _cachedHistoryOrEmpty(tankId);
    }
  }

  @override
  Future<SensorReading> submitSensorReading(
    String tankId,
    SensorReading reading,
  ) async {
    if (tankId.isEmpty) {
      throw StateError('未连接鱼缸，无法录入传感器数据。');
    }

    final data = await _postJson(
      '/api/v1/tanks/$tankId/sensors',
      reading.toSensorPayload(),
    );
    final raw = data['data'];
    if (raw is Map<String, dynamic>) {
      final created = SensorReading.fromJson(raw);
      _historyCache.remove(tankId);
      return created;
    }
    throw const FormatException('Expected created sensor data');
  }

  List<SensorReading> _cachedHistoryOrEmpty(String tankId) {
    final cached = _historyCache[tankId];
    if (cached != null && cached.isNotEmpty) return cached;
    return const [];
  }

  @override
  Future<AnalysisReport> runAnalysis(String tankId) async {
    if (tankId.isEmpty) {
      throw StateError('未连接鱼缸，无法运行 AI 分析。');
    }
    final data = await _getJson(
      '/api/v1/analyse/$tankId',
      timeout: const Duration(seconds: 60),
    );
    final report = AnalysisReport.fromJson(data);
    await _cacheAnalysisScore(tankId, report.statusScore);
    await _cacheAnalysisSummary(tankId, _analysisSummaryFromReport(report));
    _cacheDashboardScore(tankId, report.statusScore);
    return report;
  }

  @override
  Uri videoStreamUri(String tankId) {
    if (tankId.isEmpty) return Uri();
    return Uri.parse('$baseUrl/api/v1/tank/$tankId/stream');
  }

  void _cacheDashboard(AquariumDashboard dashboard) {
    _dashboardCache[dashboard.tank.id] = dashboard;
    _lastDashboard = dashboard;
  }

  void _cacheDashboardScore(String tankId, int score) {
    final dashboard = _dashboardCache[tankId];
    if (dashboard == null) return;

    final updated = AquariumDashboard(
      tank: dashboard.tank,
      latest: dashboard.latest,
      prediction: dashboard.prediction,
      score: score,
      analysisSummary: dashboard.analysisSummary,
    );
    _cacheDashboard(updated);
  }

  AquariumDashboard _cachedDashboardOrEmpty() {
    final cached = _lastDashboard;
    if (cached != null) return cached;
    return AquariumDashboard.empty();
  }

  Future<TankInfo> _loadTank() async {
    final data = await _getJson('/api/v1/tank/');
    final raw = data['tanks'];
    if (raw is List && raw.isNotEmpty && raw.first is Map<String, dynamic>) {
      return TankInfo.fromJson(raw.first as Map<String, dynamic>);
    }
    return const TankInfo(id: '', name: '', size: 0, fishCount: 0);
  }

  Future<SensorReading> _loadLatest(String tankId) async {
    if (tankId.isEmpty) return AquariumDashboard.empty().latest;
    final data = await _getJson('/api/v1/tanks/$tankId/sensors/latest');
    final reading = SensorReading.fromJson(data);
    return reading.isUseful ? reading : AquariumDashboard.empty().latest;
  }

  Future<PredictionSummary> _loadPrediction(
    String tankId,
    SensorReading latest,
  ) async {
    if (tankId.isEmpty) return AquariumDashboard.empty().prediction;
    try {
      final data = await _getJson('/api/v1/predict/$tankId?horizon=60');
      return PredictionSummary.fromJson(data, fallback: latest);
    } catch (_) {
      return PredictionSummary(
        forecast: latest,
        confidence: 0,
        warning: '',
        generatedAt: DateTime.fromMillisecondsSinceEpoch(0),
        actions: const [],
      );
    }
  }

  Future<Map<String, dynamic>> _getJson(
    String path, {
    Duration timeout = const Duration(seconds: 4),
  }) async {
    final uri = Uri.parse('$baseUrl$path');
    final response = await http
        .get(uri, headers: const {'Accept': 'application/json'})
        .timeout(timeout);
    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw Exception('HTTP ${response.statusCode}: $uri');
    }
    final decoded = jsonDecode(utf8.decode(response.bodyBytes));
    if (decoded is Map<String, dynamic>) return decoded;
    throw const FormatException('Expected JSON object');
  }

  Future<Map<String, dynamic>> _postJson(
    String path,
    Map<String, dynamic> body, {
    Duration timeout = const Duration(seconds: 8),
  }) async {
    final uri = Uri.parse('$baseUrl$path');
    final response = await http
        .post(
          uri,
          headers: const {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
          },
          body: jsonEncode(body),
        )
        .timeout(timeout);
    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw Exception(
        'HTTP ${response.statusCode}: ${utf8.decode(response.bodyBytes)}',
      );
    }
    final decoded = jsonDecode(utf8.decode(response.bodyBytes));
    if (decoded is Map<String, dynamic>) return decoded;
    throw const FormatException('Expected JSON object');
  }

  Future<int?> _loadCachedAnalysisScore(String tankId) async {
    if (tankId.isEmpty) return null;
    final preferences = await SharedPreferences.getInstance();
    return preferences.getInt('$_analysisScoreKeyPrefix$tankId');
  }

  Future<String?> _loadCachedAnalysisSummary(String tankId) async {
    if (tankId.isEmpty) return null;
    final preferences = await SharedPreferences.getInstance();
    final summary = preferences.getString('$_analysisSummaryKeyPrefix$tankId');
    if (summary == null || summary.isEmpty) return null;
    return summary;
  }

  Future<void> _cacheAnalysisScore(String tankId, int score) async {
    if (tankId.isEmpty) return;
    final preferences = await SharedPreferences.getInstance();
    await preferences.setInt('$_analysisScoreKeyPrefix$tankId', score);
  }

  Future<void> _cacheAnalysisSummary(String tankId, String summary) async {
    if (tankId.isEmpty || summary.isEmpty) return;
    final preferences = await SharedPreferences.getInstance();
    await preferences.setString('$_analysisSummaryKeyPrefix$tankId', summary);
  }

  String _analysisSummaryFromReport(AnalysisReport report) {
    final content = report.zh;
    if (content.summary.isNotEmpty) return content.summary;
    return content.reasoning;
  }
}
