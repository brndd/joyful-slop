package config

type Config struct {
	Devices DevicesConfig `yaml:"devices"`
	Rules   []RuleConfig  `yaml:"rules"`
}

type DevicesConfig struct {
	Physical []PhysicalDeviceConfig `yaml:"physical"`
	Virtual  []VirtualDeviceConfig  `yaml:"virtual"`
}

type PhysicalDeviceConfig struct {
	Name string `yaml:"name"`
	Uuid string `yaml:"uuid"`
}

type VirtualDeviceConfig struct {
	NumButtons int `yaml:"num_buttons"`
	NumAxes    int `yaml:"num_axes"`
}

type RuleConfig struct {
	Type   string           `yaml:"type"`
	Input  RuleInputConfig  `yaml:"input"`
	Output RuleOutputConfig `yaml:"output"`
}

type RuleInputConfig struct {
	Device   string   `yaml:"device"`
	Button   string   `yaml:"button,omitempty"`
	Buttons  []string `yaml:"buttons,omitempty"`
	Axis     string   `yaml:"axis,omitempty"`
	Inverted bool     `yaml:"inverted,omitempty"`
}

type RuleOutputConfig struct {
	Device string `yaml:"device"`
	Button string `yaml:"button,omitempty"`
	Axis   string `yaml:"axis,omitempty"`
}
