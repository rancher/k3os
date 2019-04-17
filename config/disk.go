package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/rancher/k3os/pkg/util"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func LoadConfig(prefix string, terminate bool) *CloudConfig {
	cfg := &CloudConfig{}
	raw, errs := LoadRawConfig(prefix, true)
	if terminate && len(errs) > 0 {
		for i, e := range errs {
			logrus.Errorf("failed to load configuration, error(#%d): %s", i, e.Error())
			return cfg
		}
	}
	if err := util.Convert(raw, cfg); err != nil {
		v, vErr := ValidateCfg(cfg)
		if vErr != nil {
			logrus.Fatal(vErr)
		}
		for i, e := range v.Errors() {
			logrus.Errorf("failed to validate configuration, error(#%d): %s", i, e.String())
		}
		return cfg
	}
	return cfg
}

func LoadRawConfig(prefix string, full bool) (map[interface{}]interface{}, []error) {
	var raw map[interface{}]interface{}
	errs := make([]error, 0)
	if full {
		raw, errs = ReadConfigFiles(nil, path.Join(prefix, OsConfigFile))
	}
	files := configDirFiles(prefix)
	files = append(files, path.Join(prefix, CloudConfigFile))
	additional, readErrs := ReadConfigFiles(nil, files...)
	if readErrs != nil && len(readErrs) > 0 {
		for _, err := range readErrs {
			errs = append(errs, err)
		}
	}
	// TODO: cmdline, extra-cmdline, debug, metadata need to be implement
	if err := mergo.Merge(&raw, additional, mergo.WithOverride); err != nil {
		errs = append(errs, err)
	}
	return raw, errs
}

func ReadConfigFiles(bytes []byte, files ...string) (map[interface{}]interface{}, []error) {
	// TODO: metadata, user-data, cloud-init need to be implement
	left := make(map[interface{}]interface{})
	errs := make([]error, 0)
	for _, file := range files {
		content, err := ReadConfigFile(file)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(content) == 0 {
			continue
		}
		right := make(map[interface{}]interface{})
		err = yaml.Unmarshal(content, &right)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		cfg := &CloudConfig{}
		if err := util.Convert(right, cfg); err != nil {
			errs = append(errs, err)
			continue
		}
		if err := mergo.Merge(&left, right, mergo.WithOverride); err != nil {
			errs = append(errs, err)
			continue
		}
	}
	if bytes == nil || len(bytes) == 0 {
		return left, errs
	}
	right := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(bytes, &right); err != nil {
		errs = append(errs, err)
		return left, errs
	}
	cfg := &CloudConfig{}
	if err := util.Convert(right, cfg); err != nil {
		errs = append(errs, err)
		return left, errs
	}
	if err := mergo.Merge(&left, right, mergo.WithOverride); err != nil {
		errs = append(errs, err)
	}
	return left, errs
}

func ReadConfigFile(file string) ([]byte, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			content = []byte{}
		} else {
			return nil, err
		}
	}
	return content, err
}

func configDirFiles(prefix string) []string {
	dir := path.Join(prefix, CloudConfigDir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Debugf("directory %s does not exist", CloudConfigDir)
		} else {
			logrus.Errorf("failed to read %s: %v", CloudConfigDir, err)
		}
		return []string{}
	}
	var final []string
	for _, file := range files {
		if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			final = append(final, path.Join(dir, file.Name()))
		}
	}
	return final
}
