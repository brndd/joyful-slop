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
	DevicePath      string   `yaml:"device_path,omitempty"`
	Preset          string   `yaml:"preset,omitempty"`
	NumButtons      int      `yaml:"num_buttons,omitempty"`
	NumAxes         int      `yaml:"num_axes,omitempty"`
	NumRelativeAxes int      `yaml:"num_rel_axes"`
	Buttons         []string `yaml:"buttons,omitempty"`
	Axes            []string `yaml:"axes,omitempty"`
	RelativeAxes    []string `yaml:"rel_axes,omitempty"`
	Lock            bool     `yaml:"lock,omitempty"`
}

type RuleConfig struct {
	Name          string             `yaml:"name,omitempty"`
	Type          string             `yaml:"type"`
	Input         RuleTargetConfig   `yaml:"input,omitempty"`
	InputLower    RuleTargetConfig   `yaml:"input_lower,omitempty"`
	InputUpper    RuleTargetConfig   `yaml:"input_upper,omitempty"`
	Inputs        []RuleTargetConfig `yaml:"inputs,omitempty"`
	Output        RuleTargetConfig   `yaml:"output"`
	Modes         []string           `yaml:"modes,omitempty"`
	RepeatRateMin int                `yaml:"repeat_rate_min,omitempty"`
	RepeatRateMax int                `yaml:"repeat_rate_max,omitempty"`
	Increment     int                `yaml:"increment,omitempty"`
}

type RuleTargetConfig struct {
	Device              string   `yaml:"device,omitempty"`
	Button              string   `yaml:"button,omitempty"`
	Axis                string   `yaml:"axis,omitempty"`
	DeadzoneCenter      int32    `yaml:"deadzone_center,omitempty"`
	DeadzoneSize        int32    `yaml:"deadzone_size,omitempty"`
	DeadzoneSizePercent int32    `yaml:"deadzone_size_percent,omitempty"`
	DeadzoneStart       int32    `yaml:"deadzone_start,omitempty"`
	DeadzoneEnd         int32    `yaml:"deadzone_end,omitempty"`
	Inverted            bool     `yaml:"inverted,omitempty"`
	Modes               []string `yaml:"modes,omitempty"`
}

// TODO: custom yaml unmarshaling is obtuse; do we really need to do all of this work
// just to set a single default value?
func (dc *DeviceConfig) UnmarshalYAML(unmarshal func(data interface{}) error) error {
	var raw struct {
		Name            string
		Type            string
		DeviceName      string `yaml:"device_name"`
		DevicePath      string `yaml:"device_path"`
		Preset          string
		NumButtons      int `yaml:"num_buttons"`
		NumAxes         int `yaml:"num_axes"`
		NumRelativeAxes int `yaml:"num_rel_axes"`
		Buttons         []string
		Axes            []string
		RelativeAxes    []string `yaml:"relative_axes"`
		Lock            bool     `yaml:"lock,omitempty"`
	}
	raw.Lock = true

	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	*dc = DeviceConfig{
		Name:            raw.Name,
		Type:            raw.Type,
		DeviceName:      raw.DeviceName,
		DevicePath:      raw.DevicePath,
		Preset:          raw.Preset,
		NumButtons:      raw.NumButtons,
		NumAxes:         raw.NumAxes,
		NumRelativeAxes: raw.NumRelativeAxes,
		Buttons:         raw.Buttons,
		Axes:            raw.Axes,
		RelativeAxes:    raw.RelativeAxes,
		Lock:            raw.Lock,
	}
	return nil
}
