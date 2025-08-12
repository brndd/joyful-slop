package configparser

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/goccy/go-yaml"
)

func ParseConfig(directory string) (*Config, error) {
	config := new(Config)

	configFiles, err := getConfigFilePaths(directory)
	if err != nil {
		return nil, err
	}

	// Open each yaml file and add its contents to the global config
	for _, filePath := range configFiles {
		data, err := os.ReadFile(filePath)
		if err != nil {
			logger.LogError(err, "Error while opening config file")
			continue
		}

		newConfig := Config{}
		err = yaml.Unmarshal(data, &newConfig)
		logger.LogIfError(err, "Error parsing YAML")
		config.Rules = append(config.Rules, newConfig.Rules...)
		config.Devices = append(config.Devices, newConfig.Devices...)
		config.Modes = append(config.Modes, newConfig.Modes...)
	}

	if len(config.Devices) == 0 {
		return nil, errors.New("Found no devices in configuration. Please add configuration at " + directory)
	}

	return config, nil
}

func getConfigFilePaths(directory string) ([]string, error) {
	paths := make([]string, 0)

	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		err = os.Mkdir(directory, 0755)
		if err != nil {
			return nil, errors.New("failed to create config directory at " + directory)
		} else {
			return nil, errors.New("no config files found at " + directory)
		}
	}

	for _, file := range dirEntries {
		name := strings.ToLower(file.Name())
		if file.IsDir() || !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")) {
			continue
		}

		paths = append(paths, filepath.Join(directory, file.Name()))
	}

	return paths, nil
}
