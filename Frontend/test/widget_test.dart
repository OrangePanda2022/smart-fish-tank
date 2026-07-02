import 'package:aquaclaw/app/aqua_claw_app.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('AquaClaw app starts on the smart aquarium home screen', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const AquaClawApp());
    await tester.pump();

    expect(find.text('我的鱼缸'), findsOneWidget);
    expect(find.text('客厅智能鱼缸'), findsNothing);
    expect(find.text('demo-tank'), findsNothing);
  });
}
