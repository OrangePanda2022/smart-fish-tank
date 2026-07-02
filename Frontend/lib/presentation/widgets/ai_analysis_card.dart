import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_shadow.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/analysis_report.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/repositories/aquarium_repository.dart';
import 'analysis_report_printer.dart';
import 'gradient_button.dart';
import 'soft_card.dart';

class AiAnalysisCard extends StatefulWidget {
  const AiAnalysisCard({
    required this.dashboard,
    required this.repository,
    required this.onAnalyze,
    required this.onAnalysisComplete,
    super.key,
  });

  final AquariumDashboard dashboard;
  final AquariumRepository repository;
  final Future<AnalysisReport> Function(String tankId) onAnalyze;
  final VoidCallback onAnalysisComplete;

  @override
  State<AiAnalysisCard> createState() => _AiAnalysisCardState();
}

class _AiAnalysisCardState extends State<AiAnalysisCard> {
  bool _isAnalyzing = false;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      padding: const EdgeInsets.fromLTRB(22, 22, 18, 18),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('AI 智能分析', style: AppText.titleSmall),
                const SizedBox(height: 8),
                const Text('上次分析：今天 09:30', style: AppText.caption),
                const SizedBox(height: 20),
                Row(
                  children: const [
                    HealthPill(icon: CupertinoIcons.drop_fill, text: '水质正常'),
                    SizedBox(width: 10),
                    HealthPill(icon: CupertinoIcons.thermometer, text: '温度适宜'),
                    SizedBox(width: 10),
                    HealthPill(icon: Icons.set_meal_rounded, text: '鱼儿活跃'),
                  ],
                ),
                const SizedBox(height: 18),
                GradientButton(
                  text: _isAnalyzing ? '分析中...' : '立即分析',
                  icon: CupertinoIcons.sparkles,
                  isBusy: _isAnalyzing,
                  onTap: _analyzeNow,
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          SizedBox(
            width: 116,
            height: 132,
            child: Image.asset(
              'assets/images/ai_robot.png',
              fit: BoxFit.contain,
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _analyzeNow() async {
    setState(() {
      _isAnalyzing = true;
    });

    try {
      final analysis = await widget.onAnalyze(widget.dashboard.tank.id);
      if (!mounted) return;
      _showAnalysis(context, analysis);
      widget.onAnalysisComplete();
    } catch (_) {
      if (!mounted) return;
      _showAnalysisError(context);
    } finally {
      if (mounted) {
        setState(() {
          _isAnalyzing = false;
        });
      }
    }
  }

  void _showAnalysis(BuildContext context, AnalysisReport report) {
    showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (_) {
        return DraggableScrollableSheet(
          initialChildSize: 0.78,
          minChildSize: 0.45,
          maxChildSize: 0.92,
          expand: false,
          builder: (context, scrollController) {
            return Container(
              padding: const EdgeInsets.fromLTRB(20, 12, 20, 26),
              decoration: const BoxDecoration(
                color: AppColors.background,
                borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
              ),
              child: ListView(
                controller: scrollController,
                children: [
                  Center(
                    child: Container(
                      width: 42,
                      height: 5,
                      decoration: BoxDecoration(
                        color: AppColors.grid,
                        borderRadius: BorderRadius.circular(999),
                      ),
                    ),
                  ),
                  const SizedBox(height: 18),
                  _AnalysisHero(report: report),
                  const SizedBox(height: 14),
                  _InfoGrid(report: report),
                  const SizedBox(height: 14),
                  _ReportSection(
                    icon: CupertinoIcons.doc_text_fill,
                    title: '状态概要',
                    child: Text(report.summary, style: AppText.body),
                  ),
                  const SizedBox(height: 14),
                  _ReportSection(
                    icon: CupertinoIcons.slider_horizontal_3,
                    title: '建议操作',
                    child: report.actions.isEmpty
                        ? const Text('暂无建议操作', style: AppText.caption)
                        : Column(
                            children: [
                              for (final action in report.actions)
                                _ActionRow(action: action),
                            ],
                          ),
                  ),
                  const SizedBox(height: 14),
                  _ReportSection(
                    icon: CupertinoIcons.lightbulb_fill,
                    title: '推理过程',
                    child: Text(report.reasoning, style: AppText.body),
                  ),
                  const SizedBox(height: 18),
                  _PrintReportButton(
                    report: report,
                    dashboard: widget.dashboard,
                    repository: widget.repository,
                  ),
                ],
              ),
            );
          },
        );
      },
    );
  }

  void _showAnalysisError(BuildContext context) {
    showModalBottomSheet<void>(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (_) {
        return const Padding(
          padding: EdgeInsets.fromLTRB(24, 24, 24, 34),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('AI 分析失败', style: AppText.titleSmall),
              SizedBox(height: 12),
              Text(
                '分析请求失败或超时，请确认后端服务可用后重试。',
                style: TextStyle(
                  color: AppColors.ink,
                  height: 1.6,
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _PrintReportButton extends StatefulWidget {
  const _PrintReportButton({
    required this.report,
    required this.dashboard,
    required this.repository,
  });

  final AnalysisReport report;
  final AquariumDashboard dashboard;
  final AquariumRepository repository;

  @override
  State<_PrintReportButton> createState() => _PrintReportButtonState();
}

class _PrintReportButtonState extends State<_PrintReportButton> {
  bool _isPrinting = false;

  @override
  Widget build(BuildContext context) {
    return GradientButton(
      text: _isPrinting ? '生成中...' : '打印报告',
      icon: Icons.print_rounded,
      isBusy: _isPrinting,
      onTap: _printReport,
    );
  }

  Future<void> _printReport() async {
    final language = await _chooseReportLanguage();
    if (language == null) return;

    setState(() {
      _isPrinting = true;
    });

    try {
      final history = await widget.repository.loadHistory(
        widget.dashboard.tank.id,
      );
      await AnalysisReportPrinter.printReport(
        report: widget.report,
        dashboard: widget.dashboard,
        history: history,
        language: language,
      );
    } catch (error) {
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('打印失败: $error')));
    } finally {
      if (mounted) {
        setState(() {
          _isPrinting = false;
        });
      }
    }
  }

  Future<ReportLanguage?> _chooseReportLanguage() {
    return showModalBottomSheet<ReportLanguage>(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (context) {
        return SafeArea(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(20, 16, 20, 18),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('选择打印语言', style: AppText.titleSmall),
                const SizedBox(height: 14),
                _PrintLanguageTile(
                  icon: CupertinoIcons.textformat,
                  title: '打印中文版',
                  subtitle: '使用中文标题、标签与页脚生成报告',
                  onTap: () =>
                      Navigator.of(context).pop(ReportLanguage.chinese),
                ),
                const SizedBox(height: 10),
                _PrintLanguageTile(
                  icon: CupertinoIcons.globe,
                  title: '打印英文版',
                  subtitle: 'Use English labels and report headings',
                  onTap: () =>
                      Navigator.of(context).pop(ReportLanguage.english),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}

class _PrintLanguageTile extends StatelessWidget {
  const _PrintLanguageTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: AppColors.paleBlue,
      borderRadius: BorderRadius.circular(18),
      child: InkWell(
        borderRadius: BorderRadius.circular(18),
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(14),
          child: Row(
            children: [
              Container(
                width: 42,
                height: 42,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(14),
                ),
                child: Icon(icon, color: AppColors.primary, size: 22),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(title, style: AppText.cardTitle),
                    const SizedBox(height: 4),
                    Text(subtitle, style: AppText.caption),
                  ],
                ),
              ),
              const Icon(
                CupertinoIcons.chevron_forward,
                color: AppColors.muted,
                size: 18,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _AnalysisHero extends StatelessWidget {
  const _AnalysisHero({required this.report});

  final AnalysisReport report;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      radius: 22,
      padding: const EdgeInsets.all(18),
      child: Row(
        children: [
          Container(
            width: 72,
            height: 72,
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [AppColors.primary, AppColors.teal],
              ),
              borderRadius: BorderRadius.circular(22),
              boxShadow: AppShadow.button,
            ),
            child: Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    '${report.statusScore}',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 26,
                      height: 1,
                      fontWeight: FontWeight.w900,
                    ),
                  ),
                  const SizedBox(height: 4),
                ],
              ),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('AI 分析结果', style: AppText.titleSmall),
                const SizedBox(height: 6),
                Text(
                  report.summary.isEmpty ? '暂无概要' : report.summary,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: AppText.caption,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _InfoGrid extends StatelessWidget {
  const _InfoGrid({required this.report});

  final AnalysisReport report;

  @override
  Widget build(BuildContext context) {
    return _InfoTile(label: 'report_id', value: report.reportId);
  }
}

class _InfoTile extends StatelessWidget {
  const _InfoTile({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      radius: 18,
      padding: const EdgeInsets.all(14),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(label, style: AppText.caption),
          const SizedBox(height: 8),
          Text(
            value.isEmpty ? '-' : value,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(
              color: AppColors.ink,
              fontSize: 15,
              fontWeight: FontWeight.w900,
            ),
          ),
        ],
      ),
    );
  }
}

class _ReportSection extends StatelessWidget {
  const _ReportSection({
    required this.icon,
    required this.title,
    required this.child,
  });

  final IconData icon;
  final String title;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      radius: 22,
      padding: const EdgeInsets.all(18),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                width: 34,
                height: 34,
                decoration: BoxDecoration(
                  color: AppColors.paleBlue,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(icon, color: AppColors.primary, size: 19),
              ),
              const SizedBox(width: 10),
              Text(title, style: AppText.cardTitle),
            ],
          ),
          const SizedBox(height: 14),
          child,
        ],
      ),
    );
  }
}

class _ActionRow extends StatelessWidget {
  const _ActionRow({required this.action});

  final AnalysisAction action;

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.paleBlue,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(
            CupertinoIcons.check_mark_circled_solid,
            color: AppColors.success,
            size: 20,
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  action.device.isEmpty ? 'unknown' : action.device,
                  style: AppText.pill,
                ),
                const SizedBox(height: 4),
                Text(action.action, style: AppText.caption),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class HealthPill extends StatelessWidget {
  const HealthPill({required this.icon, required this.text, super.key});

  final IconData icon;
  final String text;

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(16),
          boxShadow: AppShadow.soft,
        ),
        child: Column(
          children: [
            Icon(icon, color: AppColors.primary, size: 22),
            const SizedBox(height: 5),
            FittedBox(
              fit: BoxFit.scaleDown,
              child: Text(text, style: AppText.pill),
            ),
            const SizedBox(height: 4),
            const Icon(
              CupertinoIcons.check_mark_circled_solid,
              color: AppColors.success,
              size: 18,
            ),
          ],
        ),
      ),
    );
  }
}
