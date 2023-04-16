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
	pg         Postgres   `mapstructure:"db"`
	httpServer Server     `mapstructure:"http"`
	grpcServer Server     `mapstructure:"grpc"`
	storage    Storage
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
	return "postgres://" + c.pg.User + ":" + c.pg.Password + "@" + c.pg.Host
}

func (c *Config) GetHttpServerAddr() string {
	return net.JoinHostPort(c.httpServer.Host, strconv.Itoa(c.httpServer.Port))
}

func (c *Config) GetGrpcServerAddr() string {
	return net.JoinHostPort(c.grpcServer.Host, strconv.Itoa(c.grpcServer.Port))
}

func (c *Config) GetStorageType() int {
	return c.storage.Type
}
