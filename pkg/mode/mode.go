package mode

import (
	"io/ioutil"
	"os"
	"strings"
)

func Get() (string, error) {
	bytes, err := ioutil.ReadFile("/run/k3os/mode")
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}
