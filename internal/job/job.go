package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/bilibili/discovery/naming"
	"github.com/ningchengzeng/goim/internal/job/conf"

	log "github.com/go-kratos/kratos/pkg/log"
)

// Job is push job.
type Job struct {
	c            *conf.Config
	cometServers map[string]*Comet

	rooms      map[string]*Room
	roomsMutex sync.RWMutex
}

// New new a push job.
func New(c *conf.Config) *Job {
	j := &Job{
		c:     c,
		rooms: make(map[string]*Room),
	}
	j.watchComet(c.DiscoveryConfig())
	return j
}

// Close close resounces.
func (j *Job) Close() error {
	return nil
}

func (j *Job) watchComet(c *naming.Config) {
	dis := naming.New(c)
	resolver := dis.Build("goim.comet")
	event := resolver.Watch()
	select {
	case _, ok := <-event:
		if !ok {
			panic("watchComet init failed")
		}
		if ins, ok := resolver.Fetch(); ok {
			if err := j.newAddress(ins.Instances); err != nil {
				panic(err)
			}
			log.Info("watchComet init newAddress:%+v", ins)
		}
	case <-time.After(10 * time.Second):
		log.Error("watchComet init instances timeout")
	}
	go func() {
		for {
			if _, ok := <-event; !ok {
				log.Info("watchComet exit")
				return
			}
			ins, ok := resolver.Fetch()
			if ok {
				if err := j.newAddress(ins.Instances); err != nil {
					log.Error("watchComet newAddress(%+v) error(%+v)", ins, err)
					continue
				}
				log.Info("watchComet change newAddress:%+v", ins)
			}
		}
	}()
}

func (j *Job) newAddress(insMap map[string][]*naming.Instance) error {
	ins := insMap[j.c.Env.Zone]
	if len(ins) == 0 {
		return fmt.Errorf("watchComet instance is empty")
	}
	comets := map[string]*Comet{}
	for _, in := range ins {
		if old, ok := j.cometServers[in.Hostname]; ok {
			comets[in.Hostname] = old
			continue
		}
		c, err := NewComet(in, j.c.Comet)
		if err != nil {
			log.Error("watchComet NewComet(%+v) error(%v)", in, err)
			return err
		}
		comets[in.Hostname] = c
		log.Info("watchComet AddComet grpc:%+v", in)
	}
	for key, old := range j.cometServers {
		if _, ok := comets[key]; !ok {
			old.cancel()
			log.Info("watchComet DelComet:%s", key)
		}
	}
	j.cometServers = comets
	return nil
}
