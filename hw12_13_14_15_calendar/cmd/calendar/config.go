package main

import "github.com/spf13/viper"

// Config При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `yaml:"logger"`
	pg     Postgres   `yaml:"postgresql"`
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

type Postgres struct {
	Host     string
	User     string
	Password string
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
