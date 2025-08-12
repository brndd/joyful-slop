package virtualdevice

import (
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/suite"
)

type InitTests struct {
	suite.Suite
}

func TestRunnerInit(t *testing.T) {
	suite.Run(t, new(InitTests))
}

func (t *InitTests) TestMakeButtons() {
	t.Run("Maximum buttons", func() {
		buttons := makeButtons(VirtualDeviceMaxButtons, []string{})
		t.Equal(VirtualDeviceMaxButtons, len(buttons))
	})

	t.Run("Truncated buttons", func() {
		buttons := makeButtons(VirtualDeviceMaxButtons+1, []string{})
		t.Equal(VirtualDeviceMaxButtons, len(buttons))
	})

	t.Run("16 buttons", func() {
		buttons := makeButtons(16, []string{})
		t.Equal(16, len(buttons))
		t.Contains(buttons, evdev.EvCode(evdev.BTN_DEAD))
		t.NotContains(buttons, evdev.EvCode(evdev.BTN_TRIGGER_HAPPY))
	})

	t.Run("Explicit buttons", func() {
		buttonConfig := []string{"BTN_THUMB", "top", "btn_top2", "0x2fe", "0x300", "15"}
		buttons := makeButtons(0, buttonConfig)
		t.Equal(len(buttonConfig), len(buttons))
		t.Contains(buttons, evdev.EvCode(0x2fe))
		t.Contains(buttons, evdev.EvCode(0x300))
		t.Contains(buttons, evdev.EvCode(evdev.BTN_TOP))
		t.Contains(buttons, evdev.EvCode(evdev.BTN_DEAD))
	})
}

func (t *InitTests) TestMakeAxes() {
	t.Run("8 axes", func() {
		axes := makeAxes(8, []string{})
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
		axes := makeAxes(9, []string{})
		t.Equal(8, len(axes))
	})

	t.Run("3 axes", func() {
		axes := makeAxes(3, []string{})
		t.Equal(3, len(axes))
		t.Contains(axes, evdev.EvCode(evdev.ABS_X))
		t.Contains(axes, evdev.EvCode(evdev.ABS_Y))
		t.Contains(axes, evdev.EvCode(evdev.ABS_Z))
	})

	t.Run("4 explicit axis", func() {
		axes := makeAxes(0, []string{"x", "y", "throttle", "rudder"})
		t.Equal(4, len(axes))
		t.Contains(axes, evdev.EvCode(evdev.ABS_X))
		t.Contains(axes, evdev.EvCode(evdev.ABS_Y))
		t.Contains(axes, evdev.EvCode(evdev.ABS_THROTTLE))
		t.Contains(axes, evdev.EvCode(evdev.ABS_RUDDER))
	})
}

func (t *InitTests) TestMakeRelativeAxes() {
	t.Run("10 axes", func() {
		axes := makeRelativeAxes(10, []string{})
		t.Equal(10, len(axes))
		t.Contains(axes, evdev.EvCode(evdev.REL_X))
		t.Contains(axes, evdev.EvCode(evdev.REL_MISC))
	})

	t.Run("11 axes", func() {
		axes := makeRelativeAxes(11, []string{})
		t.Equal(10, len(axes))
	})

	t.Run("3 axes", func() {
		axes := makeRelativeAxes(3, []string{})
		t.Equal(3, len(axes))
		t.Contains(axes, evdev.EvCode(evdev.REL_X))
		t.Contains(axes, evdev.EvCode(evdev.REL_Y))
		t.Contains(axes, evdev.EvCode(evdev.REL_Z))
	})

	t.Run("1 explicit axis", func() {
		axes := makeRelativeAxes(0, []string{"wheel"})
		t.Equal(1, len(axes))
		t.Contains(axes, evdev.EvCode(evdev.REL_WHEEL))
	})

	t.Run("4 explicit axis", func() {
		axes := makeRelativeAxes(0, []string{"x", "y", "wheel", "hwheel"})
		t.Equal(4, len(axes))
		t.Contains(axes, evdev.EvCode(evdev.REL_X))
		t.Contains(axes, evdev.EvCode(evdev.REL_Y))
		t.Contains(axes, evdev.EvCode(evdev.REL_WHEEL))
		t.Contains(axes, evdev.EvCode(evdev.REL_HWHEEL))
	})
}
