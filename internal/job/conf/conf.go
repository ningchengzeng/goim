package conf

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bilibili/discovery/naming"
	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	log "github.com/go-kratos/kratos/pkg/log"
	xtime "github.com/ningchengzeng/goim/pkg/time"
)

var (
	configKey = "job.toml"
	// Conf config
	Conf = &Config{}
)

// Init init config.
func Init() (err error) {
	if err = paladin.Init(); err != nil {
		return
	}
	return paladin.Watch(configKey, Conf)
}

// Config is job config.
type Config struct {
	Discovery *Discovery
	Env       *Env
	Nsq       *Nsq
	Comet     *Comet
	Room      *Room
	Log       *log.Config
}

// Discovery is discovery config.
type Discovery struct {
	Nodes []string
}

// Room is room config.
type Room struct {
	Batch  int
	Signal xtime.Duration
	Idle   xtime.Duration
}

// Comet is comet config.
type Comet struct {
	RoutineChan int
	RoutineSize int
}

// Nsq is kafka config.
type Nsq struct {
	Topic   string
	Channel string
	Address []string
}

// Env is env config.
type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
}

// DiscoveryConfig 创建发现服务配置
func (c *Config) DiscoveryConfig() *naming.Config {
	return &naming.Config{
		Nodes:  c.Discovery.Nodes,
		Zone:   c.Env.Zone,
		Region: c.Env.Region,
		Env:    c.Env.DeployEnv,
		Host:   c.Env.Host,
	}
}

func (e *Env) fix() (err error) {
	if e.Region == "" {
		e.Region = env.Region
	}
	if e.Zone == "" {
		e.Zone = env.Zone
	}
	if e.Host == "" {
		e.Host = env.Hostname
	}
	if e.DeployEnv == "" {
		e.DeployEnv = env.DeployEnv
	}

	return
}

func (c *Comet) fix() (err error) {
	if c.RoutineChan == 0 {
		c.RoutineChan = 1024
	}
	if c.RoutineSize == 0 {
		c.RoutineSize = 32
	}
	return
}

func (r *Room) fix() (err error) {
	if r.Batch == 0 {
		r.Batch = 20
	}
	if r.Signal == 0 {
		r.Signal = xtime.Duration(time.Second)
	}
	if r.Idle == 0 {
		r.Idle = xtime.Duration(time.Minute * 15)
	}
	return
}

func (c *Config) fix() (err error) {
	if c.Env == nil {
		c.Env = new(Env)
	}
	if err = c.Env.fix(); err != nil {
		return
	}

	if c.Comet == nil {
		c.Comet = &Comet{RoutineChan: 1024, RoutineSize: 32}
	}
	if err = c.Comet.fix(); err != nil {
		return
	}

	if c.Room == nil {
		c.Room = &Room{
			Batch:  20,
			Signal: xtime.Duration(time.Second),
			Idle:   xtime.Duration(time.Minute * 15),
		}
	}
	if err = c.Room.fix(); err != nil {
		return
	}

	return
}

// Set config setter.
func (c *Config) Set(content string) (err error) {
	var tmpConf *Config
	if _, err = toml.Decode(content, &tmpConf); err != nil {
		log.Error("decode config fail %v", err)
		return
	}
	if err = tmpConf.fix(); err != nil {
		return
	}
	*Conf = *tmpConf
	return nil
}
