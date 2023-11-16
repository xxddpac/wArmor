package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"rule_engine/config"
	"rule_engine/global"
	"rule_engine/log"
	"syscall"
	"time"
)

func Run(router *gin.Engine) error {
	var (
		addr string
		conf = config.CoreConf
	)
	addr = fmt.Sprintf(":%d", conf.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	log.Info(fmt.Sprintf("start http server, listen:%s", addr))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error", zap.Error(err))
		}
	}()
	tryDisConn(srv)
	return nil
}

func tryDisConn(srv *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT,
	)
	select {
	case sig := <-signals:
		go func() {
			select {
			case <-time.After(time.Second * 10):
				log.Warn("Shutdown gracefully timeout, application will shutdown immediately.")
				os.Exit(0)
			}
		}()
		log.Info(fmt.Sprintf("get signal %s, application will shutdown.", sig))
		log.Debug("stop HttpServer...")
		global.Cancel()
		_ = srv.Shutdown(context.Background())
		os.Exit(0)
	}
}
