import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/domain/models/tank_history.dart';
import 'package:aqua/data/services/tank_history_service.dart';

final tankHistoryServiceProvider = TankHistoryService();

final tankHistoryProvider = FutureProvider.family<TankHistoryData, StatType>((
  ref,
  type,
) async {
  final service = tankHistoryServiceProvider;
  return service.getHistory(type);
});
