package configparser

type RuleConfig struct {
	Type   RuleType
	Name   string
	Modes  []string
	Config interface{}
}

func (dc *RuleConfig) UnmarshalYAML(unmarshal func(data interface{}) error) error {
	metaConfig := &struct {
		Type  RuleType
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
	}

	return err
}
