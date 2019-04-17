package config

import (
	"errors"
	"strings"

	"github.com/rancher/k3os/pkg/util"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

const (
	validatePrompt = "error configuration, check configuration please use `k3os c validate`"
)

func Get(key string) (interface{}, error) {
	cfg := LoadConfig("", false)
	d := map[interface{}]interface{}{}
	if err := util.Convert(cfg, &d); err != nil {
		return nil, err
	}
	v, _ := util.GetValue(key, d)
	return v, nil
}

func Set(key string, value interface{}) error {
	exist, errs := ReadConfigFiles(nil, CloudConfigFile)
	if len(errs) > 0 {
		return errors.New(validatePrompt)
	}
	_, modified := util.SetValue(key, exist, value)
	cfg := &CloudConfig{}
	if err := util.Convert(modified, cfg); err != nil {
		return err
	}
	return util.WriteToFile(modified, CloudConfigFile)
}

func Export(prefix string, full bool) (string, error) {
	raw, errs := LoadRawConfig(prefix, full)
	if len(errs) > 0 {
		return "", errors.New(validatePrompt)
	}
	raw = filterAdditional(raw)
	bytes, err := yaml.Marshal(raw)
	return string(bytes), err
}

func Merge(bytes []byte) error {
	data, errs := ReadConfigFiles(bytes)
	if len(errs) > 0 {
		return errors.New(validatePrompt)
	}
	exist, errs := ReadConfigFiles(nil, CloudConfigFile)
	if len(errs) > 0 {
		return errors.New(validatePrompt)
	}
	err := mergo.Merge(&exist, data)
	if err != nil {
		return err
	}
	return util.WriteToFile(exist, CloudConfigFile)
}

func filterAdditional(data map[interface{}]interface{}) map[interface{}]interface{} {
	for _, additional := range Additional {
		_, data = util.FilterKeys(data, strings.Split(additional, "."))
	}
	return data
}
