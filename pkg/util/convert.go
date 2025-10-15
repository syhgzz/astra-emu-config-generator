package util

import (
	"encoding/json"
	"fmt"

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

func YamlToJson(yamlData []byte) ([]byte, error) {
	var yamlDataInterface interface{}
	if err := yaml.Unmarshal(yamlData, &yamlDataInterface); err != nil {
		return nil, fmt.Errorf("解析 YAML 失败: %v", err)
	}

	// 递归转换
	jsonCompatibleData := convertMap(yamlDataInterface)

	jsonData, err := json.MarshalIndent(jsonCompatibleData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("转换为 JSON 失败: %v", err)
	}

	return jsonData, nil
}

func JsonToYaml(jsonData []byte) ([]byte, error) {
	var jsonDataInterface interface{}
	if err := json.Unmarshal(jsonData, &jsonDataInterface); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	// 递归转换
	yamlCompatibleData := convertMap(jsonDataInterface)

	yamlData, err := yaml.Marshal(yamlCompatibleData)
	if err != nil {
		return nil, fmt.Errorf("转换为 YAML 失败: %v", err)
	}

	return yamlData, nil
}
