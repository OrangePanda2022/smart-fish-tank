import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/providers/ai_analysis_provider.dart';
import 'package:aqua/presentation/widgets/AnalysisResultCard.dart';

class AnalysisButton extends ConsumerWidget {
  const AnalysisButton({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final analysisState = ref.watch(aiAnalysisProvider);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '智能分析',
          style: TextStyle(
            fontWeight: FontWeight.bold,
            fontSize: 18,
            color: Colors.black87,
          ),
        ),
        const SizedBox(height: 12),
        _buildButton(context, ref, analysisState),
        if (analysisState.status == AnalysisStatus.success &&
            analysisState.result != null) ...[
          const SizedBox(height: 16),
          AnalysisResultCard(result: analysisState.result!),
        ],
        if (analysisState.status == AnalysisStatus.error) ...[
          const SizedBox(height: 16),
          _buildErrorCard(context, ref, analysisState.errorMessage),
        ],
      ],
    );
  }

  Widget _buildButton(
    BuildContext context,
    WidgetRef ref,
    AnalysisState state,
  ) {
    final isLoading = state.status == AnalysisStatus.loading;

    return GestureDetector(
      onTap: isLoading
          ? null
          : () => ref.read(aiAnalysisProvider.notifier).analyze('tank-001'),
      child: Container(
        width: double.infinity,
        height: 56,
        decoration: BoxDecoration(
          color: isLoading
              ? Theme.of(context).colorScheme.surfaceContainerHighest
              : Theme.of(context).colorScheme.primaryContainer,
          borderRadius: BorderRadius.circular(32),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            if (isLoading) ...[
              SizedBox(
                width: 24,
                height: 24,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  color: Theme.of(context).colorScheme.onPrimaryContainer,
                ),
              ),
              const SizedBox(width: 12),
              Text(
                '分析中...',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                  color: Theme.of(context).colorScheme.onPrimaryContainer,
                ),
              ),
            ] else ...[
              Icon(
                Icons.psychology,
                color: Theme.of(context).colorScheme.onPrimaryContainer,
              ),
              const SizedBox(width: 12),
              Text(
                '开始分析',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                  color: Theme.of(context).colorScheme.onPrimaryContainer,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildErrorCard(BuildContext context, WidgetRef ref, String? message) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.errorContainer,
        borderRadius: BorderRadius.circular(32),
      ),
      child: Row(
        children: [
          Icon(
            Icons.error_outline,
            color: Theme.of(context).colorScheme.onErrorContainer,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              message ?? '分析失败，请重试',
              style: TextStyle(
                color: Theme.of(context).colorScheme.onErrorContainer,
              ),
            ),
          ),
          IconButton(
            onPressed: () => ref.read(aiAnalysisProvider.notifier).reset(),
            icon: Icon(
              Icons.refresh,
              color: Theme.of(context).colorScheme.onErrorContainer,
            ),
          ),
        ],
      ),
    );
  }
}
