import 'dart:async';
import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

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
  final List<int> _buffer = [];

  http.Client? _client;
  StreamSubscription<List<int>>? _subscription;
  Uint8List? _frame;
  Object? _error;

  @override
  void initState() {
    super.initState();
    _connect();
  }

  @override
  void didUpdateWidget(MjpegStreamView oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.url != widget.url) {
      _connect();
    }
  }

  @override
  void dispose() {
    _disconnect();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final frame = _frame;
    if (frame != null) {
      return Image.memory(frame, fit: widget.fit, gaplessPlayback: true);
    }

    return ColoredBox(
      color: const Color(0xFF061B35),
      child: Center(
        child: _error == null
            ? const CircularProgressIndicator.adaptive()
            : const Icon(
                Icons.videocam_off_rounded,
                color: Colors.white70,
                size: 34,
              ),
      ),
    );
  }

  Future<void> _connect() async {
    _disconnect();
    _buffer.clear();
    setState(() {
      _frame = null;
      _error = null;
    });

    final uri = Uri.tryParse(widget.url);
    if (uri == null) {
      setState(() {
        _error = const FormatException('Invalid MJPEG stream URL');
      });
      return;
    }

    final client = http.Client();
    _client = client;

    try {
      final request = http.Request('GET', uri)
        ..headers['Accept'] = 'multipart/x-mixed-replace,image/jpeg,*/*';
      final response = await client.send(request);
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw HttpExceptionStatus(response.statusCode);
      }

      _subscription = response.stream.listen(
        _onChunk,
        onError: _setError,
        onDone: () {
          if (mounted && _frame == null) _setError(StateError('Stream ended'));
        },
        cancelOnError: true,
      );
    } catch (error) {
      _setError(error);
    }
  }

  void _disconnect() {
    _subscription?.cancel();
    _subscription = null;
    _client?.close();
    _client = null;
  }

  void _onChunk(List<int> chunk) {
    _buffer.addAll(chunk);

    while (true) {
      final start = _indexOfMarker(_buffer, 0xFF, 0xD8, 0);
      if (start < 0) {
        if (_buffer.length > 4096) _buffer.clear();
        return;
      }

      if (start > 0) _buffer.removeRange(0, start);

      final end = _indexOfMarker(_buffer, 0xFF, 0xD9, 2);
      if (end < 0) return;

      final frame = Uint8List.fromList(_buffer.sublist(0, end + 2));
      _buffer.removeRange(0, end + 2);

      if (mounted) {
        setState(() {
          _frame = frame;
          _error = null;
        });
      }
    }
  }

  int _indexOfMarker(List<int> bytes, int first, int second, int from) {
    for (var index = from; index < bytes.length - 1; index += 1) {
      if (bytes[index] == first && bytes[index + 1] == second) return index;
    }
    return -1;
  }

  void _setError(Object error) {
    if (!mounted) return;
    setState(() {
      _error = error;
    });
  }
}

class HttpExceptionStatus implements Exception {
  const HttpExceptionStatus(this.statusCode);

  final int statusCode;

  @override
  String toString() => 'HTTP $statusCode';
}
