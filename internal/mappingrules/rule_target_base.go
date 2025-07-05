package mappingrules

import "github.com/holoplot/go-evdev"

type RuleTargetBase struct {
	DeviceName string
	Device     *evdev.InputDevice
	Code       evdev.EvCode
	Inverted   bool
}

func NewRuleTargetBase(device_name string,
	device *evdev.InputDevice,
	code evdev.EvCode,
	inverted bool) RuleTargetBase {

	return RuleTargetBase{
		DeviceName: device_name,
		Device:     device,
		Code:       code,
		Inverted:   inverted,
	}
}

func (target *RuleTargetBase) GetCode() evdev.EvCode {
	return target.Code
}

func (target *RuleTargetBase) GetDeviceName() string {
	return target.DeviceName
}

func (target *RuleTargetBase) GetDevice() *evdev.InputDevice {
	return target.Device
}
