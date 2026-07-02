import 'package:flutter/cupertino.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_text.dart';

class SectionHeader extends StatelessWidget {
  const SectionHeader({
    required this.title,
    required this.actionText,
    required this.onAction,
    this.showChevron = true,
    this.compact = false,
    super.key,
  });

  final String title;
  final String actionText;
  final VoidCallback onAction;
  final bool showChevron;
  final bool compact;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(title, style: compact ? AppText.sectionCompact : AppText.section),
        const Spacer(),
        GestureDetector(
          onTap: onAction,
          child: Row(
            children: [
              Text(actionText, style: AppText.caption),
              if (showChevron) ...[
                const SizedBox(width: 4),
                const Icon(
                  CupertinoIcons.chevron_right,
                  color: AppColors.muted,
                  size: 18,
                ),
              ],
            ],
          ),
        ),
      ],
    );
  }
}
