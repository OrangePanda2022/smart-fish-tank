import 'package:flutter/material.dart';

import '../core/theme/app_theme.dart';
import '../presentation/shell/aquarium_shell.dart';

class AquaClawApp extends StatelessWidget {
  const AquaClawApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'AquaClaw',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light(),
      home: const AquariumShell(),
    );
  }
}
