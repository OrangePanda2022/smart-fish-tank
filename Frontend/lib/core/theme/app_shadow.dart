import 'package:flutter/material.dart';

import 'app_colors.dart';

class AppShadow {
  static final card = [
    BoxShadow(
      color: AppColors.primary.withAlpha(20),
      blurRadius: 24,
      offset: const Offset(0, 12),
    ),
  ];

  static final soft = [
    BoxShadow(
      color: AppColors.primary.withAlpha(14),
      blurRadius: 18,
      offset: const Offset(0, 8),
    ),
  ];

  static final button = [
    BoxShadow(
      color: AppColors.primary.withAlpha(72),
      blurRadius: 18,
      offset: const Offset(0, 8),
    ),
  ];

  static final nav = [
    BoxShadow(
      color: AppColors.primary.withAlpha(18),
      blurRadius: 30,
      offset: const Offset(0, -8),
    ),
  ];
}
