package configparser

// These top-level structs use custom unmarshaling to unpack each available sub-type
type DeviceConfig struct {
	Type   DeviceType
	Config interface{}
}

func (dc *DeviceConfig) UnmarshalYAML(unmarshal func(data interface{}) error) error {
	metaConfig := &struct {
		Type DeviceType
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
	}
	return err
}
