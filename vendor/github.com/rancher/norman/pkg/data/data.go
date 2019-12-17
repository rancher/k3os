package data

import (
	"github.com/rancher/norman/pkg/types/convert"
	"github.com/rancher/norman/pkg/types/values"
)

type List []map[string]interface{}

type Object map[string]interface{}

func New() Object {
	return map[string]interface{}{}
}

func (o Object) Map(names ...string) Object {
	v := values.GetValueN(o, names...)
	m := convert.ToMapInterface(v)
	return m
}

func (o Object) Slice(names ...string) (result []Object) {
	v := values.GetValueN(o, names...)
	for _, item := range convert.ToInterfaceSlice(v) {
		result = append(result, convert.ToMapInterface(item))
	}
	return
}

func (o Object) Values() (result []Object) {
	for k := range o {
		result = append(result, o.Map(k))
	}
	return
}

func (o Object) String(names ...string) string {
	v := values.GetValueN(o, names...)
	return convert.ToString(v)
}

func (o Object) StringSlice(names ...string) []string {
	v := values.GetValueN(o, names...)
	return convert.ToStringSlice(v)
}

func (o Object) Set(key string, obj interface{}) {
	if o == nil {
		return
	}
	o[key] = obj
}

func (o Object) SetNested(obj interface{}, key ...string) {
	values.PutValue(o, obj, key...)
}

func (o Object) Bool(key ...string) bool {
	return convert.ToBool(values.GetValueN(o, key...))
}
