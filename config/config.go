package config

import (
	"github.com/BurntSushi/toml"
	"time"
)

type config struct {
	IP        string        `toml:"server_ip"`
	Port      string        `toml:"server_port"`
	ValidThru time.Duration `toml:"cache_valid_period_sec"`
	Servers   []string      `toml:"servers"`
}

var Get config

func ReadConfig() {
	if _, err := toml.DecodeFile("config/app.toml", &Get); err != nil {
		panic(err)
	}
}
