package kubenet_config_util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
)

func ReadConfig(filePath string) (*kubenet_config.NetworkExperiment, error) {
	// check if file is json or yaml
	if !strings.HasSuffix(filePath, ".yaml") && !strings.HasSuffix(filePath, ".json") {
		return nil, fmt.Errorf("unsupported file format: %v", filePath)
	}

	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if strings.HasSuffix(filePath, ".yaml") {
		jsonData, err := YamlToJson(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to convert YAML to JSON: %v", err)
		}
		rawData = jsonData
	}

	var exp kubenet_config.NetworkExperiment
	err = exp.UnmarshalJSON(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to struct: %v", err)
	}

	return &exp, nil
}

func WriteConfig(filePath string, exp *kubenet_config.NetworkExperiment) error {
	// check if file is json or yaml
	if !strings.HasSuffix(filePath, ".yaml") && !strings.HasSuffix(filePath, ".json") {
		return fmt.Errorf("unsupported file format: %v", filePath)
	}

	// create directory if not exists
	dir := filepath.Dir(filePath)
	// check if directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	rawData, err := exp.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal struct to JSON: %v", err)
	}

	if strings.HasSuffix(filePath, ".yaml") {
		rawData, err = JsonToYaml(rawData)
		if err != nil {
			return fmt.Errorf("failed to convert YAML to JSON: %v", err)
		}
	}

	err = os.WriteFile(filePath, rawData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
