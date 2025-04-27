package util

import (
	"encoding/json"
	"fmt"

	"dslab.sjtu/kubenet-config-generator/api/pkg/kubenet_config"
	"gopkg.in/yaml.v2"
)

func JsonToYaml(jsonData []byte) (string, error) {
	var jsonDataInterface interface{}
	if err := json.Unmarshal(jsonData, &jsonDataInterface); err != nil {
		return "", fmt.Errorf("解析 JSON 失败: %v", err)
	}

	// 递归转换
	yamlCompatibleData := convertMap(jsonDataInterface)

	yamlData, err := yaml.Marshal(yamlCompatibleData)
	if err != nil {
		return "", fmt.Errorf("转换为 YAML 失败: %v", err)
	}

	return string(yamlData), nil
}

// StructToYAML 函数将结构体转换为 YAML 数据
func StructToYAML(exp *kubenet_config.NetworkExperiment) ([]byte, error) {
	var gotJSON []byte
	gotJSON, err := exp.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to JSON: %v", err)
	}

	yamlData, err := JsonToYaml(gotJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JSON to YAML: %v", err)
	}

	return []byte(yamlData), nil
}
