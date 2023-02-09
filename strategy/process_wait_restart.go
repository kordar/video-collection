package strategy

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var WaitRestartHandler = ProcessWaitRestart{
	strategy: make(map[*BaseStrategy]time.Time),
	locker:   sync.Mutex{},
}

// ProcessWaitRestart Process等待固定时间进行重启
type ProcessWaitRestart struct {
	strategy map[*BaseStrategy]time.Time
	locker   sync.Mutex
}

func (r *ProcessWaitRestart) Set(key *BaseStrategy) {
	// TODO 配置该值需大于5分钟
	if key.RestartProcessSeconds > 300 {
		r.strategy[key] = time.Now()
	}
}

func (r *ProcessWaitRestart) Clear(key *BaseStrategy) {
	r.locker.Lock()
	defer r.locker.Unlock()
	delete(r.strategy, key)
}

func (r *ProcessWaitRestart) Run() {
	c := cron.New()
	_, _ = c.AddFunc("@every 5m", func() {
		r.locker.Lock()
		defer r.locker.Unlock()
		for b, t := range r.strategy {
			now := time.Now().Unix()
			if b.RestartProcessSeconds > 0 && now-t.Unix() > b.RestartProcessSeconds {
				log.Info("强制停止编码器，id = ", b.GetId())
				b.GetTrans().Stop()
			}
		}
	})
	c.Start()
}
