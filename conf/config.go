package conf

import (
	"errors"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HTTPAddr string   `toml:"http_addr"`
	Resource string   `toml:"static_resource"`
	DB       DBConfig `toml:"database"`
	ES       ESConfig `toml:"es"`
}

type DBConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DBName       string `toml:"db_name"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxConns     int    `toml:"max_conns"`
}

type ESConfig struct {
	Endpoint string `toml:"endpoint"`
}

func Load(confPath string) (*Config, error) {
	if confPath == "" {
		return nil, errors.New("config path is missing")
	}

	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	_, err = toml.Decode(string(data), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
