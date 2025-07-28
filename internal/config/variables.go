package config

import (
	"github.com/holoplot/go-evdev"
)

const (
	DeviceTypePhysical = "physical"
	DeviceTypeVirtual  = "virtual"

	RuleTypeButton        = "button"
	RuleTypeButtonCombo   = "button-combo"
	RuleTypeLatched       = "button-latched"
	RuleTypeAxis          = "axis"
	RuleTypeAxisCombined  = "axis-combined"
	RuleTypeModeSelect    = "mode-select"
	RuleTypeAxisToButton  = "axis-to-button"
	RuleTypeAxisToRelaxis = "axis-to-relaxis"

	CodePrefixButton  = "BTN"
	CodePrefixAxis    = "ABS"
	CodePrefixRelaxis = "REL"

	VirtualDeviceMaxButtons = 74
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
		evdev.EvCode(0x12c), // decimal 300
		evdev.EvCode(0x12d), // decimal 301
		evdev.EvCode(0x12e), // decimal 302
		evdev.BTN_DEAD,
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
		evdev.EvCode(0x2e8),
		evdev.EvCode(0x2e9),
		evdev.EvCode(0x2f0),
		evdev.EvCode(0x2f1),
		evdev.EvCode(0x2f2),
		evdev.EvCode(0x2f3),
		evdev.EvCode(0x2f4),
		evdev.EvCode(0x2f5),
		evdev.EvCode(0x2f6),
		evdev.EvCode(0x2f7),
		evdev.EvCode(0x2f8),
		evdev.EvCode(0x2f9),
		evdev.EvCode(0x2fa),
		evdev.EvCode(0x2fb),
		evdev.EvCode(0x2fc),
		evdev.EvCode(0x2fd),
		evdev.EvCode(0x2fe),
		evdev.EvCode(0x2ff),
	}
)
