import 'package:flutter/material.dart';

import '../../core/theme/app_shadow.dart';

class CircleIconButton extends StatelessWidget {
  const CircleIconButton({
    required this.icon,
    required this.fill,
    required this.color,
    this.onPressed,
    super.key,
  });

  final IconData icon;
  final Color fill;
  final Color color;
  final VoidCallback? onPressed;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        customBorder: const CircleBorder(),
        onTap: onPressed,
        child: Container(
          width: 44,
          height: 44,
          alignment: Alignment.center,
          decoration: BoxDecoration(
            color: fill,
            shape: BoxShape.circle,
            boxShadow: fill == Colors.white ? AppShadow.soft : AppShadow.button,
          ),
          child: Icon(icon, color: color, size: 25),
        ),
      ),
    );
  }
}
