package control

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rancher/k3os/config"
	"github.com/rancher/k3os/pkg/util"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func configCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "get",
			Usage:  "get a value",
			Action: get,
		},
		{
			Name:   "set",
			Usage:  "set a value",
			Action: set,
		},
		{
			Name:  "export",
			Usage: "export configuration",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Usage: "file to which to save",
				},
				cli.BoolFlag{
					Name:  "full, f",
					Usage: "export full configuration, including internal and default settings",
				},
			},
			Action: export,
		},
		{
			Name:     "generate",
			Usage:    "generate a configuration file from a template",
			Action:   generate,
			HideHelp: true,
		},
		{
			Name:   "merge",
			Usage:  "merge configuration from stdin",
			Action: merge,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "file from which to read",
				},
			},
		},
		{
			Name:   "validate",
			Usage:  "validate configuration from stdin",
			Action: validate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "file from which to read",
				},
			},
		},
	}
}

func set(c *cli.Context) error {
	if c.NArg() < 2 {
		return nil
	}
	key := c.Args().Get(0)
	value := c.Args().Get(1)
	if key == "" {
		return nil
	}
	err := config.Set(key, value)
	if err != nil {
		logrus.Fatal(err)
	}
	return nil
}

func get(c *cli.Context) error {
	arg := c.Args().Get(0)
	if arg == "" {
		return nil
	}
	val, err := config.Get(arg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"key": arg, "val": val, "err": err}).Fatal("failed to get value")
	}
	printYaml := false
	switch val.(type) {
	case []interface{}:
		printYaml = true
	case map[interface{}]interface{}:
		printYaml = true
	}
	if printYaml {
		bytes, err := yaml.Marshal(val)
		if err != nil {
			logrus.Fatal(err)
		}
		fmt.Println(string(bytes))
	} else {
		fmt.Println(val)
	}
	return nil
}

func generate(c *cli.Context) error {
	if err := util.GenTemplate(os.Stdin, os.Stdout); err != nil {
		logrus.Fatalf("failed to generate config, err: '%s'", err)
	}
	return nil
}

func export(c *cli.Context) error {
	content, err := config.Export("", c.Bool("full"))
	if err != nil {
		logrus.Fatal(err)
	}
	output := c.String("output")
	if output == "" {
		fmt.Println(content)
	} else {
		err := util.WriteFileAtomic(output, []byte(content), 0400)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	return nil
}

func merge(c *cli.Context) error {
	bytes, err := inputBytes(c)
	if err != nil {
		logrus.Fatal(err)
	}
	if err = config.Merge(bytes); err != nil {
		logrus.Error(err)
		validationErrors, err := config.ValidateBytes(bytes)
		if err != nil {
			logrus.Fatal(err)
		}
		for _, validationError := range validationErrors.Errors() {
			logrus.Error(validationError)
		}
	}
	return nil
}

func validate(c *cli.Context) error {
	bytes, err := inputBytes(c)
	if err != nil {
		logrus.Fatal(err)
	}
	validationErrors, err := config.ValidateBytes(bytes)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, validationError := range validationErrors.Errors() {
		logrus.Error(validationError)
	}
	return nil
}

func inputBytes(c *cli.Context) ([]byte, error) {
	input := os.Stdin
	inputFile := c.String("input")
	if inputFile != "" {
		var err error
		input, err = os.Open(inputFile)
		if err != nil {
			return nil, err
		}
		defer input.Close()
	}
	content, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(content, []byte{13, 10}) {
		return nil, errors.New("file format shouldn't contain CRLF characters")
	}
	return content, nil
}
