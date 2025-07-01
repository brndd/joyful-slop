// These types comprise the YAML schema for configuring Joyful.
// The config files will be combined and then unmarshalled into this

package config

type Config struct {
	Devices []DeviceConfig `yaml:"devices"`
	// TODO: add groups
	// Groups  []GroupConfig  `yaml:"groups,omitempty"`
	Rules []RuleConfig `yaml:"rules"`
}

type DeviceConfig struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	DeviceName string `yaml:"device_name,omitempty"`
	Uuid       string `yaml:"uuid,omitempty"`
	Buttons    int    `yaml:"buttons,omitempty"`
	Axes       int    `yaml:"axes,omitempty"`
}

type RuleConfig struct {
	Name   string            `yaml:"name,omitempty"`
	Type   string            `yaml:"type"`
	Input  RuleInputConfig   `yaml:"input,omitempty"`
	Inputs []RuleInputConfig `yaml:"inputs,omitempty"`
	Output RuleOutputConfig  `yaml:"output"`
}

type RuleInputConfig struct {
	Device   string   `yaml:"device"`
	Button   string   `yaml:"button,omitempty"`
	Buttons  []string `yaml:"buttons,omitempty"`
	Axis     string   `yaml:"axis,omitempty"`
	Inverted bool     `yaml:"inverted,omitempty"`
}

type RuleOutputConfig struct {
	Device string `yaml:"device,omitempty"`
	Button string `yaml:"button,omitempty"`
	Axis   string `yaml:"axis,omitempty"`
	Groups string `yaml:"groups,omitempty"`
}
