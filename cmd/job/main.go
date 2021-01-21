package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bilibili/discovery/naming"
	log "github.com/go-kratos/kratos/pkg/log"
	"github.com/ningchengzeng/goim/internal/job"
	"github.com/ningchengzeng/goim/internal/job/conf"

	resolver "github.com/bilibili/discovery/naming/grpc"
)

var (
	ver = "2.0.0"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)

	log.Info("goim-job [version: %s env: %+v] start", ver, conf.Conf.Env)
	// grpc register naming
	dis := naming.New(conf.Conf.DiscoveryConfig())
	resolver.Register(dis)

	// job
	j := job.New(conf.Conf)
	consumer, _ := job.NewNsqConsumer(conf.Conf.Nsq, j)
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("goim-job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			j.Close()
			consumer.Close()
			log.Info("goim-job [version: %s] exit", ver)
			log.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
