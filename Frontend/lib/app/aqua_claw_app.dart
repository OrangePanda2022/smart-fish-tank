import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import '../core/i18n/app_language_scope.dart';
import '../core/i18n/app_localizations.dart';
import '../core/theme/app_theme.dart';
import '../presentation/shell/aquarium_shell.dart';

class AquaClawApp extends StatefulWidget {
  const AquaClawApp({super.key});

  @override
  State<AquaClawApp> createState() => _AquaClawAppState();
}

class _AquaClawAppState extends State<AquaClawApp> {
  Locale? _locale;

  @override
  Widget build(BuildContext context) {
    return AppLanguageScope(
      locale: _locale,
      setLocale: (locale) {
        setState(() {
          _locale = locale;
        });
      },
      child: MaterialApp(
        title: 'AquaClaw',
        debugShowCheckedModeBanner: false,
        theme: AppTheme.light(),
        locale: _locale,
        supportedLocales: AppLocalizations.supportedLocales,
        localizationsDelegates: const [
          AppLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        home: const AquariumShell(),
      ),
    );
  }
}
