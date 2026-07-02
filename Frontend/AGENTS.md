# Repository Guidelines

## Project Structure & Module Organization

This is a Flutter smart aquarium app using a lightweight DDD-style layout.

- `lib/main.dart` is only the app entry point.
- `lib/app/` contains the root `AquaClawApp`.
- `lib/domain/` contains entities and repository interfaces.
- `lib/infrastructure/` contains API and demo-data implementations.
- `lib/presentation/` contains screens, shell navigation, widgets, painters, and UI view models.
- `lib/core/` contains shared theme and utility code.
- `assets/images/` stores local UI images registered in `pubspec.yaml`.
- `test/` contains Flutter widget tests.

Keep new code in the layer that owns the responsibility. UI should depend on domain abstractions, not directly on backend clients.
Local persistence uses `shared_preferences`; keep cache key ownership in the infrastructure layer unless a broader storage abstraction is introduced.

## Build, Test, and Development Commands

- `flutter pub get` installs dependencies after `pubspec.yaml` changes.
- `dart format lib test` formats Dart source and tests.
- `flutter analyze` runs static analysis using `analysis_options.yaml`.
- `flutter test` runs widget/unit tests.
- `flutter run -d web-server --web-port 5173` starts a local web preview.
- `flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080` runs against a backend from Android emulator.

The default API base URL is `http://localhost:8080`.

## Coding Style & Naming Conventions

Use standard Dart formatting with two-space indentation. Prefer small, focused files over large mixed-responsibility files.

- Classes and widgets: `PascalCase`, e.g. `PredictionCard`.
- Files: `snake_case.dart`, e.g. `aqua_api_repository.dart`.
- Private helpers: prefix with `_` when scoped to one file.
- Keep constants and theme values in `lib/core/theme/`.

Do not add authentication, token, cookie, login, or registration code unless explicitly requested.

## Testing Guidelines

Use `flutter_test`. Add focused widget tests for user-visible flows and unit tests for parsing/repository behavior when adding API models.

Test names should describe behavior, e.g. `AquaClaw app starts on the smart aquarium home screen`.

Before handing off changes, run:

```bash
dart format lib test
flutter analyze
flutter test
```

## Commit & Pull Request Guidelines

This workspace now has a local Git repository on the `main` branch. Use concise conventional-style commit messages:

- `feat: add prediction summary parsing`
- `fix: isolate analyse API fallback`
- `refactor: split aquarium status widgets`

Pull requests should include a short summary, test results, screenshots for UI changes, and notes about backend endpoints or required `--dart-define` values.

## Architecture & Configuration Notes

Backend access lives in `AquaApiRepository`, which implements `AquariumRepository`. Keep API response parsing in infrastructure/domain entities and keep screens presentation-focused. Optional backend calls should fail gracefully and preserve real data already loaded from required endpoints.

The AI analyse action is exposed through `AquariumRepository.runAnalysis`. It calls `GET /api/v1/analyse/:tank_id` with a 60 second timeout and returns the structured `AnalysisReport` parsed from the response `data` field. Keep analyse response parsing in `lib/domain/entities/analysis_report.dart`.

The home health score is not sensor-estimated. It is read from the local `shared_preferences` cache key `analysis_score_<tankId>`, updated only after a successful analyse response using `decision.status_score`. If no cached score exists, the home status card should display `未分析`.

Dashboard loading should preserve the last valid dashboard and last known analysis score during refreshes. High-frequency analysis or refresh failures must not replace an already analysed home score with the demo dashboard's `null` score.

## Device & Profile Page Notes

`DeviceScreen` is the real device-management page shown from the bottom navigation. Keep device display data in `lib/presentation/view_models/device_catalog.dart` and `DeviceControl` until a backend device-control API is introduced. Device toggles are currently local UI state only; do not add backend control, authentication, token, cookie, login, or registration flows unless explicitly requested.

`DeviceDetailScreen` owns single-device presentation state such as the local enabled switch. Keep shared device UI in `lib/presentation/widgets/device_cards.dart`, and avoid duplicating the same device catalog in screens or widgets.

`ProfileScreen` is a presentation-only "我的" page. It should show profile, fish tank archive, preference toggles, and maintenance/support entries. Do not add a health score, shortcut entry section, cache status, backend URL display, or account/authentication features unless explicitly requested.

The home device grid and data-center sensor grid must stay responsive across narrow devices. When changing `DeviceGrid` or `SensorMetricGrid`, preserve their `LayoutBuilder`-based narrow-screen behavior so small screens avoid fixed three-column overflow.

## History Data Page Notes

`HistoryScreen` owns presentation-only history state: selected metric, selected date range, custom date range, refresh loading, and table expansion. Keep filtering local to the screen unless the backend API is explicitly changed to support server-side filters.

`AquariumRepository.loadHistory` should continue returning the full available history for a tank. The refresh button should reload that history without resetting the currently selected metric or date filter.

History loading should preserve the last non-empty history result as a cache. During high-frequency refreshes, empty or failed backend responses should keep showing the cached history instead of falling back to mismatched demo timestamps that can make the active filter show `暂无匹配数据`.

History metric definitions live in `lib/presentation/view_models/chart_series.dart` as `HistoryMetric` entries. Add new chart/table-selectable sensor metrics there so labels, units, precision, safe ranges, colors, and value extraction stay consistent across the filter picker, chart, stats, and table.

The history chart should use real reading timestamps for the x-axis and adapt the y-axis to the selected metric's values plus its safe range. Empty or single-point data sets must render gracefully without crashing.

The history details table should follow the active filtered history and show at most 10 rows by default. If more rows are available, keep them collapsed behind an expand/collapse action instead of rendering the full list immediately.

The history-data AI analysis card is currently disabled. Keep the original `HistoryAiCard` usage commented in `lib/presentation/screens/history_screen.dart` and show a clear `暂不启用` placeholder unless this feature is explicitly re-enabled.
