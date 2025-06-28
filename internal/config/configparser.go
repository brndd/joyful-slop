package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/goccy/go-yaml"
)

type ConfigParser struct {
	config      Config
	configFiles []string
}

// Parse all the config files and store the config data for further use
func (parser *ConfigParser) Parse(directory string) {
	// Find the config files in the directory
	dirEntries, err := os.ReadDir(directory)
	logger.FatalIfError(err, "Failed to read config directory")

	for _, file := range dirEntries {
		name := file.Name()
		if file.IsDir() || !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")) {
			continue
		}

		filePath := filepath.Join(directory, name)
		if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
			parser.configFiles = append(parser.configFiles, filePath)
		}
	}

	rawData := parser.parseConfigFiles()
	err = yaml.Unmarshal(rawData, &parser.config)
	logger.FatalIfError(err, "Failed to parse config")
}

// Open each config file and concatenate their contents
func (parser *ConfigParser) parseConfigFiles() []byte {
	var rawData []byte

	for _, filePath := range parser.configFiles {
		data, err := os.ReadFile(filePath)
		if err != nil {
			logger.LogIfError(err, fmt.Sprintf("Failed to read config file '%s'", filePath))
			continue
		}

		rawData = append(rawData, data...)
		rawData = append(rawData, '\n')
	}

	if len(rawData) == 0 {
		logger.Log("No config data found")
		return nil
	}

	return rawData
}
