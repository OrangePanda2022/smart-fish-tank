class Device {
  final String name;
  final String type;
  bool isOn;
  double? intensity; // Main Lighting
  double? colorTemp; // Main Lighting
  String? flow;      // Water Pump
  double? setTemp;   // Heater
  String? schedule;  // Air Pump
  String? nextFeed;  // Auto Feeder

  Device({
    required this.name,
    required this.type,
    this.isOn = false,
    this.intensity,
    this.colorTemp,
    this.flow,
    this.setTemp,
    this.schedule,
    this.nextFeed,
  });
}