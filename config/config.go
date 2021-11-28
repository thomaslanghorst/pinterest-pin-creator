package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AccessTokenPath  string `yaml:"access_token_path"`
	ScheduleFilePath string `yaml:"schedule_file_path"`
	BrowserPath      string `yaml:"browser_path"`
	RedirectPort     int    `yaml:"redirect_port"`
}

type ConfigReaderInterface interface {
	Read() (*Config, error)
}

type ConfigReader struct {
	configFilePath string
}

func NewReader(configFilePath string) *ConfigReader {
	return &ConfigReader{
		configFilePath: configFilePath,
	}
}

func (r *ConfigReader) Read() (*Config, error) {

	yamlFile, err := ioutil.ReadFile(r.configFilePath)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
