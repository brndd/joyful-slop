package config

import (
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/suite"
)

type DevicesConfigTests struct {
	suite.Suite
}

func TestRunnerDevicesConfig(t *testing.T) {
	suite.Run(t, new(DevicesConfigTests))
}

func (t *DevicesConfigTests) TestMakeAxes() {
	t.Run("8 axes", func() {
		axes := makeAxes(8)
		t.Equal(8, len(axes))
		t.Contains(axes, evdev.EvCode(evdev.ABS_X))
		t.Contains(axes, evdev.EvCode(evdev.ABS_Y))
		t.Contains(axes, evdev.EvCode(evdev.ABS_Z))
		t.Contains(axes, evdev.EvCode(evdev.ABS_RX))
		t.Contains(axes, evdev.EvCode(evdev.ABS_RY))
		t.Contains(axes, evdev.EvCode(evdev.ABS_RZ))
		t.Contains(axes, evdev.EvCode(evdev.ABS_THROTTLE))
		t.Contains(axes, evdev.EvCode(evdev.ABS_RUDDER))
	})

	t.Run("9 axes is truncated", func() {
		axes := makeAxes(9)
		t.Equal(8, len(axes))
	})

	t.Run("3 axes", func() {
		axes := makeAxes(3)
		t.Equal(3, len(axes))
		t.Contains(axes, evdev.EvCode(evdev.ABS_X))
		t.Contains(axes, evdev.EvCode(evdev.ABS_Y))
		t.Contains(axes, evdev.EvCode(evdev.ABS_Z))
	})
}

func (t *DevicesConfigTests) TestMakeButtons() {
	t.Run("Maximum buttons", func() {
		buttons := makeButtons(VirtualDeviceMaxButtons)
		t.Equal(VirtualDeviceMaxButtons, len(buttons))
	})

	t.Run("Truncated buttons", func() {
		buttons := makeButtons(VirtualDeviceMaxButtons + 1)
		t.Equal(VirtualDeviceMaxButtons, len(buttons))
	})

	t.Run("16 buttons", func() {
		buttons := makeButtons(16)
		t.Equal(16, len(buttons))
		t.Contains(buttons, evdev.EvCode(evdev.BTN_DEAD))
		t.NotContains(buttons, evdev.EvCode(evdev.BTN_TRIGGER_HAPPY))
	})
}
