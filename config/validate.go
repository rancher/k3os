package config

import (
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

func KeysToStrings(item interface{}) interface{} {
	switch d := item.(type) {
	case map[string]interface{}:
		for key, value := range d {
			d[key] = KeysToStrings(value)
		}
		return d
	case map[interface{}]interface{}:
		n := make(map[string]interface{})
		for key, value := range d {
			s := key.(string)
			n[s] = KeysToStrings(value)
		}
		return n
	case []interface{}:
		for i, value := range d {
			d[i] = KeysToStrings(value)
		}
		return d
	default:
		return item
	}
}

func ValidateCfg(rawCfg interface{}) (*gojsonschema.Result, error) {
	rawCfg = KeysToStrings(rawCfg).(map[string]interface{})
	loader := gojsonschema.NewGoLoader(rawCfg)
	schemaLoader := gojsonschema.NewStringLoader(schema)
	return gojsonschema.Validate(schemaLoader, loader)
}

func ValidateBytes(bytes []byte) (*gojsonschema.Result, error) {
	var rawCfg map[string]interface{}
	if err := yaml.Unmarshal([]byte(bytes), &rawCfg); err != nil {
		return nil, err
	}
	return ValidateCfg(rawCfg)
}
