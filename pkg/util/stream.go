package util

import (
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func Stream(url string, fn func(body io.Reader) error) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	logrus.Debugf("# %s %s %s", req.Method, req.URL.String(), req.Proto)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	status := res.Proto + ` ` + res.Status
	logrus.Debugf("# %s", status)
	if res.StatusCode/100 > 2 {
		return errors.New(status)
	}
	return fn(res.Body)
}
