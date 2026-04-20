import 'package:flutter/material.dart';

class ImageCard extends StatelessWidget {
  final String imageUrl;
  final double borderRadius;
  final double elevation;
  final double width;
  final double height;
  final BoxFit fit;

  const ImageCard({
    required this.imageUrl,
    this.borderRadius = 16.0,
    this.elevation = 0.0,
    this.width = 200,
    this.height = 200,
    this.fit = BoxFit.cover,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Material(
      elevation: elevation,
      borderRadius: BorderRadius.circular(borderRadius),
      color: Colors.transparent,
      child: ClipRRect(
        borderRadius: BorderRadius.circular(borderRadius),
        child: Image.asset(
          imageUrl,
          width: width,
          height: height,
          fit: BoxFit.contain,
          // loadingBuilder: (context, child, progress) {
          //   if (progress == null) return child;
          //   return Center(
          //     child: CircularProgressIndicator(
          //       value: progress.expectedTotalBytes != null
          //           ? progress.cumulativeBytesLoaded / progress.expectedTotalBytes!
          //           : null,
          //     ),
          //   );
          // },
          errorBuilder: (context, error, stackTrace) {
            return Container(
              width: width,
              height: height,
              color: Colors.grey[300],
              child: const Icon(Icons.broken_image, size: 40, color: Colors.grey),
            );
          },
        ),
      ),
    );
  }
}