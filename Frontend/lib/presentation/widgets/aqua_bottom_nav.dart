import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import '../../core/theme/app_shadow.dart';

class AquaBottomNav extends StatelessWidget {
  const AquaBottomNav({
    required this.selectedIndex,
    required this.onChanged,
    super.key,
  });

  final int selectedIndex;
  final ValueChanged<int> onChanged;

  @override
  Widget build(BuildContext context) {
    final items = [
      (CupertinoIcons.house_fill, '首页'),
      (CupertinoIcons.chart_bar_fill, '数据'),
      (CupertinoIcons.cube_box, '设备'),
      (CupertinoIcons.person, '我的'),
    ];

    return ClipRRect(
      borderRadius: BorderRadius.circular(24),
      child: Container(
        height: 78,
        decoration: BoxDecoration(
          color: Colors.white.withAlpha(245),
          boxShadow: AppShadow.nav,
        ),
        child: Row(
          children: List.generate(items.length, (index) {
            final selected = selectedIndex == index;
            return Expanded(
              child: GestureDetector(
                behavior: HitTestBehavior.opaque,
                onTap: () => onChanged(index),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      items[index].$1,
                      color: selected
                          ? AppColors.primary
                          : AppColors.secondaryText,
                      size: 28,
                    ),
                    const SizedBox(height: 6),
                    Text(
                      items[index].$2,
                      style: TextStyle(
                        color: selected
                            ? AppColors.primary
                            : AppColors.secondaryText,
                        fontSize: 13,
                        fontWeight: selected
                            ? FontWeight.w800
                            : FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            );
          }),
        ),
      ),
    );
  }
}
