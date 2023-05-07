package main

import (
	"net"
	"strconv"

	"github.com/spf13/viper"
)

const (
	StorageInMemory = 1
	StorageDB       = 2
)

// Config При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf `mapstructure:"logger"`
	Pg         Postgres   `mapstructure:"db"`
	HttpServer Server     `mapstructure:"http"`
	GrpcServer Server     `mapstructure:"grpc"`
	Storage    Storage    `mapstructure:"storage"`
}

type LoggerConf struct {
	Level string
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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

func (c *Config) GetHTTPServerAddr() string {
	return net.JoinHostPort(c.HttpServer.Host, strconv.Itoa(c.HttpServer.Port))
}

func (c *Config) GetGrpcServerAddr() string {
	return net.JoinHostPort(c.GrpcServer.Host, strconv.Itoa(c.GrpcServer.Port))
}
