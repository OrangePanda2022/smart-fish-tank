import '../../core/utils/json_parsing.dart';

class TankInfo {
  const TankInfo({
    required this.id,
    required this.name,
    required this.size,
    required this.fishCount,
  });

  final String id;
  final String name;
  final int size;
  final int fishCount;

  factory TankInfo.fromJson(Map<String, dynamic> json) {
    return TankInfo(
      id: asString(json['tank_id']),
      name: asString(json['tank_name']),
      size: asInt(json['tank_size']),
      fishCount: asInt(json['fish_count']),
    );
  }
}
