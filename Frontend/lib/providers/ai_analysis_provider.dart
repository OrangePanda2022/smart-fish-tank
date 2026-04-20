import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/data/services/ai_api_service.dart';
import 'package:aqua/data/models/ai_api_response.dart';
import 'package:aqua/domain/models/analysis_result.dart';

enum AnalysisStatus { idle, loading, success, error }

class AnalysisState {
  final AnalysisStatus status;
  final AnalysisResult? result;
  final String? errorMessage;

  const AnalysisState({
    this.status = AnalysisStatus.idle,
    this.result,
    this.errorMessage,
  });

  AnalysisState copyWith({
    AnalysisStatus? status,
    AnalysisResult? result,
    String? errorMessage,
  }) {
    return AnalysisState(
      status: status ?? this.status,
      result: result ?? this.result,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }
}

final aiApiServiceProvider = AiApiService();

final aiAnalysisProvider = NotifierProvider<AiAnalysisNotifier, AnalysisState>(
  AiAnalysisNotifier.new,
);

class AiAnalysisNotifier extends Notifier<AnalysisState> {
  @override
  AnalysisState build() {
    return const AnalysisState();
  }

  AiApiService get _service => aiApiServiceProvider;

  Future<void> analyze(String tankId) async {
    state = state.copyWith(status: AnalysisStatus.loading);
    try {
      final apiResponse = await _service.analyse(tankId);
      final result = _convertToAnalysisResult(apiResponse);
      state = AnalysisState(status: AnalysisStatus.success, result: result);
    } on AiApiException catch (e) {
      state = AnalysisState(
        status: AnalysisStatus.error,
        errorMessage: e.message,
      );
    } catch (e) {
      state = AnalysisState(
        status: AnalysisStatus.error,
        errorMessage: e.toString(),
      );
    }
  }

  AnalysisResult _convertToAnalysisResult(AiApiResponse apiResponse) {
    final decision = apiResponse.data!.decision;
    return AnalysisResult(
      reportId: apiResponse.data!.reportId,
      statusScore: decision.statusScore,
      summary: decision.summary,
      suggestions: decision.actions
          .map(
            (action) =>
                AnalysisAction(device: action.device, action: action.action),
          )
          .toList(),
      reasoning: decision.reasoning,
      timestamp: DateTime.now(),
    );
  }

  void reset() {
    state = const AnalysisState();
  }
}
