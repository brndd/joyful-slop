package configparser

type DeviceConfigPhysical struct {
	Name       string
	DeviceName string `yaml:"device_name,omitempty"`
	DevicePath string `yaml:"device_path,omitempty"`
	Lock       bool
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
