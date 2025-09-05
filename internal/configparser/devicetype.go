package configparser

import (
	"fmt"
	"strings"
)

type DeviceType string

const (
	DeviceTypeNone     DeviceType = ""
	DeviceTypePhysical DeviceType = "physical"
	DeviceTypeVirtual  DeviceType = "virtual"
)

var (
	deviceTypeMap = map[string]DeviceType{
		"physical": DeviceTypePhysical,
		"virtual":  DeviceTypeVirtual,
	}
)

func ParseDeviceType(in string) (DeviceType, error) {
	deviceType, ok := deviceTypeMap[strings.ToLower(in)]
	if !ok {
		return DeviceTypeNone, fmt.Errorf("invalid rule type '%s'", in)
	}
	return deviceType, nil
}

func (rt *DeviceType) UnmarshalYAML(unmarshal func(data interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	*rt, err = ParseDeviceType(raw)
	return err
}
