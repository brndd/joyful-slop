// The ConfigParser is the main structure you'll interact with when using this package.
//
// Example usage:
//     config := &config.ConfigParser{}
//     config.Parse(<some directory containing YAML files>)
//     virtualDevices, err := config.CreateVirtualDevices()
//     physicalDevices, err := config.ConnectVirtualDevices()
//     modes, err := config.GetModes()
//     rules, err := config.BuildRules(physicalDevices, virtualDevices, modes)
//
// nb: there are methods defined on ConfigParser in other files in this package!

package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/goccy/go-yaml"
)

type ConfigParser struct {
	config Config
}

// Parse all the config files and store the config data for further use
func (parser *ConfigParser) Parse(directory string) error {
	parser.config = Config{}

	// Find the config files in the directory
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		err = os.Mkdir(directory, 0755)
		if err != nil {
			return errors.New("Failed to create config directory at " + directory)
		}
	}

	// Open each yaml file and add its contents to the global config
	for _, file := range dirEntries {
		name := file.Name()
		if file.IsDir() || !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")) {
			continue
		}

		filePath := filepath.Join(directory, name)
		if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
			data, err := os.ReadFile(filePath)
			if err != nil {
				logger.LogError(err, "Error while opening config file")
				continue
			}
			newConfig := Config{}
			err = yaml.Unmarshal(data, &newConfig)
			logger.LogIfError(err, "Error parsing YAML")
			parser.config.Rules = append(parser.config.Rules, newConfig.Rules...)
			parser.config.Devices = append(parser.config.Devices, newConfig.Devices...)
			// parser.config.Groups = append(parser.config.Groups, newConfig.Groups...)
		}
	}

	if len(parser.config.Devices) == 0 {
		return errors.New("Found no devices in configuration. Please add configuration at " + directory)
	}

	return nil
}
