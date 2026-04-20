class AnalysisResult {
  final String reportId;
  final int statusScore;
  final String summary;
  final List<AnalysisAction> suggestions;
  final String reasoning;
  final DateTime timestamp;

  AnalysisResult({
    required this.reportId,
    required this.statusScore,
    required this.summary,
    required this.suggestions,
    required this.reasoning,
    required this.timestamp,
  });
}

class AnalysisAction {
  final String device;
  final String action;

  AnalysisAction({required this.device, required this.action});

  String get deviceDisplayName {
    switch (device) {
      case 'feeder':
        return '投喂器';
      case 'light':
        return '灯光';
      case 'heater':
        return '加热器';
      case 'pump':
        return '水泵';
      case 'aerator':
        return '增氧机';
      default:
        return device;
    }
  }

  String get actionDisplayName {
    switch (action) {
      case 'feed':
        return '投喂';
      case 'dim':
        return '调暗';
      case 'brighten':
        return '调亮';
      case 'heat':
        return '加热';
      case 'cool':
        return '降温';
      case 'increase_aeration':
        return '增加曝气';
      case 'change_water':
        return '换水';
      case 'add_oxygen':
        return '加氧';
      case 'maintain':
        return '保持现状';
      default:
        return action;
    }
  }
}
