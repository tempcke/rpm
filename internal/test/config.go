package test

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/pkg/log"
)

var _conf config.Config

func GetConfig() config.Config {
	if _conf.PostgresHost() == "" {
		var err error

		cb := config.NewConfBuilderGlobal().AutomaticEnv()
		if envFile := findConfigFile(); envFile != "" {
			cb.ImportFile(envFile)
		}

		conf, err := cb.Build()
		if err != nil {
			log.WithError(err).
				WithField("func", "test.getConfig").
				Fatal("could not build configuration")
		}
		_conf = conf
	}
	return _conf
}

func findConfigFile() string {
	var file = ".env"
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(file); err == nil {
			_ = godotenv.Load(file)
			return file
		}
		file = fmt.Sprintf("../%s", file)
	}
	return ""
}
