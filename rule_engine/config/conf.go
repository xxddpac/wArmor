package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"rule_engine/log"
	"rule_engine/mysql"
	"rule_engine/redis"
)

var (
	CoreConf *Settings
)

func Init(conf string) {
	_, err := toml.DecodeFile(conf, &CoreConf)
	if err != nil {
		fmt.Printf("Err %v", err)
		os.Exit(1)
	}
}

type Settings struct {
	Log    log.Config
	Mysql  mysql.Config
	Redis  redis.Config
	Server Server
}

type Server struct {
	Port int
	Mode string
}
