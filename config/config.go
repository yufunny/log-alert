package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type FileConfig struct {
	File     string `yaml:"file"`
	Desc     string `yaml:"desc"`
	Bound 	 string `yaml:"bound"`
	Rules    []RuleConfig `yaml:"rules"`
}

type RuleConfig struct {
	Rule     string `yaml:"rule"`
	Desc     string `yaml:"desc"`
	Duration string `yaml:"duration"`
	Times    int    `yaml:"times"`
	Interval string `yaml:"interval"`
	Receiver []string     `yaml:"receivers"`
}

type NotifyConfig struct {
	Driver string `yaml:"driver"`
	Url    string `yaml:"url"`
}

type SystemConfig struct {
	Mode     string       `yaml:"mode"`
	Receiver []string     `yaml:"receivers"`
	Notify   NotifyConfig `yaml:"notify"`
	Files    []FileConfig `yaml:"files"`
}

// LoadConfig  加载系统配置
func LoadConfig(file string) (*SystemConfig, error) {
	b, e := ioutil.ReadFile(file)
	if nil != e {
		return nil, errors.New("Config->Read config file[" + file + "] error; " + e.Error())
	}
	config := &SystemConfig{}
	e = yaml.Unmarshal(b, config)
	if nil != e {
		return nil, errors.New("Config->Unmarshal config from config file[" + file + "] error; " + e.Error())
	}
	return config, nil
}
