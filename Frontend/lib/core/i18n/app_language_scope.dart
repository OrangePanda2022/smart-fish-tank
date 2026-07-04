import 'package:flutter/material.dart';

class AppLanguageScope extends InheritedWidget {
  const AppLanguageScope({
    required this.locale,
    required this.setLocale,
    required super.child,
    super.key,
  });

  final Locale? locale;
  final ValueChanged<Locale?> setLocale;

  static AppLanguageScope of(BuildContext context) {
    final scope = context
        .dependOnInheritedWidgetOfExactType<AppLanguageScope>();
    assert(scope != null, 'No AppLanguageScope found in context');
    return scope!;
  }

  @override
  bool updateShouldNotify(AppLanguageScope oldWidget) {
    return oldWidget.locale != locale;
  }
}
