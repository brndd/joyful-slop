package configparser

type RuleTargetConfigButton struct {
	Device   string
	Button   string
	Inverted bool
}

type RuleTargetConfigAxis struct {
	Device    string
	Axis      string
	Inverted  bool
	Deadzones []DeadzoneConfig
}

type DeadzoneConfig struct {
	Center      int32 `yaml:"center,omitempty"`
	Size        int32 `yaml:"size,omitempty"`
	SizePercent int32 `yaml:"size_percent,omitempty"`
	Start       int32 `yaml:"start,omitempty"`
	End         int32 `yaml:"end,omitempty"`
	Emit        bool  `yaml:"emit,omitempty"`
	Value       int32 `yaml:"emit_value,omitempty"`
}

type RuleTargetConfigRelaxis struct {
	Device string
	Axis   string
}

type RuleTargetConfigModeSelect struct {
	Modes []string
}

type RuleTargetConfigHat struct {
	Device   string
	Hat      string
	Inverted bool
}
