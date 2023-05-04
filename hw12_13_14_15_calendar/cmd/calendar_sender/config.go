package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf `mapstructure:"logger"`
	Queue  Queue      `mapstructure:"queue"`
}

type LoggerConf struct {
	Level string
}

type Queue struct {
	BrokerAddr string `mapstructure:"broker"`
	ReadTopic  string `mapstructure:"read_topic"`
	WriteTopic string `mapstructure:"write_topic"`
}

func NewConfig() (Config, error) {
	c := Config{}
	return c, c.init()
}

func (c *Config) init() error {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(c); err != nil {
		return err
	}
	return nil
}
