package main

import (
	"github.com/spf13/viper"
)

const (
	StorageInMemory = 1
	StorageDB       = 2
)

type Config struct {
	Notify  Time       `mapstructure:"notify"`
	Logger  LoggerConf `mapstructure:"logger"`
	Pg      Postgres   `mapstructure:"db"`
	Queue   Queue      `mapstructure:"queue"`
	Storage Storage    `mapstructure:"storage"`
}

type Time struct {
	Hour int `mapstructure:"hour"`
	Min  int `mapstructure:"min"`
}

type LoggerConf struct {
	Level string
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Queue struct {
	BrokerAddr string `mapstructure:"broker"`
	Topic      string `mapstructure:"topic"`
}

type Storage struct {
	Type int `mapstructure:"type"`
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

func (c *Config) GetDBURL() string {
	return "postgres://" + c.Pg.User + ":" + c.Pg.Password + "@" + c.Pg.Host
}

func (c *Config) GetStorageType() int {
	return c.Storage.Type
}

func (c *Config) GetBroker() string {
	return c.Queue.BrokerAddr
}

func (c *Config) GetTopic() string {
	return c.Queue.Topic
}

func (c *Config) GetNotifyTime() (int, int) {
	return c.Notify.Hour, c.Notify.Min
}
