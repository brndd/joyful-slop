// These types comprise the YAML schema that doesn't need custom unmarshalling.

package configparser

type Config struct {
	Devices []DeviceConfig
	Modes   []string
	Rules   []RuleConfig
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

type RuleConfigHat struct {
	Input  RuleTargetConfigHat
	Output RuleTargetConfigHat
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
