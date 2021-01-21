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
	configKey = "comet.toml"

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

// Config is comet config.
type Config struct {
	Debug     bool
	Env       *Env
	Discovery *Discovery
	TCP       *TCP
	Websocket *Websocket
	Protocol  *Protocol
	Bucket    *Bucket
	RPCClient *RPCClient
	RPCServer *RPCServer
	Whitelist *Whitelist
	Log       *log.Config
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
	Offline   bool
	Addrs     []string
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

// TCP is tcp config.
type TCP struct {
	Bind         []string
	Sndbuf       int
	Rcvbuf       int
	KeepAlive    bool
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

// Websocket is websocket config.
type Websocket struct {
	Bind        []string
	TLSOpen     bool
	TLSBind     []string
	CertFile    string
	PrivateFile string
}

// Protocol is protocol config.
type Protocol struct {
	Timer            int
	TimerSize        int
	SvrProto         int
	CliProto         int
	HandshakeTimeout xtime.Duration
}

// Bucket is bucket config.
type Bucket struct {
	Size          int
	Channel       int
	Room          int
	RoutineAmount uint64
	RoutineSize   int
}

// Whitelist is white list config.
type Whitelist struct {
	Whitelist []int64
	WhiteLog  string
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

func (e *Env) fix() error {
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
	return nil
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
func (t *TCP) fix() error {
	if t.Bind == nil {
		t.Bind = []string{":3101"}
	}
	if t.Sndbuf == 0 {
		t.Sndbuf = 4096
	}
	if t.Rcvbuf == 0 {
		t.Rcvbuf = 4096
	}
	if t.Reader == 0 {
		t.Reader = 32
	}
	if t.ReadBuf == 0 {
		t.ReadBuf = 32
	}
	if t.ReadBufSize == 0 {
		t.ReadBufSize = 1024
	}
	if t.Writer == 0 {
		t.Writer = 32
	}
	if t.WriteBuf == 0 {
		t.WriteBuf = 1024
	}
	if t.WriteBufSize == 0 {
		t.WriteBufSize = 8192
	}
	return nil
}
func (w *Websocket) fix() error {
	if w.Bind == nil {
		w.Bind = []string{":3101"}
	}
	return nil
}
func (p *Protocol) fix() error {
	if p.Timer == 0 {
		p.Timer = 32
	}
	if p.TimerSize == 0 {
		p.TimerSize = 2048
	}
	if p.CliProto == 0 {
		p.CliProto = 5
	}
	if p.SvrProto == 0 {
		p.SvrProto = 10
	}
	if p.HandshakeTimeout == 0 {
		p.HandshakeTimeout = xtime.Duration(time.Second * 5)
	}
	return nil
}
func (b *Bucket) fix() error {
	if b.Size == 0 {
		b.Size = 32
	}
	if b.Channel == 0 {
		b.Channel = 1024
	}
	if b.Room == 0 {
		b.Room = 1024
	}
	if b.RoutineAmount == 0 {
		b.RoutineAmount = 32
	}
	if b.RoutineSize == 0 {
		b.RoutineSize = 1024
	}
	return nil
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

	if c.TCP == nil {
		c.TCP = &TCP{
			Bind:         []string{":3101"},
			Sndbuf:       4096,
			Rcvbuf:       4096,
			KeepAlive:    false,
			Reader:       32,
			ReadBuf:      1024,
			ReadBufSize:  8192,
			Writer:       32,
			WriteBuf:     1024,
			WriteBufSize: 8192,
		}
	}
	if err = c.TCP.fix(); err != nil {
		return
	}

	if c.Websocket == nil {
		c.Websocket = &Websocket{
			Bind: []string{":3102"},
		}
	}
	if err = c.Websocket.fix(); err != nil {
		return
	}

	if c.Protocol == nil {
		c.Protocol = &Protocol{
			Timer:            32,
			TimerSize:        2048,
			CliProto:         5,
			SvrProto:         10,
			HandshakeTimeout: xtime.Duration(time.Second * 5),
		}
	}
	if err = c.Protocol.fix(); err != nil {
		return
	}

	if c.Bucket == nil {
		c.Bucket = &Bucket{
			Size:          32,
			Channel:       1024,
			Room:          1024,
			RoutineAmount: 32,
			RoutineSize:   1024,
		}
	}

	if err = c.Bucket.fix(); err != nil {
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
