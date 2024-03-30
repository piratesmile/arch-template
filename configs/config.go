package configs

import (
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	Dev = "dev"
	Pro = "pro"
)

type Config struct {
	Database Database
	APP      APP
	Auth     Auth
	Log      Log
}

type APP struct {
	Env      string
	Host     string
	ServerID int64
}

type Database struct {
	Driver string
	Source string
}

type Log struct {
	Level    string
	ErrFile  string
	LogFile  string
	FileSize int
	FileAge  int
}

type Auth struct {
	SecretKey  []byte
	Expiration time.Duration
}

func Setup(file string) (*Config, error) {
	var c Config
	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err := viper.Unmarshal(&c, matchKey())
	return &c, err
}

func matchKey() viper.DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.MatchName = func(mapKey, fieldName string) bool {
			newKey := strings.Replace(mapKey, "-", "", -1)
			return strings.EqualFold(newKey, fieldName)
		}
	}
}
