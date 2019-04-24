package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rancher/mapper"
	"github.com/rancher/mapper/convert"
	merge2 "github.com/rancher/mapper/convert/merge"
	"github.com/rancher/mapper/values"
	"gopkg.in/yaml.v3"
)

const (
	systemConfig = "/k3os/system/config.yaml"
	localConfig  = "/var/lib/rancher/k3os/config.yaml"
	localConfigs = "/var/lib/rancher/k3os/config.d"
)

var (
	schemas = mapper.NewSchemas().Init(func(s *mapper.Schemas) *mapper.Schemas {
		s.DefaultMappers = func() []mapper.Mapper {
			return []mapper.Mapper{
				NewToSlice(),
				NewToBool(),
				&FuzzyNames{},
			}
		}
		return s
	}).MustImport(CloudConfig{})
	schema  = schemas.Schema("cloudConfig")
	readers = []reader{
		readSystemConfig,
		readCmdline,
		readLocalConfig,
	}
)

func ToEnv(cfg CloudConfig) ([]string, error) {
	data, err := convert.EncodeToMap(&cfg)
	if err != nil {
		return nil, err
	}

	return mapToEnv("", data), nil
}

func mapToEnv(prefix string, data map[string]interface{}) []string {
	var result []string
	for k, v := range data {
		keyName := strings.ToUpper(prefix + convert.ToYAMLKey(k))
		if data, ok := v.(map[string]interface{}); ok {
			subResult := mapToEnv(keyName+"_", data)
			result = append(result, subResult...)
		} else {
			result = append(result, fmt.Sprintf("%s=%v", keyName, v))
		}
	}
	return result
}

func ReadConfig() (CloudConfig, error) {
	var result CloudConfig

	data, err := merge(append(readers, readLocalConfigs()...)...)
	if err != nil {
		return result, err
	}

	return result, convert.ToObj(data, &result)
}

type reader func() (map[string]interface{}, error)

func merge(readers ...reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	for _, r := range readers {
		newData, err := r()
		if err != nil {
			return nil, err
		}

		data = merge2.UpdateMerge(schema, schemas, data, newData, false)
	}
	return data, nil
}

func readSystemConfig() (map[string]interface{}, error) {
	return readFile(systemConfig)
}

func readLocalConfig() (map[string]interface{}, error) {
	return readFile(localConfig)
}

func readLocalConfigs() []reader {
	var result []reader

	files, err := ioutil.ReadDir(localConfigs)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return []reader{
			func() (map[string]interface{}, error) {
				return nil, err
			},
		}
	}

	for _, f := range files {
		p := filepath.Join(localConfigs, f.Name())
		result = append(result, func() (map[string]interface{}, error) {
			return readFile(p)
		})
	}

	return result
}

func readFile(path string) (map[string]interface{}, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer f.Close()

	data := map[string]interface{}{}
	if err := yaml.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}

	return data, schema.Mapper.ToInternal(data)
}

func readCmdline() (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile("/proc/cmdline")
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	for _, item := range strings.Fields(string(bytes)) {
		parts := strings.SplitN(item, "=", 2)
		value := "true"
		if len(parts) > 1 {
			value = parts[1]
		}
		keys := strings.Split(parts[0], ".")
		existing, ok := values.GetValue(data, keys...)
		if ok {
			switch v := existing.(type) {
			case string:
				values.PutValue(data, []string{v, value}, keys...)
			case []string:
				values.PutValue(data, append(v, value), keys...)
			}
		} else {
			values.PutValue(data, value, keys...)
		}
	}

	return data, schema.Mapper.ToInternal(data)
}
