import 'dart:math' as math;

import 'package:flutter/services.dart' show rootBundle;
import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart' as pw;
import 'package:printing/printing.dart';

import '../../domain/entities/analysis_report.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/entities/sensor_reading.dart';
import '../view_models/chart_series.dart';

enum ReportLanguage { chinese, english }

class AnalysisReportPrinter {
  const AnalysisReportPrinter._();

  static Future<void> printReport({
    required AnalysisReport report,
    required AquariumDashboard dashboard,
    required List<SensorReading> history,
    ReportLanguage language = ReportLanguage.chinese,
  }) async {
    final doc = await _buildDocument(
      report: report,
      dashboard: dashboard,
      history: history,
      language: language,
    );
    final bytes = await doc.save();
    final suffix = language == ReportLanguage.english ? 'EN' : 'ZH';
    await Printing.layoutPdf(
      name: 'AquaClaw-${report.reportId}-$suffix.pdf',
      onLayout: (_) async => bytes,
    );
  }

  static Future<pw.Document> _buildDocument({
    required AnalysisReport report,
    required AquariumDashboard dashboard,
    required List<SensorReading> history,
    required ReportLanguage language,
  }) async {
    final fontData = await rootBundle.load(
      'assets/fonts/NotoSerifSC-Medium.ttf',
    );
    final cjkFont = pw.Font.ttf(fontData);
    final theme = pw.ThemeData.withFont(
      base: cjkFont,
      bold: cjkFont,
      italic: cjkFont,
      boldItalic: cjkFont,
    );
    final doc = pw.Document(theme: theme);
    final now = DateTime.now();
    final sortedHistory = [...history]
      ..sort((a, b) => a.timestamp.compareTo(b.timestamp));
    final chartHistory = sortedHistory.length > 30
        ? sortedHistory.sublist(sortedHistory.length - 30)
        : sortedHistory;
    final latest = dashboard.latest;
    final forecast = dashboard.prediction.forecast;
    final scoreColor = _scoreColor(report.statusScore);
    final labels = _ReportLabels.forLanguage(language);

    doc.addPage(
      pw.Page(
        pageFormat: PdfPageFormat.a4,
        margin: const pw.EdgeInsets.fromLTRB(22, 20, 22, 18),
        build: (context) {
          return pw.Column(
            crossAxisAlignment: pw.CrossAxisAlignment.start,
            children: [
              pw.Row(
                crossAxisAlignment: pw.CrossAxisAlignment.start,
                children: [
                  pw.Expanded(
                    child: pw.Column(
                      crossAxisAlignment: pw.CrossAxisAlignment.start,
                      children: [
                        pw.Text(
                          labels.title,
                          style: pw.TextStyle(
                            fontSize: 16,
                            fontWeight: pw.FontWeight.bold,
                            color: _pdfColor(0xFF071631),
                          ),
                        ),
                        pw.SizedBox(height: 4),
                        pw.Text(
                          '${labels.reportId}: ${_emptyDash(report.reportId)}  |  ${labels.testTime}: ${_formatDateTime(now)}  |  ${labels.tank}: ${dashboard.tank.name}',
                          style: pw.TextStyle(
                            fontSize: 8,
                            color: _pdfColor(0xFF66738A),
                          ),
                        ),
                      ],
                    ),
                  ),
                  pw.Container(
                    width: 88,
                    padding: const pw.EdgeInsets.symmetric(vertical: 8),
                    decoration: pw.BoxDecoration(
                      color: _scoreBackgroundColor(report.statusScore),
                      borderRadius: pw.BorderRadius.circular(12),
                      border: pw.Border.all(color: scoreColor, width: 0.8),
                    ),
                    child: pw.Column(
                      children: [
                        pw.Text(
                          '${report.statusScore}',
                          style: pw.TextStyle(
                            fontSize: 30,
                            fontWeight: pw.FontWeight.bold,
                            color: scoreColor,
                          ),
                        ),
                        pw.Text(
                          '${labels.health}: ${_scoreLabel(report.statusScore, labels)}',
                          style: pw.TextStyle(
                            fontSize: 7.5,
                            fontWeight: pw.FontWeight.bold,
                            color: scoreColor,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              pw.SizedBox(height: 8),
              _summaryBlock(report.summary, labels),
              _sectionTitle(labels.visualAnalysis),
              pw.SizedBox(height: 7),
              _chartFrame(
                title: labels.trendTitle,
                minHeight: 215,
                child: pw.Column(
                  crossAxisAlignment: pw.CrossAxisAlignment.start,
                  children: [
                    _legend([
                      _LegendItem(labels.ph, _pdfColor(0xFF2F7AF6)),
                      _LegendItem(labels.temperature, _pdfColor(0xFFFF822E)),
                    ]),
                    pw.SizedBox(height: 4),
                    pw.SvgImage(
                      svg: _trendChartSvg(chartHistory),
                      width: double.infinity,
                      height: 160,
                    ),
                  ],
                ),
              ),
              pw.SizedBox(height: 8),
              pw.Row(
                crossAxisAlignment: pw.CrossAxisAlignment.start,
                children: [
                  pw.Expanded(
                    child: _chartFrame(
                      title: labels.nitrogenTitle,
                      minHeight: 195,
                      child: pw.SvgImage(
                        svg: _nitrogenBarChartSvg(latest),
                        width: double.infinity,
                        height: 150,
                      ),
                    ),
                  ),
                  pw.SizedBox(width: 12),
                  pw.Expanded(
                    child: _chartFrame(
                      title: labels.radarTitle,
                      minHeight: 195,
                      child: pw.Column(
                        crossAxisAlignment: pw.CrossAxisAlignment.start,
                        children: [
                          _legend([
                            _LegendItem(labels.current, _pdfColor(0xFF2F7AF6)),
                            _LegendItem(labels.forecast, _pdfColor(0xFF45C77A)),
                          ]),
                          pw.SizedBox(height: 2),
                          pw.SvgImage(
                            svg: _radarChartSvg(latest, forecast),
                            width: double.infinity,
                            height: 128,
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
              pw.SizedBox(height: 8),
              pw.Row(
                crossAxisAlignment: pw.CrossAxisAlignment.start,
                children: [
                  pw.Expanded(
                    child: _textPanel(
                      title: labels.reasoningTitle,
                      color: _pdfColor(0xFF2F7AF6),
                      minHeight: 116,
                      children: [
                        for (final item in _splitReasoning(
                          report.reasoning,
                          labels.noReasoning,
                        ).take(3))
                          _bullet(item),
                      ],
                    ),
                  ),
                  pw.SizedBox(width: 12),
                  pw.Expanded(
                    child: _textPanel(
                      title: labels.actionTitle,
                      color: _pdfColor(0xFF45C77A),
                      background: _pdfColor(0xFFF0FDF4),
                      minHeight: 116,
                      children: report.actions.isEmpty
                          ? [_bullet(labels.noActions)]
                          : [
                              for (final entry
                                  in report.actions
                                      .take(3)
                                      .toList()
                                      .asMap()
                                      .entries)
                                _actionLine(entry.key + 1, entry.value),
                            ],
                    ),
                  ),
                ],
              ),
              pw.Spacer(),
              pw.Center(
                child: pw.Text(
                  '${labels.footer} - ${_formatDateTime(now)}',
                  style: pw.TextStyle(
                    fontSize: 7.5,
                    color: _pdfColor(0xFF8D99AA),
                  ),
                ),
              ),
            ],
          );
        },
      ),
    );

    return doc;
  }

  static pw.Widget _summaryBlock(String summary, _ReportLabels labels) {
    return pw.Container(
      margin: const pw.EdgeInsets.only(bottom: 8),
      padding: const pw.EdgeInsets.all(10),
      decoration: pw.BoxDecoration(
        color: _pdfColor(0xFFF7FBFF),
        border: pw.Border(
          left: pw.BorderSide(color: _pdfColor(0xFFFF822E), width: 4),
        ),
      ),
      child: pw.Column(
        crossAxisAlignment: pw.CrossAxisAlignment.start,
        children: [
          pw.Text(
            labels.summaryTitle,
            style: pw.TextStyle(
              fontSize: 11,
              fontWeight: pw.FontWeight.bold,
              color: _pdfColor(0xFF071631),
            ),
          ),
          pw.SizedBox(height: 4),
          pw.Text(
            summary.isEmpty ? labels.noSummary : summary,
            maxLines: 3,
            overflow: pw.TextOverflow.clip,
            style: pw.TextStyle(
              fontSize: 8.5,
              lineSpacing: 1.2,
              color: _pdfColor(0xFF334155),
            ),
          ),
        ],
      ),
    );
  }

  static pw.Widget _sectionTitle(String text) {
    return pw.Text(
      text,
      style: pw.TextStyle(
        fontSize: 12,
        fontWeight: pw.FontWeight.bold,
        color: _pdfColor(0xFF071631),
      ),
    );
  }

  static pw.Widget _chartFrame({
    required String title,
    required pw.Widget child,
    double? minHeight,
  }) {
    return pw.Container(
      constraints: minHeight == null
          ? null
          : pw.BoxConstraints(minHeight: minHeight),
      padding: const pw.EdgeInsets.all(8),
      decoration: pw.BoxDecoration(
        border: pw.Border.all(color: _pdfColor(0xFFE1EAF6), width: 0.6),
        borderRadius: pw.BorderRadius.circular(8),
      ),
      child: pw.Column(
        crossAxisAlignment: pw.CrossAxisAlignment.start,
        children: [
          pw.Text(
            title,
            style: pw.TextStyle(
              fontSize: 8.5,
              fontWeight: pw.FontWeight.bold,
              color: _pdfColor(0xFF66738A),
            ),
          ),
          pw.SizedBox(height: 5),
          child,
        ],
      ),
    );
  }

  static pw.Widget _textPanel({
    required String title,
    required PdfColor color,
    required List<pw.Widget> children,
    PdfColor? background,
    double? minHeight,
  }) {
    return pw.Container(
      constraints: minHeight == null
          ? null
          : pw.BoxConstraints(minHeight: minHeight),
      padding: const pw.EdgeInsets.all(8),
      decoration: pw.BoxDecoration(
        color: background,
        border: pw.Border.all(color: _pdfColor(0xFFE1EAF6), width: 0.6),
        borderRadius: pw.BorderRadius.circular(8),
      ),
      child: pw.Column(
        crossAxisAlignment: pw.CrossAxisAlignment.start,
        children: [
          pw.Text(
            title,
            style: pw.TextStyle(
              fontSize: 9.5,
              fontWeight: pw.FontWeight.bold,
              color: color,
            ),
          ),
          pw.SizedBox(height: 5),
          ...children,
        ],
      ),
    );
  }

  static pw.Widget _bullet(String text) {
    return pw.Padding(
      padding: const pw.EdgeInsets.only(bottom: 4),
      child: pw.Row(
        crossAxisAlignment: pw.CrossAxisAlignment.start,
        children: [
          pw.Text('- ', style: pw.TextStyle(fontSize: 8)),
          pw.Expanded(
            child: pw.Text(
              text,
              maxLines: 2,
              overflow: pw.TextOverflow.clip,
              style: pw.TextStyle(
                fontSize: 7.5,
                lineSpacing: 1.15,
                color: _pdfColor(0xFF334155),
              ),
            ),
          ),
        ],
      ),
    );
  }

  static pw.Widget _actionLine(int index, AnalysisAction action) {
    return pw.Padding(
      padding: const pw.EdgeInsets.only(bottom: 4),
      child: pw.Row(
        crossAxisAlignment: pw.CrossAxisAlignment.start,
        children: [
          pw.Container(
            width: 16,
            height: 12,
            decoration: pw.BoxDecoration(
              color: _pdfColor(0xFF45C77A),
              borderRadius: pw.BorderRadius.circular(3),
            ),
            child: pw.Center(
              child: pw.Text(
                index.toString().padLeft(2, '0'),
                style: pw.TextStyle(
                  fontSize: 6,
                  fontWeight: pw.FontWeight.bold,
                  color: PdfColors.white,
                ),
              ),
            ),
          ),
          pw.SizedBox(width: 6),
          pw.Expanded(
            child: pw.Text(
              '${_emptyDash(action.device)}: ${_emptyDash(action.action)}',
              maxLines: 2,
              overflow: pw.TextOverflow.clip,
              style: pw.TextStyle(fontSize: 7.5, lineSpacing: 1.15),
            ),
          ),
        ],
      ),
    );
  }

  static pw.Widget _legend(List<_LegendItem> items) {
    return pw.Row(
      children: [
        for (final item in items) ...[
          pw.Container(width: 8, height: 8, color: item.color),
          pw.SizedBox(width: 4),
          pw.Text(
            item.label,
            style: pw.TextStyle(fontSize: 7, color: _pdfColor(0xFF66738A)),
          ),
          pw.SizedBox(width: 10),
        ],
      ],
    );
  }

  static String _trendChartSvg(List<SensorReading> history) {
    const width = 720.0;
    const height = 240.0;
    const left = 44.0;
    const right = 18.0;
    const top = 18.0;
    const bottom = 30.0;
    final chartWidth = width - left - right;
    final chartHeight = height - top - bottom;

    if (history.isEmpty) {
      return _emptyChartSvg(width, height, 'No data');
    }

    final phMetric = ChartSeries.metricFor('ph');
    final tempMetric = ChartSeries.metricFor('temperature');
    final phValues = history.map((reading) => reading.ph).toList();
    final tempValues = history.map((reading) => reading.temperature).toList();
    final phRange = _Range.fromValues(
      phValues,
      phMetric.safeLow,
      phMetric.safeHigh,
    );
    final tempRange = _Range.fromValues(
      tempValues,
      tempMetric.safeLow,
      tempMetric.safeHigh,
    );

    String polyline(List<double> values, _Range range) {
      if (values.length == 1) {
        final y = top + chartHeight * (1 - range.normalize(values.first));
        return '${left.toStringAsFixed(1)},${y.toStringAsFixed(1)} '
            '${(left + chartWidth).toStringAsFixed(1)},${y.toStringAsFixed(1)}';
      }

      return values
          .asMap()
          .entries
          .map((entry) {
            final x = left + chartWidth * entry.key / (values.length - 1);
            final y = top + chartHeight * (1 - range.normalize(entry.value));
            return '${x.toStringAsFixed(1)},${y.toStringAsFixed(1)}';
          })
          .join(' ');
    }

    final grid = List.generate(5, (index) {
      final y = top + chartHeight * index / 4;
      return '<line x1="$left" y1="$y" x2="${width - right}" y2="$y" '
          'stroke="#E1EAF6" stroke-width="1"/>';
    }).join();

    return '''
<svg xmlns="http://www.w3.org/2000/svg" width="$width" height="$height" viewBox="0 0 $width $height">
  <rect width="$width" height="$height" fill="#FFFFFF"/>
  $grid
  <line x1="$left" y1="${height - bottom}" x2="${width - right}" y2="${height - bottom}" stroke="#CBD5E1" stroke-width="1.2"/>
  <line x1="$left" y1="$top" x2="$left" y2="${height - bottom}" stroke="#CBD5E1" stroke-width="1.2"/>
  <text x="6" y="$top" fill="#64748B" font-size="18">${phRange.high.toStringAsFixed(1)}</text>
  <text x="6" y="${height - bottom}" fill="#64748B" font-size="18">${phRange.low.toStringAsFixed(1)}</text>
  <text x="${width - 76}" y="$top" fill="#64748B" font-size="18">${tempRange.high.toStringAsFixed(1)}C</text>
  <text x="${width - 76}" y="${height - bottom}" fill="#64748B" font-size="18">${tempRange.low.toStringAsFixed(1)}C</text>
  <polyline points="${polyline(phValues, phRange)}" fill="none" stroke="#2F7AF6" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
  <polyline points="${polyline(tempValues, tempRange)}" fill="none" stroke="#FF822E" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
</svg>
''';
  }

  static String _nitrogenBarChartSvg(SensorReading latest) {
    const width = 340.0;
    const height = 240.0;
    const top = 24.0;
    const bottom = 44.0;
    const chartHeight = height - top - bottom;
    final items = [
      _BarItem('NH3', latest.ammonia, 0.08, '#FF525D'),
      _BarItem('NO2', latest.nitrite, 0.10, '#FFCF3A'),
      _BarItem('NO3', latest.nitrate, 20.0, '#39CE89'),
    ];
    final bars = items.asMap().entries.map((entry) {
      final item = entry.value;
      final barWidth = 58.0;
      final gap = 38.0;
      final x = 46 + entry.key * (barWidth + gap);
      final ratio = (item.value / item.safeHigh).clamp(0.0, 1.25);
      final barHeight = math.min(chartHeight, chartHeight * ratio);
      final y = top + chartHeight - barHeight;
      return '''
  <rect x="$x" y="$y" width="$barWidth" height="$barHeight" rx="9" fill="${item.color}"/>
  <line x1="$x" y1="$top" x2="${x + barWidth}" y2="$top" stroke="#94A3B8" stroke-width="1" stroke-dasharray="5 5"/>
  <text x="${x + barWidth / 2}" y="${y - 8}" text-anchor="middle" fill="#334155" font-size="18">${item.value.toStringAsFixed(item.value < 1 ? 2 : 1)}</text>
  <text x="${x + barWidth / 2}" y="${height - 14}" text-anchor="middle" fill="#64748B" font-size="18">${item.label}</text>
''';
    }).join();

    return '''
<svg xmlns="http://www.w3.org/2000/svg" width="$width" height="$height" viewBox="0 0 $width $height">
  <rect width="$width" height="$height" fill="#FFFFFF"/>
  <line x1="24" y1="${height - bottom}" x2="${width - 22}" y2="${height - bottom}" stroke="#CBD5E1" stroke-width="1.2"/>
  $bars
</svg>
''';
  }

  static String _radarChartSvg(SensorReading latest, SensorReading forecast) {
    const width = 340.0;
    const height = 210.0;
    const centerX = width / 2;
    const centerY = 100.0;
    const radius = 74.0;
    final metrics = [
      _RadarMetric('T', latest.temperature, forecast.temperature, 23.5, 29),
      _RadarMetric('pH', latest.ph, forecast.ph, 6.5, 7.8),
      _RadarMetric('O2', latest.oxygen, forecast.oxygen, 6, 10),
      _RadarMetric('NH3', latest.ammonia, forecast.ammonia, 0, 0.08),
      _RadarMetric('WL', latest.waterLevel, forecast.waterLevel, 60, 100),
    ];

    List<math.Point<double>> points(double Function(_RadarMetric) valueOf) {
      return metrics.asMap().entries.map((entry) {
        final angle = -math.pi / 2 + math.pi * 2 * entry.key / metrics.length;
        final metric = entry.value;
        final ratio = metric.normalize(valueOf(metric));
        return math.Point(
          centerX + math.cos(angle) * radius * ratio,
          centerY + math.sin(angle) * radius * ratio,
        );
      }).toList();
    }

    String pointString(List<math.Point<double>> points) {
      return points
          .map(
            (point) =>
                '${point.x.toStringAsFixed(1)},${point.y.toStringAsFixed(1)}',
          )
          .join(' ');
    }

    final grid = List.generate(4, (index) {
      final r = radius * (index + 1) / 4;
      final ring = List.generate(metrics.length, (metricIndex) {
        final angle = -math.pi / 2 + math.pi * 2 * metricIndex / metrics.length;
        return math.Point(
          centerX + math.cos(angle) * r,
          centerY + math.sin(angle) * r,
        );
      });
      return '<polygon points="${pointString(ring)}" fill="none" stroke="#E1EAF6" stroke-width="1"/>';
    }).join();

    final axes = metrics.asMap().entries.map((entry) {
      final angle = -math.pi / 2 + math.pi * 2 * entry.key / metrics.length;
      final endX = centerX + math.cos(angle) * radius;
      final endY = centerY + math.sin(angle) * radius;
      final labelX = centerX + math.cos(angle) * (radius + 20);
      final labelY = centerY + math.sin(angle) * (radius + 20);
      return '''
  <line x1="$centerX" y1="$centerY" x2="$endX" y2="$endY" stroke="#E1EAF6" stroke-width="1"/>
  <text x="$labelX" y="$labelY" text-anchor="middle" dominant-baseline="middle" fill="#64748B" font-size="16">${entry.value.label}</text>
''';
    }).join();

    return '''
<svg xmlns="http://www.w3.org/2000/svg" width="$width" height="$height" viewBox="0 0 $width $height">
  <rect width="$width" height="$height" fill="#FFFFFF"/>
  $grid
  $axes
  <polygon points="${pointString(points((m) => m.current))}" fill="#2F7AF6" fill-opacity="0.18" stroke="#2F7AF6" stroke-width="3"/>
  <polygon points="${pointString(points((m) => m.forecast))}" fill="#45C77A" fill-opacity="0.16" stroke="#45C77A" stroke-width="3"/>
</svg>
''';
  }

  static String _emptyChartSvg(double width, double height, String label) {
    return '''
<svg xmlns="http://www.w3.org/2000/svg" width="$width" height="$height" viewBox="0 0 $width $height">
  <rect width="$width" height="$height" fill="#FFFFFF"/>
  <text x="${width / 2}" y="${height / 2}" text-anchor="middle" fill="#94A3B8" font-size="18">$label</text>
</svg>
''';
  }

  static List<String> _splitReasoning(String text, String fallback) {
    final parts = text
        .replaceAll('；', ';')
        .split(RegExp(r'[;。\n]'))
        .map((value) => value.trim())
        .where((value) => value.isNotEmpty)
        .toList();
    if (parts.isEmpty) return [fallback];
    return parts;
  }

  static PdfColor _scoreColor(int score) {
    if (score >= 80) return _pdfColor(0xFF45C77A);
    if (score >= 60) return _pdfColor(0xFFFF822E);
    return _pdfColor(0xFFFF525D);
  }

  static PdfColor _scoreBackgroundColor(int score) {
    if (score >= 80) return _pdfColor(0xFFEAFBF1);
    if (score >= 60) return _pdfColor(0xFFFFF3E8);
    return _pdfColor(0xFFFFEEF0);
  }

  static String _scoreLabel(int score, _ReportLabels labels) {
    if (score >= 80) return labels.scoreGood;
    if (score >= 60) return labels.scoreWarning;
    return labels.scoreAlert;
  }

  static PdfColor _pdfColor(int value) => PdfColor.fromInt(value);

  static String _formatDateTime(DateTime value) {
    return '${value.year}-${value.month.toString().padLeft(2, '0')}-'
        '${value.day.toString().padLeft(2, '0')} '
        '${value.hour.toString().padLeft(2, '0')}:'
        '${value.minute.toString().padLeft(2, '0')}';
  }

  static String _emptyDash(String value) => value.isEmpty ? '-' : value;
}

class _ReportLabels {
  const _ReportLabels({
    required this.title,
    required this.reportId,
    required this.testTime,
    required this.tank,
    required this.health,
    required this.scoreGood,
    required this.scoreWarning,
    required this.scoreAlert,
    required this.summaryTitle,
    required this.noSummary,
    required this.visualAnalysis,
    required this.trendTitle,
    required this.nitrogenTitle,
    required this.radarTitle,
    required this.ph,
    required this.temperature,
    required this.current,
    required this.forecast,
    required this.reasoningTitle,
    required this.noReasoning,
    required this.actionTitle,
    required this.noActions,
    required this.footer,
  });

  final String title;
  final String reportId;
  final String testTime;
  final String tank;
  final String health;
  final String scoreGood;
  final String scoreWarning;
  final String scoreAlert;
  final String summaryTitle;
  final String noSummary;
  final String visualAnalysis;
  final String trendTitle;
  final String nitrogenTitle;
  final String radarTitle;
  final String ph;
  final String temperature;
  final String current;
  final String forecast;
  final String reasoningTitle;
  final String noReasoning;
  final String actionTitle;
  final String noActions;
  final String footer;

  factory _ReportLabels.forLanguage(ReportLanguage language) {
    switch (language) {
      case ReportLanguage.english:
        return const _ReportLabels(
          title: 'Aquarium Water Quality Health Report',
          reportId: 'Report ID',
          testTime: 'Generated',
          tank: 'Tank',
          health: 'Health',
          scoreGood: 'Good',
          scoreWarning: 'Watch',
          scoreAlert: 'Alert',
          summaryTitle: 'Diagnostic Summary',
          noSummary: 'No summary available.',
          visualAnalysis: 'Multi-dimensional Data Analysis',
          trendTitle: 'Recent Key Water Metrics',
          nitrogenTitle: 'Nitrogen Cycle Residue (mg/L)',
          radarTitle: 'Five-factor Water Quality Forecast',
          ph: 'pH',
          temperature: 'Temperature',
          current: 'Current',
          forecast: 'Forecast',
          reasoningTitle: 'Correlation Diagnosis',
          noReasoning: 'No reasoning available.',
          actionTitle: 'Recommended Actions',
          noActions: 'Keep the current maintenance plan.',
          footer: 'Generated automatically by AquaClaw Smart Aquarium',
        );
      case ReportLanguage.chinese:
        return const _ReportLabels(
          title: '鱼缸水质全谱系深度体检报告',
          reportId: '报告编号',
          testTime: '测定时间',
          tank: '缸体',
          health: '健康度',
          scoreGood: '良好',
          scoreWarning: '注意',
          scoreAlert: '预警',
          summaryTitle: '综合诊断结论',
          noSummary: '暂无概要。',
          visualAnalysis: '多维数据可视化分析',
          trendTitle: '近30条关键理化指标波动',
          nitrogenTitle: '氮循环毒素残留量 (mg/L)',
          radarTitle: '水质预测五维雷达',
          ph: 'pH值',
          temperature: '温度',
          current: '当前',
          forecast: '预测',
          reasoningTitle: '数据深度关联诊断',
          noReasoning: '暂无推理说明。',
          actionTitle: '针对性调水建议与行动指南',
          noActions: '继续保持当前养护方案。',
          footer: '此报告由 AquaClaw 智能鱼缸系统自动生成',
        );
    }
  }
}

class _LegendItem {
  const _LegendItem(this.label, this.color);

  final String label;
  final PdfColor color;
}

class _BarItem {
  const _BarItem(this.label, this.value, this.safeHigh, this.color);

  final String label;
  final double value;
  final double safeHigh;
  final String color;
}

class _Range {
  const _Range(this.low, this.high);

  final double low;
  final double high;

  factory _Range.fromValues(
    List<double> values,
    double safeLow,
    double safeHigh,
  ) {
    final clean = values.where((value) => value > 0).toList();
    if (clean.isEmpty) return _Range(safeLow, safeHigh);
    final minValue = clean.reduce(math.min);
    final maxValue = clean.reduce(math.max);
    final range = maxValue - minValue;
    final padding = math.max(range * 0.35, (safeHigh - safeLow) * 0.12);
    return _Range(
      math.min(safeLow, minValue - padding),
      math.max(safeHigh, maxValue + padding),
    );
  }

  double normalize(double value) {
    final span = high - low;
    if (span <= 0) return 0.5;
    return ((value - low) / span).clamp(0.0, 1.0);
  }
}

class _RadarMetric {
  const _RadarMetric(
    this.label,
    this.current,
    this.forecast,
    this.safeLow,
    this.safeHigh,
  );

  final String label;
  final double current;
  final double forecast;
  final double safeLow;
  final double safeHigh;

  double normalize(double value) {
    final span = safeHigh - safeLow;
    if (span <= 0) return 0.5;
    return ((value - safeLow) / span).clamp(0.0, 1.0);
  }
}
