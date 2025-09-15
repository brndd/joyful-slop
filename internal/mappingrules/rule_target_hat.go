package mappingrules

import (
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/eventcodes"
	"github.com/holoplot/go-evdev"
)

type RuleTargetHat struct {
	Device   Device
	Hat      evdev.EvCode
	Inverted bool
}

func NewRuleTargetHatFromConfig(config configparser.RuleTargetConfigHat, devs map[string]Device) (*RuleTargetHat, error) {
	dev, ok := devs[config.Device]
	if !ok {
		return nil, fmt.Errorf("device '%s' not found", config.Device)
	}

	code, err := eventcodes.ParseCode(config.Hat, eventcodes.CodePrefixAxis)
	if err != nil {
		return nil, err
	}

	return &RuleTargetHat{
		Device:   dev,
		Hat:      code,
		Inverted: config.Inverted,
	}, nil
}

func (target *RuleTargetHat) NormalizeValue(value int32) int32 {
	if !target.Inverted {
		return value
	}

	return value * -1
}

func (target *RuleTargetHat) CreateEvent(value int32, _ *string) *evdev.InputEvent {
	return &evdev.InputEvent{
		Type:  evdev.EV_ABS,
		Code:  target.Hat,
		Value: value,
	}
}

func (target *RuleTargetHat) MatchEvent(device Device, event *evdev.InputEvent) bool {
	return device == target.Device && event.Code == target.Hat
}
