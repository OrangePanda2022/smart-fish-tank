import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import '../../domain/entities/aquarium_dashboard.dart';
import '../../domain/entities/sensor_reading.dart';
import '../../domain/repositories/aquarium_repository.dart';
import '../view_models/chart_series.dart';
import '../widgets/app_background.dart';
import '../widgets/chart_card.dart';
import '../widgets/headers.dart';
// import '../widgets/history_ai_card.dart';
import '../widgets/history_filters.dart';
import '../widgets/soft_card.dart';

class HistoryScreen extends StatefulWidget {
  const HistoryScreen({
    required this.dashboard,
    required this.repository,
    super.key,
  });

  final AquariumDashboard dashboard;
  final AquariumRepository repository;

  @override
  State<HistoryScreen> createState() => _HistoryScreenState();
}

class _HistoryScreenState extends State<HistoryScreen> {
  String _metric = 'ph';
  int _rangeIndex = 1;
  DateTimeRange? _customRange;
  List<SensorReading> _cachedHistory = const [];
  late Future<List<SensorReading>> _historyFuture;

  @override
  void initState() {
    super.initState();
    _historyFuture = _loadHistory();
  }

  Future<List<SensorReading>> _loadHistory() async {
    final history = await widget.repository.loadHistory(
      widget.dashboard.tank.id,
    );
    if (history.isNotEmpty) {
      _cachedHistory = List.unmodifiable(history);
    }
    return history;
  }

  void _refresh() {
    setState(() {
      _historyFuture = _loadHistory();
    });
  }

  Future<void> _pickCustomRange(List<SensorReading> history) async {
    final bounds = _historyBounds(history);
    final initialRange =
        _customRange ??
        DateTimeRange(
          start: bounds.end.subtract(const Duration(days: 6)),
          end: bounds.end,
        );

    final selected = await showDateRangePicker(
      context: context,
      firstDate: bounds.start.subtract(const Duration(days: 365)),
      lastDate: bounds.end.add(const Duration(days: 365)),
      initialDateRange: initialRange,
      helpText: '选择日期范围',
      cancelText: '取消',
      confirmText: '确定',
      saveText: '确定',
      fieldStartHintText: '开始日期',
      fieldEndHintText: '结束日期',
    );
    if (selected == null || !mounted) return;

    setState(() {
      _customRange = selected;
      _rangeIndex = 3;
    });
  }

  List<SensorReading> _filteredHistory(List<SensorReading> history) {
    if (history.isEmpty) return const [];

    final range = _activeRange(history);
    final filtered = history.where((reading) {
      final time = reading.timestamp;
      return !time.isBefore(range.start) && time.isBefore(range.end);
    }).toList()..sort((a, b) => a.timestamp.compareTo(b.timestamp));

    return filtered;
  }

  DateTimeRange _activeRange(List<SensorReading> history) {
    if (_rangeIndex == 3 && _customRange != null) {
      return DateTimeRange(
        start: _startOfDay(_customRange!.start),
        end: _dayAfter(_customRange!.end),
      );
    }

    final today = _startOfDay(DateTime.now());
    return switch (_rangeIndex) {
      0 => DateTimeRange(start: today, end: _dayAfter(today)),
      2 => DateTimeRange(
        start: today.subtract(const Duration(days: 29)),
        end: _dayAfter(today),
      ),
      _ => DateTimeRange(
        start: today.subtract(const Duration(days: 6)),
        end: _dayAfter(today),
      ),
    };
  }

  DateTimeRange _historyBounds(List<SensorReading> history) {
    if (history.isEmpty) {
      final now = DateTime.now();
      return DateTimeRange(start: now, end: now);
    }

    final sorted = [...history]
      ..sort((a, b) => a.timestamp.compareTo(b.timestamp));
    return DateTimeRange(
      start: sorted.first.timestamp,
      end: sorted.last.timestamp,
    );
  }

  DateTime _startOfDay(DateTime value) {
    return DateTime(value.year, value.month, value.day);
  }

  DateTime _dayAfter(DateTime value) {
    return _startOfDay(value).add(const Duration(days: 1));
  }

  String? _customRangeLabel() {
    final range = _customRange;
    if (range == null) return null;
    return '${_formatDate(range.start)} - ${_formatDate(range.end)}';
  }

  @override
  Widget build(BuildContext context) {
    return AppBackground(
      child: SafeArea(
        bottom: false,
        child: FutureBuilder<List<SensorReading>>(
          future: _historyFuture,
          builder: (context, snapshot) {
            final snapshotHistory = snapshot.data;
            final history = snapshotHistory?.isNotEmpty == true
                ? snapshotHistory!
                : _cachedHistory.isNotEmpty
                ? _cachedHistory
                : const <SensorReading>[];
            final filteredHistory = _filteredHistory(history);
            final series = ChartSeries.from(filteredHistory, _metric);
            final isLoading =
                snapshot.connectionState == ConnectionState.waiting;

            return ListView(
              padding: const EdgeInsets.fromLTRB(20, 8, 20, 36),
              children: [
                HistoryTopBar(onRefresh: _refresh),
                const SizedBox(height: 26),
                HistoryFilters(
                  metric: _metric,
                  rangeIndex: _rangeIndex,
                  customRangeLabel: _customRangeLabel(),
                  onMetricChanged: (value) => setState(() => _metric = value),
                  onRangeChanged: (index) =>
                      setState(() => _rangeIndex = index),
                  onCustomRangePressed: () => _pickCustomRange(history),
                ),
                if (isLoading) ...[
                  const SizedBox(height: 14),
                  const LinearProgressIndicator(minHeight: 3),
                ],
                const SizedBox(height: 24),
                ChartCard(series: series),
                const SizedBox(height: 18),
                HistoryDataTable(
                  history: filteredHistory,
                  metric: ChartSeries.metricFor(_metric),
                  totalCount: history.length,
                ),
                const SizedBox(height: 18),
                // HistoryAiCard(series: series),
                const _DisabledHistoryAiCard(),
                const SizedBox(height: 18),
                const Center(
                  child: Text(
                    '数据仅供参考，具体以实际情况为准',
                    style: TextStyle(color: AppColors.muted, fontSize: 13),
                  ),
                ),
              ],
            );
          },
        ),
      ),
    );
  }
}

class HistoryDataTable extends StatefulWidget {
  const HistoryDataTable({
    required this.history,
    required this.metric,
    required this.totalCount,
    super.key,
  });

  final List<SensorReading> history;
  final HistoryMetric metric;
  final int totalCount;

  @override
  State<HistoryDataTable> createState() => _HistoryDataTableState();
}

class _HistoryDataTableState extends State<HistoryDataTable> {
  static const _collapsedLimit = 10;

  bool _expanded = false;

  @override
  Widget build(BuildContext context) {
    final hasMore = widget.history.length > _collapsedLimit;
    final visibleHistory = _expanded || !hasMore
        ? widget.history
        : widget.history.take(_collapsedLimit).toList();

    return SoftCard(
      radius: 12,
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 10),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Expanded(child: Text('历史明细', style: AppText.titleSmall)),
              Text(
                '显示 ${visibleHistory.length}/${widget.totalCount} 条',
                style: AppText.caption,
              ),
            ],
          ),
          const SizedBox(height: 12),
          if (widget.history.isEmpty)
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 28),
              child: Center(child: Text('当前筛选范围内暂无数据', style: AppText.caption)),
            )
          else
            Scrollbar(
              thumbVisibility: true,
              child: SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: DataTable(
                  headingRowHeight: 42,
                  dataRowMinHeight: 44,
                  dataRowMaxHeight: 48,
                  columnSpacing: 22,
                  headingTextStyle: AppText.caption,
                  dataTextStyle: const TextStyle(
                    color: AppColors.ink,
                    fontSize: 13,
                    fontWeight: FontWeight.w700,
                  ),
                  columns: [
                    const DataColumn(label: Text('时间')),
                    DataColumn(label: Text(widget.metric.label)),
                    const DataColumn(label: Text('温度')),
                    const DataColumn(label: Text('pH')),
                    const DataColumn(label: Text('溶解氧')),
                    const DataColumn(label: Text('TDS')),
                    const DataColumn(label: Text('水位')),
                  ],
                  rows: visibleHistory.reversed.map((reading) {
                    return DataRow(
                      cells: [
                        DataCell(Text(_formatDateTime(reading.timestamp))),
                        DataCell(Text(metricValue(reading))),
                        DataCell(
                          Text('${reading.temperature.toStringAsFixed(1)}°C'),
                        ),
                        DataCell(Text(reading.ph.toStringAsFixed(2))),
                        DataCell(
                          Text('${reading.oxygen.toStringAsFixed(1)} mg/L'),
                        ),
                        DataCell(Text('${reading.tds.toStringAsFixed(0)} ppm')),
                        DataCell(
                          Text('${reading.waterLevel.toStringAsFixed(0)}%'),
                        ),
                      ],
                    );
                  }).toList(),
                ),
              ),
            ),
          if (hasMore) ...[
            const SizedBox(height: 8),
            TextButton(
              onPressed: () => setState(() => _expanded = !_expanded),
              child: Text(_expanded ? '收起' : '展开全部 ${widget.history.length} 条'),
            ),
          ],
        ],
      ),
    );
  }

  String metricValue(SensorReading reading) {
    final metric = widget.metric;
    return '${metric.valueOf(reading).toStringAsFixed(metric.precision)}${metric.unit}';
  }
}

class _DisabledHistoryAiCard extends StatelessWidget {
  const _DisabledHistoryAiCard();

  @override
  Widget build(BuildContext context) {
    return const SoftCard(
      radius: 12,
      padding: EdgeInsets.fromLTRB(18, 18, 18, 18),
      child: Row(
        children: [
          Icon(Icons.pause_circle_outline, color: AppColors.muted, size: 28),
          SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('AI 分析', style: AppText.blueHeading),
                SizedBox(height: 4),
                Text('暂不启用', style: AppText.caption),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

String _formatDate(DateTime value) {
  return '${value.month}/${value.day}';
}

String _formatDateTime(DateTime value) {
  return '${value.month}/${value.day} ${_two(value.hour)}:${_two(value.minute)}';
}

String _two(int value) => value.toString().padLeft(2, '0');
