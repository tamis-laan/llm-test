package utils

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

// Turn struct into pretty yaml string
func PrettyYAML(data interface{}) string {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return ""
	}
	return string(yamlBytes)
}

// Turn struct into pretty json string
func PrettyJSON(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// String array contains string
func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// String slice A contains All strings from B
func Subset(A []string, B []string) bool {
	for _, str := range A {
		if !Contains(B, str) {
			return false
		}
	}
	return true
}
