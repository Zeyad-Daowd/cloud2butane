package cloud2butane

import (
	"strconv"

	"gopkg.in/yaml.v3"
)

type Butane struct {
	Variant string        `yaml:"variant"`
	Version string        `yaml:"version"`
	Passwd  ButanePasswd  `yaml:"passwd"`
	Storage ButaneStorage `yaml:"storage"`
}

type ButanePasswd struct {
	Users []ButaneUser `yaml:"users"`
}

type ButaneStorage struct {
	Files []ButaneFile `yaml:"files"`
}
type ButaneUser struct {
	Name   string   `yaml:"name"`
	Groups []string `yaml:"groups"`
	Shell  string   `yaml:"shell"`
}

type Octal int

func (o Octal) MarshalYAML() (interface{}, error) {
	octalStr := "0" + strconv.FormatInt(int64(o), 8)

	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!int", // Tells YAML it is a number, not a string
		Value: octalStr,
	}, nil
}

type ButaneFile struct {
	Path     string        `yaml:"path"`
	Contents ButaneContent `yaml:"contents"`
	// mode should be displayed as an octal number in the butane file, but it is stored as an int in the struct
	Mode Octal `yaml:"mode"`
}
type ButaneContent struct {
	Inline string `yaml:"inline"`
}
