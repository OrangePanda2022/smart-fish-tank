import 'package:aqua/presentation/widgets/SnackBar.dart';
import 'package:flutter/material.dart';

class AppearancePage extends StatefulWidget {
  const AppearancePage({super.key});

  @override
  State<AppearancePage> createState() => _AppearancePageState();
}

class _AppearancePageState extends State<AppearancePage> {
  final bool _darkMode = false;

  void _onDarkModeChanged(bool value) {
    AppSnackBar.warning(context, '此功能还没做完');
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(28)),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Center(
            child: Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.outline.withAlpha(77),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 20),
          Text('外观设置', style: Theme.of(context).textTheme.headlineSmall),
          const SizedBox(height: 24),
          Card(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(20),
            ),
            child: SwitchListTile(
              title: const Text('深色模式'),
              subtitle: const Text('开启后界面将变为深色主题'),
              secondary: Icon(
                _darkMode ? Icons.dark_mode : Icons.light_mode,
                color: Theme.of(context).colorScheme.onPrimaryContainer,
              ),
              value: _darkMode,
              onChanged: (value) => _onDarkModeChanged(value),
            ),
          ),
        ],
      ),
    );
  }
}
