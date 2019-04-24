package config

import (
	"strings"

	"github.com/rancher/mapper"
	"github.com/rancher/mapper/convert"
	"github.com/rancher/mapper/mappers"
)

type FuzzyNames struct {
	mappers.DefaultMapper
	names map[string]string
}

func (f *FuzzyNames) ToInternal(data map[string]interface{}) error {
	for k, v := range data {
		if newK, ok := f.names[k]; ok && newK != k {
			data[newK] = v
		}
	}
	return nil
}

func (f *FuzzyNames) addName(name string) {
	f.names[strings.ToLower(name)] = name
	f.names[convert.ToYAMLKey(name)] = name
	f.names[strings.ToLower(convert.ToYAMLKey(name))] = name
}

func (f *FuzzyNames) ModifySchema(schema *mapper.Schema, schemas *mapper.Schemas) error {
	f.names = map[string]string{}

	for name := range schema.ResourceFields {
		if strings.HasSuffix(name, "s") && len(name) > 1 {
			f.addName(name[:len(name)-1])
		}
		if strings.HasSuffix(name, "es") && len(name) > 2 {
			f.addName(name[:len(name)-2])
		}
		f.addName(name)
	}

	return nil
}
