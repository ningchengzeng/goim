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
	configKey = "logic.toml"
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

// Config config.
type Config struct {
	Env        *Env
	Discovery  *Discovery
	RPCClient  *RPCClient
	RPCServer  *RPCServer
	HTTPServer *HTTPServer
	Nsq        *Nsq
	Redis      *Redis
	Node       *Node
	Backoff    *Backoff
	Regions    map[string][]string
}

// Discovery is discovery config.
type Discovery struct {
	Nodes []string
}

// Env is env config.
type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
	Weight    int64
}

// Node node config.
type Node struct {
	DefaultDomain string
	HostDomain    string
	TCPPort       int
	WSPort        int
	WSSPort       int
	HeartbeatMax  int
	Heartbeat     xtime.Duration
	RegionWeight  float64
}

// Backoff backoff.
type Backoff struct {
	MaxDelay  int32
	BaseDelay int32
	Factor    float32
	Jitter    float32
}

// Redis .
type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
	Expire       xtime.Duration
}

// Nsq .
type Nsq struct {
	Topic   string
	Address string
}

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    xtime.Duration
	Timeout xtime.Duration
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string
	Addr              string
	Timeout           xtime.Duration
	IdleTimeout       xtime.Duration
	MaxLifeTime       xtime.Duration
	ForceCloseWait    xtime.Duration
	KeepAliveInterval xtime.Duration
	KeepAliveTimeout  xtime.Duration
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
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

func (r *RPCClient) fix() error {
	if r.Dial == 0 {
		r.Dial = xtime.Duration(time.Second)
	}
	if r.Timeout == 0 {
		r.Timeout = xtime.Duration(time.Second)
	}
	return nil
}
func (r *RPCServer) fix() error {
	if r.Network == "" {
		r.Network = "tcp"
	}
	if r.Addr == "" {
		r.Addr = ":3109"
	}
	if r.Timeout == 0 {
		r.Timeout = xtime.Duration(time.Second)
	}
	if r.IdleTimeout == 0 {
		r.IdleTimeout = xtime.Duration(time.Second * 60)
	}
	if r.MaxLifeTime == 0 {
		r.MaxLifeTime = xtime.Duration(time.Hour * 2)
	}
	if r.ForceCloseWait == 0 {
		r.ForceCloseWait = xtime.Duration(time.Second * 20)
	}
	if r.KeepAliveInterval == 0 {
		r.KeepAliveInterval = xtime.Duration(time.Second * 60)
	}
	if r.KeepAliveTimeout == 0 {
		r.KeepAliveTimeout = xtime.Duration(time.Second * 20)
	}
	return nil
}
func (h *HTTPServer) fix() (err error) {

	if h.Network == "" {
		h.Network = "tcp"
	}
	if h.Addr == "" {
		h.Addr = "3111"
	}
	if h.ReadTimeout == 0 {
		h.ReadTimeout = xtime.Duration(time.Second)
	}
	if h.WriteTimeout == 0 {
		h.WriteTimeout = xtime.Duration(time.Second)
	}
	return
}
func (b *Backoff) fix() (err error) {
	if b.MaxDelay == 0 {
		b.MaxDelay = 300
	}
	if b.BaseDelay == 0 {
		b.BaseDelay = 3
	}
	if b.Factor == 0 {
		b.Factor = 1.8
	}
	if b.Jitter == 0 {
		b.Jitter = 1.3
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

	if c.RPCClient == nil {
		c.RPCClient = &RPCClient{
			Dial:    xtime.Duration(time.Second),
			Timeout: xtime.Duration(time.Second),
		}
	}
	if err = c.RPCClient.fix(); err != nil {
		return
	}

	if c.RPCServer == nil {
		c.RPCServer = &RPCServer{
			Network:           "tcp",
			Addr:              ":3109",
			Timeout:           xtime.Duration(time.Second),
			IdleTimeout:       xtime.Duration(time.Second * 60),
			MaxLifeTime:       xtime.Duration(time.Hour * 2),
			ForceCloseWait:    xtime.Duration(time.Second * 20),
			KeepAliveInterval: xtime.Duration(time.Second * 60),
			KeepAliveTimeout:  xtime.Duration(time.Second * 20),
		}
	}
	if err = c.RPCServer.fix(); err != nil {
		return
	}

	if c.Backoff == nil {
		c.Backoff = &Backoff{MaxDelay: 300, BaseDelay: 3, Factor: 1.8, Jitter: 1.3}
	}
	if err = c.Backoff.fix(); err != nil {
		return
	}

	if c.HTTPServer == nil {
		c.HTTPServer = &HTTPServer{
			Network:      "tcp",
			Addr:         "3111",
			ReadTimeout:  xtime.Duration(time.Second),
			WriteTimeout: xtime.Duration(time.Second),
		}
	}
	if err = c.HTTPServer.fix(); err != nil {
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
