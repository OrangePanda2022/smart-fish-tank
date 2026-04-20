import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aqua/domain/models/DeviceInfo.dart';

final devicesProvider = NotifierProvider<DevicesNotifier, List<Device>>(
  DevicesNotifier.new,
);

class DevicesNotifier extends Notifier<List<Device>> {
  @override
  List<Device> build() {
    return [
      Device(
        name: "主要照明",
        type: "light",
        isOn: true,
        intensity: 80,
        colorTemp: 6500,
      ),
      Device(name: "水泵", type: "pump", isOn: true, flow: "High"),
      Device(name: "加热器", type: "heater", isOn: true, setTemp: 26.5),
      Device(name: "氧气泵", type: "air", isOn: false, schedule: "Night"),
      Device(name: "自动投喂器", type: "feeder", isOn: true, nextFeed: "18:00"),
    ];
  }

  void toggleDevice(int index, bool value) {
    state = [
      for (int i = 0; i < state.length; i++)
        if (i == index)
          Device(
            name: state[i].name,
            type: state[i].type,
            isOn: value,
            intensity: state[i].intensity,
            colorTemp: state[i].colorTemp,
            flow: state[i].flow,
            setTemp: state[i].setTemp,
            schedule: state[i].schedule,
            nextFeed: state[i].nextFeed,
          )
        else
          state[i],
    ];
  }

  void updateIntensity(int index, double value) {
    state = [
      for (int i = 0; i < state.length; i++)
        if (i == index)
          Device(
            name: state[i].name,
            type: state[i].type,
            isOn: state[i].isOn,
            intensity: value,
            colorTemp: state[i].colorTemp,
            flow: state[i].flow,
            setTemp: state[i].setTemp,
            schedule: state[i].schedule,
            nextFeed: state[i].nextFeed,
          )
        else
          state[i],
    ];
  }

  void updateColorTemp(int index, double value) {
    state = [
      for (int i = 0; i < state.length; i++)
        if (i == index)
          Device(
            name: state[i].name,
            type: state[i].type,
            isOn: state[i].isOn,
            intensity: state[i].intensity,
            colorTemp: value,
            flow: state[i].flow,
            setTemp: state[i].setTemp,
            schedule: state[i].schedule,
            nextFeed: state[i].nextFeed,
          )
        else
          state[i],
    ];
  }
}
