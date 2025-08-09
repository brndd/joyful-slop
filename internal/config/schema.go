// These types comprise the YAML schema for configuring Joyful.
// The config files will be combined and then unmarshalled into this

package config

import (
	"fmt"
)

type Config struct {
	Devices []DeviceConfig
	Modes   []string
	Rules   []RuleConfig
}

// These top-level structs use custom unmarshaling to unpack each available sub-type
type DeviceConfig struct {
	Type   string
	Config interface{}
}

type RuleConfig struct {
	Type   string
	Name   string
	Modes  []string
	Config interface{}
}

type DeviceConfigPhysical struct {
	Name       string
	DeviceName string `yaml:"device_name,omitempty"`
	DevicePath string `yaml:"device_path,omitempty"`
	Lock       bool
}

// TODO: configure custom unmarshaling so we can overload Buttons, Axes, and RelativeAxes...
type DeviceConfigVirtual struct {
	Name            string
	Preset          string
	NumButtons      int `yaml:"num_buttons,omitempty"`
	NumAxes         int `yaml:"num_axes,omitempty"`
	NumRelativeAxes int `yaml:"num_rel_axes"`
	Buttons         []string
	Axes            []string
	RelativeAxes    []string `yaml:"rel_axes,omitempty"`
}

type RuleConfigButton struct {
	Input  RuleTargetConfigButton
	Output RuleTargetConfigButton
}

type RuleConfigButtonCombo struct {
	Inputs []RuleTargetConfigButton
	Output RuleTargetConfigButton
}

type RuleConfigButtonLatched struct {
	Input  RuleTargetConfigButton
	Output RuleTargetConfigButton
}

type RuleConfigAxis struct {
	Input  RuleTargetConfigAxis
	Output RuleTargetConfigAxis
}

type RuleConfigAxisCombined struct {
	InputLower RuleTargetConfigAxis `yaml:"input_lower,omitempty"`
	InputUpper RuleTargetConfigAxis `yaml:"input_upper,omitempty"`
	Output     RuleTargetConfigAxis
}

type RuleConfigAxisToButton struct {
	RepeatRateMin int `yaml:"repeat_rate_min,omitempty"`
	RepeatRateMax int `yaml:"repeat_rate_max,omitempty"`
	Input         RuleTargetConfigAxis
	Output        RuleTargetConfigButton
}

type RuleConfigAxisToRelaxis struct {
	RepeatRateMin int `yaml:"repeat_rate_min"`
	RepeatRateMax int `yaml:"repeat_rate_max"`
	Increment     int
	Input         RuleTargetConfigAxis
	Output        RuleTargetConfigRelaxis
}

type RuleConfigModeSelect struct {
	Input  RuleTargetConfigButton
	Output RuleTargetConfigModeSelect
}

type RuleTargetConfigButton struct {
	Device   string
	Button   string
	Inverted bool
}

type RuleTargetConfigAxis struct {
	Device              string
	Axis                string
	DeadzoneCenter      int32 `yaml:"deadzone_center,omitempty"`
	DeadzoneSize        int32 `yaml:"deadzone_size,omitempty"`
	DeadzoneSizePercent int32 `yaml:"deadzone_size_percent,omitempty"`
	DeadzoneStart       int32 `yaml:"deadzone_start,omitempty"`
	DeadzoneEnd         int32 `yaml:"deadzone_end,omitempty"`
	Inverted            bool
}

type RuleTargetConfigRelaxis struct {
	Device string
	Axis   string
}

type RuleTargetConfigModeSelect struct {
	Modes []string
}

func (dc *DeviceConfig) UnmarshalYAML(unmarshal func(data interface{}) error) error {
	metaConfig := &struct {
		Type string
	}{}
	err := unmarshal(metaConfig)
	if err != nil {
		return err
	}
	dc.Type = metaConfig.Type

	err = nil
	switch metaConfig.Type {
	case DeviceTypePhysical:
		config := DeviceConfigPhysical{}
		err = unmarshal(&config)
		dc.Config = config
	case DeviceTypeVirtual:
		config := DeviceConfigVirtual{}
		err = unmarshal(&config)
		dc.Config = config
	default:
		err = fmt.Errorf("invalid device type '%s'", dc.Type)
	}
	return err
}

func (dc *RuleConfig) UnmarshalYAML(unmarshal func(data interface{}) error) error {
	metaConfig := &struct {
		Type  string
		Name  string
		Modes []string
	}{}
	err := unmarshal(metaConfig)
	if err != nil {
		return err
	}
	dc.Type = metaConfig.Type
	dc.Name = metaConfig.Name
	dc.Modes = metaConfig.Modes

	switch dc.Type {
	case RuleTypeButton:
		config := RuleConfigButton{}
		err = unmarshal(&config)
		dc.Config = config
	case RuleTypeButtonCombo:
		config := RuleConfigButtonCombo{}
		err = unmarshal(&config)
		dc.Config = config
	case RuleTypeButtonLatched:
		config := RuleConfigButtonLatched{}
		err = unmarshal(&config)
		dc.Config = config
	case RuleTypeAxis:
		config := RuleConfigAxis{}
		err = unmarshal(&config)
		dc.Config = config
	case RuleTypeAxisCombined:
		config := RuleConfigAxisCombined{}
		err = unmarshal(&config)
		dc.Config = config
	case RuleTypeAxisToButton:
		config := RuleConfigAxisToButton{}
		err = unmarshal(&config)
		dc.Config = config
	case RuleTypeAxisToRelaxis:
		config := RuleConfigAxisToRelaxis{}
		err = unmarshal(&config)
		dc.Config = config
	case RuleTypeModeSelect:
		config := RuleConfigModeSelect{}
		err = unmarshal(&config)
		dc.Config = config
	default:
		err = fmt.Errorf("invalid rule type '%s'", dc.Type)
	}

	return err
}

// TODO: custom yaml unmarshaling is obtuse; do we really need to do all of this work
// just to set a single default value?
func (dc *DeviceConfigPhysical) UnmarshalYAML(unmarshal func(data interface{}) error) error {
	var raw struct {
		Name       string
		DeviceName string `yaml:"device_name"`
		DevicePath string `yaml:"device_path"`
		Lock       bool   `yaml:"lock,omitempty"`
	}

	// Set non-standard defaults
	raw.Lock = true

	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	*dc = DeviceConfigPhysical{
		Name:       raw.Name,
		DeviceName: raw.DeviceName,
		DevicePath: raw.DevicePath,
		Lock:       raw.Lock,
	}
	return nil
}
