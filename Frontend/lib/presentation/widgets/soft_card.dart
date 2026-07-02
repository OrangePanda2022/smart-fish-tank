import 'package:flutter/material.dart';

import '../../core/theme/app_shadow.dart';

class SoftCard extends StatelessWidget {
  const SoftCard({
    required this.child,
    this.padding = EdgeInsets.zero,
    this.radius = 24,
    super.key,
  });

  final Widget child;
  final EdgeInsetsGeometry padding;
  final double radius;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: padding,
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(radius),
        boxShadow: AppShadow.card,
      ),
      child: child,
    );
  }
}
