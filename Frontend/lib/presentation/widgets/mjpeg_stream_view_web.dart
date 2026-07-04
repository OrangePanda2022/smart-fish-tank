import 'dart:ui_web' as ui_web;

import 'package:flutter/material.dart';
import 'package:web/web.dart' as web;

class MjpegStreamView extends StatefulWidget {
  const MjpegStreamView({
    required this.url,
    this.fit = BoxFit.cover,
    super.key,
  });

  final String url;
  final BoxFit fit;

  @override
  State<MjpegStreamView> createState() => _MjpegStreamViewState();
}

class _MjpegStreamViewState extends State<MjpegStreamView> {
  static final Set<String> _registeredViewTypes = {};

  late String _viewType;

  @override
  void initState() {
    super.initState();
    _registerView();
  }

  @override
  void didUpdateWidget(MjpegStreamView oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.url != widget.url || oldWidget.fit != widget.fit) {
      _registerView();
    }
  }

  @override
  Widget build(BuildContext context) {
    return HtmlElementView(key: ValueKey(_viewType), viewType: _viewType);
  }

  void _registerView() {
    _viewType = 'aquaclaw-mjpeg-${widget.url.hashCode}-${widget.fit.name}';
    if (_registeredViewTypes.contains(_viewType)) return;

    ui_web.platformViewRegistry.registerViewFactory(_viewType, (viewId) {
      return web.HTMLImageElement()
        ..src = widget.url
        ..alt = 'Aquarium camera stream'
        ..style.width = '100%'
        ..style.height = '100%'
        ..style.display = 'block'
        ..style.objectFit = _objectFit(widget.fit)
        ..style.backgroundColor = '#061B35';
    });
    _registeredViewTypes.add(_viewType);
  }

  String _objectFit(BoxFit fit) {
    return switch (fit) {
      BoxFit.contain => 'contain',
      BoxFit.fill => 'fill',
      BoxFit.fitWidth => 'cover',
      BoxFit.fitHeight => 'cover',
      BoxFit.none => 'none',
      BoxFit.scaleDown => 'scale-down',
      BoxFit.cover => 'cover',
    };
  }
}
