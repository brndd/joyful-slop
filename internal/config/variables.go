package config

import (
	"github.com/holoplot/go-evdev"
)

const (
	DeviceTypePhysical = "physical"
	DeviceTypeVirtual  = "virtual"

	RuleTypeButton       = "button"
	RuleTypeButtonCombo  = "button-combo"
	RuleTypeLatched      = "button-latched"
	RuleTypeAxis         = "axis"
	RuleTypeModeSelect   = "mode-select"
	RuleTypeAxisToButton = "axis-to-button"
)

var (
	ButtonFromIndex = []evdev.EvCode{
		evdev.BTN_TRIGGER,
		evdev.BTN_THUMB,
		evdev.BTN_THUMB2,
		evdev.BTN_TOP,
		evdev.BTN_TOP2,
		evdev.BTN_PINKIE,
		evdev.BTN_BASE,
		evdev.BTN_BASE2,
		evdev.BTN_BASE3,
		evdev.BTN_BASE4,
		evdev.BTN_BASE5,
		evdev.BTN_BASE6,
		evdev.BTN_TRIGGER_HAPPY1,
		evdev.BTN_TRIGGER_HAPPY2,
		evdev.BTN_TRIGGER_HAPPY3,
		evdev.BTN_TRIGGER_HAPPY4,
		evdev.BTN_TRIGGER_HAPPY5,
		evdev.BTN_TRIGGER_HAPPY6,
		evdev.BTN_TRIGGER_HAPPY7,
		evdev.BTN_TRIGGER_HAPPY8,
		evdev.BTN_TRIGGER_HAPPY9,
		evdev.BTN_TRIGGER_HAPPY10,
		evdev.BTN_TRIGGER_HAPPY11,
		evdev.BTN_TRIGGER_HAPPY12,
		evdev.BTN_TRIGGER_HAPPY13,
		evdev.BTN_TRIGGER_HAPPY14,
		evdev.BTN_TRIGGER_HAPPY15,
		evdev.BTN_TRIGGER_HAPPY16,
		evdev.BTN_TRIGGER_HAPPY17,
		evdev.BTN_TRIGGER_HAPPY18,
		evdev.BTN_TRIGGER_HAPPY19,
		evdev.BTN_TRIGGER_HAPPY20,
		evdev.BTN_TRIGGER_HAPPY21,
		evdev.BTN_TRIGGER_HAPPY22,
		evdev.BTN_TRIGGER_HAPPY23,
		evdev.BTN_TRIGGER_HAPPY24,
		evdev.BTN_TRIGGER_HAPPY25,
		evdev.BTN_TRIGGER_HAPPY26,
		evdev.BTN_TRIGGER_HAPPY27,
		evdev.BTN_TRIGGER_HAPPY28,
		evdev.BTN_TRIGGER_HAPPY29,
		evdev.BTN_TRIGGER_HAPPY30,
		evdev.BTN_TRIGGER_HAPPY31,
		evdev.BTN_TRIGGER_HAPPY32,
		evdev.BTN_TRIGGER_HAPPY33,
		evdev.BTN_TRIGGER_HAPPY34,
		evdev.BTN_TRIGGER_HAPPY35,
		evdev.BTN_TRIGGER_HAPPY36,
		evdev.BTN_TRIGGER_HAPPY37,
		evdev.BTN_TRIGGER_HAPPY38,
		evdev.BTN_TRIGGER_HAPPY39,
		evdev.BTN_TRIGGER_HAPPY40,
	}
)
