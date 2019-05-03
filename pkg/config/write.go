package config

import (
	"github.com/ghodss/yaml"
	"github.com/rancher/mapper/convert"
)

func PrintInstall(cfg CloudConfig) ([]byte, error) {
	data := map[string]interface{}{}
	data, err := convert.EncodeToMap(cfg.K3OS.Install)
	if err != nil {
		return nil, err
	}

	toYAMLKeys(data)
	return yaml.Marshal(data)
}

func ToBytes(cfg CloudConfig) ([]byte, error) {
	cfg.K3OS.Install = nil
	data := map[string]interface{}{}
	data, err := convert.EncodeToMap(cfg)
	if err != nil {
		return nil, err
	}

	toYAMLKeys(data)
	return yaml.Marshal(data)
}

func toYAMLKeys(data map[string]interface{}) {
	for k, v := range data {
		if sub, ok := v.(map[string]interface{}); ok {
			toYAMLKeys(sub)
		}
		newK := convert.ToYAMLKey(k)
		if newK != k {
			delete(data, k)
			data[newK] = v
		}
	}
}
