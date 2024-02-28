package test

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/tempcke/rpm/internal/configs"
)

var (
	_buildConfigOnce sync.Once
	_conf            configs.Config
)

func Config() configs.Config {
	_buildConfigOnce.Do(func() {
		if file := findConfigFile(); file != "" {
			m, err := godotenv.Read(findConfigFile())
			if err != nil {
				panic(err) // this function is only called from tests, so it is okay to panic
			}
			_conf = configs.New(configs.WithEnvFromMap(m))
		}
	})
	return _conf
}

func findConfigFile() string {
	var file = ".env"
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(file); err == nil {
			return file
		}
		file = fmt.Sprintf("../%s", file)
	}
	return ""
}
