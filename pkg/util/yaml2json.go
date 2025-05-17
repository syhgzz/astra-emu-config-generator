package util

import (
	"encoding/json"
	"fmt"

	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"gopkg.in/yaml.v2"
)

// convertMap 递归转换 map[interface{}]interface{} 为 map[string]interface{}
func convertMap(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range x {
			newMap[fmt.Sprint(k)] = convertMap(v)
		}
		return newMap
	case []interface{}:
		for i, v := range x {
			x[i] = convertMap(v)
		}
	}
	return i
}

func YamlToJson(yamlData []byte) (string, error) {
	var yamlDataInterface interface{}
	if err := yaml.Unmarshal(yamlData, &yamlDataInterface); err != nil {
		return "", fmt.Errorf("解析 YAML 失败: %v", err)
	}

	// 递归转换
	jsonCompatibleData := convertMap(yamlDataInterface)

	jsonData, err := json.MarshalIndent(jsonCompatibleData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("转换为 JSON 失败: %v", err)
	}

	return string(jsonData), nil
}

// YAMLToStruct 函数将 YAML 数据转换为结构体
func YAMLToStruct(yamlData []byte) (*kubenet_config.NetworkExperiment, error) {
	gotJSON, err := YamlToJson(yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %v", err)
	}

	var exp kubenet_config.NetworkExperiment
	err = exp.UnmarshalJSON([]byte(gotJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to struct: %v", err)
	}

	return &exp, nil
}
