package config

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	ErrEnvFileImport = errors.New("could not import env file")
)

// ConfBuilder is used to help construct a Config by importing ConfMap 's and config files
type ConfBuilder struct {
	viper *viper.Viper
	err   error
}

// NewConfBuilder constructs a new config mapper
func NewConfBuilder() *ConfBuilder {
	return &ConfBuilder{
		viper: viper.New(),
	}
}

// NewConfBuilderGlobal is the same as NewConfBuilder except
// it uses the global viper instance
func NewConfBuilderGlobal() *ConfBuilder {
	return &ConfBuilder{
		viper: viper.GetViper(),
	}
}

var _conf Config

func GetConfig() Config {
	if _conf.v == nil {
		c, err := NewConfBuilderGlobal().AutomaticEnv().Build()

		if err != nil {
			logrus.WithError(err).
				WithField("func", "config.GetConfig").
				Error("Failed to initialize configuration")
		}

		_conf = c
	}
	return _conf
}

// Build builds and returns the Config object
func (b *ConfBuilder) Build() (Config, error) {
	c := &Config{v: b.viper}
	if err := b.err; err != nil {
		return *c, err
	}
	if err := b.viper.Unmarshal(c); err != nil {
		return *c, err
	}
	return *c, nil
}

// AutomaticEnv has Viper check ENV variables for all.
// keys set in config, default & flags
func (b *ConfBuilder) AutomaticEnv() *ConfBuilder {
	b.viper.AutomaticEnv()
	return b
}

// ImportMap into the config
// this is useful for setting defaults
func (b *ConfBuilder) ImportMap(cm ConfMap) *ConfBuilder {
	for k, v := range cm {
		b.viper.SetDefault(k, v)
	}
	return b
}

// ImportFile into the config
func (b *ConfBuilder) ImportFile(file string) *ConfBuilder {
	b.viper.SetConfigFile(file)
	if err := b.viper.ReadInConfig(); err != nil {
		b.err = fmt.Errorf("%w: %s %s", ErrEnvFileImport, file, err.Error())
	}
	return b
}
