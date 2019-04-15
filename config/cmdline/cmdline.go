package cmdline

import (
	"io/ioutil"
	"strings"

	"github.com/niusmallnan/k3os/pkg/util"
)

func GetCmdLine(key string) interface{} {
	parse := true
	if strings.HasPrefix(key, "k3os.") {
		parse = false
	}
	cmdline, _ := ReadCmdLine(parse)
	v, _ := util.GetValue(key, cmdline)
	return v
}

func ParseCmdLine(cmdLine string, parse bool) map[interface{}]interface{} {
	result := map[interface{}]interface{}{}
	for _, part := range strings.Split(cmdLine, " ") {
		if !strings.HasPrefix(part, "k3os.") && !parse {
			continue
		}
		var value string
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 1 {
			value = "true"
		} else {
			value = kv[1]
		}
		current := result
		keys := strings.Split(kv[0], ".")
		for i, key := range keys {
			if i == len(keys)-1 {
				current[key] = util.UnmarshalValue(value)
			} else {
				if obj, ok := current[key]; ok {
					if temp, ok := obj.(map[interface{}]interface{}); ok {
						current = temp
					} else {
						break
					}
				} else {
					temp := make(map[interface{}]interface{})
					current[key] = temp
					current = temp
				}
			}
		}
	}
	return result
}

func ReadCmdLine(parse bool) (m map[interface{}]interface{}, err error) {
	bytes, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, nil
	}
	return ParseCmdLine(strings.TrimSpace(util.UnescapeKernelParams(string(bytes))), parse), nil
}
