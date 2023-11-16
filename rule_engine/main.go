package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"rule_engine/config"
	"rule_engine/global"
	"rule_engine/log"
	"rule_engine/model"
	"rule_engine/mysql"
	"rule_engine/redis"
	"rule_engine/router"
	"rule_engine/server"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		cfg string
	)
	flag.StringVar(&cfg, "config", "", "server config [toml]")
	flag.Parse()
	if len(cfg) == 0 {
		fmt.Println("config is empty")
		os.Exit(0)
	}
	config.Init(cfg)
	conf := config.CoreConf
	log.Init(&conf.Log)
	gin.SetMode(conf.Server.Mode)
	if err := mysql.InitMysql(&conf.Mysql); err != nil {
		log.Fatal("Init mysql error:%s", zap.Error(err))
	}
	if err := redis.Init(&conf.Redis); err != nil {
		log.Fatal("Init redis error:%s", zap.Error(err))
	}
	heartbeat() //给resty_redis subscribe提供心跳避免超时
	if err := server.Run(router.NewHttpRouter()); nil != err {
		log.Error("server run error", zap.Error(err))
	}
}

func heartbeat() {
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				model.Message{Event: global.Heartbeat}.Update()
			case <-global.Ctx.Done():
				return
			}
		}
	}()
}
