import 'package:flutter/material.dart';

class SettingsDetailPage extends StatelessWidget {
  final String title;
  final Widget child;
  final Widget? trailing;

  const SettingsDetailPage({
    super.key,
    required this.title,
    required this.child,
    this.trailing,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(title),
        actions: trailing != null ? [trailing!] : null,
      ),
      body: child,
    );
  }
}
