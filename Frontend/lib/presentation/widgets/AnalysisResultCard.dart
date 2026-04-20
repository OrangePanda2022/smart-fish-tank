import 'package:flutter/material.dart';
import 'package:aqua/domain/models/analysis_result.dart';
import 'package:aqua/presentation/widgets/TypewriterText.dart';

class AnalysisResultCard extends StatefulWidget {
  final AnalysisResult result;

  const AnalysisResultCard({super.key, required this.result});

  @override
  State<AnalysisResultCard> createState() => _AnalysisResultCardState();
}

class _AnalysisResultCardState extends State<AnalysisResultCard> {
  bool _showReport = false;
  bool _showSuggestions = false;
  bool _showReasoning = false;

  int get _reportDuration => widget.result.summary.length * 30;
  int get _suggestionsDuration =>
      _formatSuggestions(widget.result.suggestions).length * 30;

  @override
  void initState() {
    super.initState();
    _startAnimationSequence();
  }

  void _startAnimationSequence() {
    setState(() => _showReport = true);

    Future.delayed(Duration(milliseconds: _reportDuration + 300), () {
      if (mounted) setState(() => _showSuggestions = true);
    });

    Future.delayed(
      Duration(
        milliseconds: _reportDuration + 300 + _suggestionsDuration + 300,
      ),
      () {
        if (mounted) setState(() => _showReasoning = true);
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildScoreCard(context),
        const SizedBox(height: 12),
        _buildCard(
          context,
          icon: Icons.summarize,
          iconColor: Theme.of(context).colorScheme.primary,
          title: '报告',
          content: widget.result.summary,
          showAnimation: _showReport,
        ),
        const SizedBox(height: 12),
        _buildCard(
          context,
          icon: Icons.lightbulb,
          iconColor: Colors.amber,
          title: '建议',
          content: _formatSuggestions(widget.result.suggestions),
          showAnimation: _showSuggestions,
        ),
        const SizedBox(height: 12),
        _buildCard(
          context,
          icon: Icons.help_outline,
          iconColor: Theme.of(context).colorScheme.tertiary,
          title: '原因',
          content: widget.result.reasoning,
          showAnimation: _showReasoning,
        ),
      ],
    );
  }

  Widget _buildScoreCard(BuildContext context) {
    final scoreColor = widget.result.statusScore >= 80
        ? Colors.green
        : widget.result.statusScore >= 60
        ? Colors.orange
        : Colors.red;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.primaryContainer,
        borderRadius: BorderRadius.circular(32),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: scoreColor.withAlpha(51),
              borderRadius: BorderRadius.circular(24),
            ),
            child: Icon(Icons.analytics, color: scoreColor, size: 28),
          ),
          const SizedBox(width: 16),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '综合评分',
                style: TextStyle(
                  fontSize: 14,
                  color: Theme.of(context).colorScheme.onPrimaryContainer,
                ),
              ),
              Row(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '${widget.result.statusScore}',
                    style: TextStyle(
                      fontSize: 36,
                      fontWeight: FontWeight.bold,
                      color: scoreColor,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.only(bottom: 6),
                    child: Text(
                      ' / 100',
                      style: TextStyle(
                        fontSize: 16,
                        color: Theme.of(context).colorScheme.onPrimaryContainer,
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _formatSuggestions(List<AnalysisAction> suggestions) {
    if (suggestions.isEmpty) return '继续保持当前养护方案';
    return suggestions
        .map((s) => '${s.deviceDisplayName}: ${s.actionDisplayName}')
        .join('\n');
  }

  Widget _buildCard(
    BuildContext context, {
    required IconData icon,
    required Color iconColor,
    required String title,
    required String content,
    required bool showAnimation,
  }) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.secondaryContainer,
        borderRadius: BorderRadius.circular(32),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(10),
                decoration: BoxDecoration(
                  color: iconColor.withAlpha(51),
                  borderRadius: BorderRadius.circular(24),
                ),
                child: Icon(icon, color: iconColor, size: 20),
              ),
              const SizedBox(width: 12),
              Text(
                title,
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: Theme.of(context).colorScheme.onSecondaryContainer,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          if (showAnimation)
            TypewriterText(
              key: ValueKey('typewriter_$title'),
              text: content,
              style: TextStyle(
                fontSize: 14,
                height: 1.6,
                color: Theme.of(context).colorScheme.onSecondaryContainer,
              ),
              durationPerChar: 30,
            )
          else
            Text(
              '',
              style: TextStyle(
                fontSize: 14,
                height: 1.6,
                color: Theme.of(context).colorScheme.onSecondaryContainer,
              ),
            ),
        ],
      ),
    );
  }
}
