import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';
import 'soft_card.dart';

class WideActionButton extends StatelessWidget {
  const WideActionButton({
    required this.text,
    required this.icon,
    required this.onTap,
    super.key,
  });

  final String text;
  final IconData icon;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: SoftCard(
        radius: 12,
        padding: const EdgeInsets.symmetric(horizontal: 22, vertical: 20),
        child: Row(
          children: [
            Expanded(
              child: Center(child: Text(text, style: AppText.cardTitle)),
            ),
            Icon(icon, color: AppColors.secondaryText),
          ],
        ),
      ),
    );
  }
}
