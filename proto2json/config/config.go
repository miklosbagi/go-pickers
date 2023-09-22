package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// FieldOverrideConfig represents the YAML configuration structure
type FieldOverrideConfig struct {
	Service string            `yaml:"service"`
	Method  string            `yaml:"method"`
	Fields  map[string]string `yaml:"fields"`
}

func ReadConfigs(filePath string) ([]FieldOverrideConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var configWrapper struct {
		Overrides []FieldOverrideConfig `yaml:"overrides"`
	}
	if err := yaml.Unmarshal(data, &configWrapper); err != nil {
		return nil, err
	}

	return configWrapper.Overrides, nil
}
