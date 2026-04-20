import 'package:dio/dio.dart';
import 'package:aqua/data/models/ai_api_response.dart';

class AiApiService {
  static const String baseUrl = 'http://192.168.1.4:6788';

  final Dio _dio;

  AiApiService()
    : _dio = Dio(
        BaseOptions(
          baseUrl: baseUrl,
          connectTimeout: const Duration(seconds: 30),
          receiveTimeout: const Duration(seconds: 120),
          headers: {'Content-Type': 'application/json'},
        ),
      );

  Future<AiApiResponse> analyse(String tankId) async {
    try {
      final response = await _dio.get('/api/v1/analyse/$tankId');
      return AiApiResponse.fromJson(response.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw AiApiException(
        message: e.message ?? 'Network error',
        statusCode: e.response?.statusCode,
      );
    } catch (e) {
      throw AiApiException(message: e.toString());
    }
  }
}

class AiApiException implements Exception {
  final String message;
  final int? statusCode;

  AiApiException({required this.message, this.statusCode});

  @override
  String toString() => 'AiApiException: $message (status: $statusCode)';
}
