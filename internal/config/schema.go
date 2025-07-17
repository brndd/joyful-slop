// These types comprise the YAML schema for configuring Joyful.
// The config files will be combined and then unmarshalled into this
//
// TODO: currently the types in here aren't especially strong; each one is
// decomposed into a different object based on the Type fields. We should implement
// some sort of delayed unmarshalling technique, for example see ideas at
// https://stackoverflow.com/questions/70635636/unmarshaling-yaml-into-different-struct-based-off-yaml-field
// Then we can be more explicit about the interface here.

package config

type Config struct {
	Devices []DeviceConfig `yaml:"devices"`
	Modes   []string       `yaml:"modes,omitempty"`
	Rules   []RuleConfig   `yaml:"rules"`
}

type DeviceConfig struct {
	Name            string   `yaml:"name"`
	Type            string   `yaml:"type"`
	DeviceName      string   `yaml:"device_name,omitempty"`
	Uuid            string   `yaml:"uuid,omitempty"`
	NumButtons      int      `yaml:"num_buttons,omitempty"`
	NumAxes         int      `yaml:"num_axes,omitempty"`
	NumRelativeAxes int      `yaml:"num_rel_axes"`
	Buttons         []string `yaml:"buttons,omitempty"`
	Axes            []string `yaml:"axes,omitempty"`
	RelativeAxes    []string `yaml:"rel_axes,omitempty"`
}

type RuleConfig struct {
	Name          string             `yaml:"name,omitempty"`
	Type          string             `yaml:"type"`
	Input         RuleTargetConfig   `yaml:"input,omitempty"`
	Inputs        []RuleTargetConfig `yaml:"inputs,omitempty"`
	Output        RuleTargetConfig   `yaml:"output"`
	Modes         []string           `yaml:"modes,omitempty"`
	RepeatRateMin int                `yaml:"repeat_rate_min,omitempty"`
	RepeatRateMax int                `yaml:"repeat_rate_max,omitempty"`
	Increment     int                `yaml:"increment,omitempty"`
}

type RuleTargetConfig struct {
	Device        string   `yaml:"device,omitempty"`
	Button        string   `yaml:"button,omitempty"`
	Axis          string   `yaml:"axis,omitempty"`
	DeadzoneStart int32    `yaml:"deadzone_start,omitempty"`
	DeadzoneEnd   int32    `yaml:"deadzone_end,omitempty"`
	Inverted      bool     `yaml:"inverted,omitempty"`
	Modes         []string `yaml:"modes,omitempty"`
}
