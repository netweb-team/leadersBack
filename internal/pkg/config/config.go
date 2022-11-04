package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/gommon/log"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db,omitempty"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	Postgres  *DBConfig `json:"postgres"`
	Host      string    `json:"host"`
	Port      string    `json:"port"`
	Path      string    `json:"path"`
	ApiKey    string    `json:"api_token"`
	SearchApi string    `json:"search_api"`
	MapApi    string    `json:"map_api"`
	CoordApi  string    `json:"coord_api"`
}

var config = &Config{}

func New() *Config {
	return config
}

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Warn(err)
		return
	}
	json.Unmarshal(data, config)
}
