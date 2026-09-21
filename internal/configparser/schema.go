// These types comprise the YAML schema that doesn't need custom unmarshalling.

package configparser

type Config struct {
	Devices []DeviceConfig
	Modes   []string
	Rules   []RuleConfig
}

// TODO: configure custom unmarshaling so we can overload Buttons, Axes, and RelativeAxes...
type DeviceConfigVirtual struct {
	Name string
	// VendorId is the uint16 vendor id, for default see internal/virtualdevice/variables.go
	VendorId string `yaml:"vendor_id,omitempty"`
	// DeviceId is the uint16 device id, for default see internal/virtualdevice/variables.go
	DeviceId        string `yaml:"device_id,omitempty"`
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
	Hold          bool
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

type RuleConfigButtonTempoBranch struct {
	Outputs []RuleTargetConfigButton
	Mode    string `yaml:"mode,omitempty"`
}

type RuleConfigButtonTempo struct {
	ThresholdMs int `yaml:"threshold_ms"`
	Input       RuleTargetConfigButton
	Tap         RuleConfigButtonTempoBranch
	Hold        RuleConfigButtonTempoBranch
}

type RuleConfigModeShift struct {
	Input RuleTargetConfigButton
	Mode  string
}
