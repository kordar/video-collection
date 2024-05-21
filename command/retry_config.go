package command

import (
	logger "github.com/kordar/gologger"
	"time"
)

// RetryConfig 重试策略
type RetryConfig struct {
	RetryId          string
	RetryStatus      int       // 状态: 0 正常 1 完成
	RetryRefreshTime time.Time // 刷新时间
	RetrySeconds     int64     // 重试间隔秒数
	RetryMaxTimes    int       // 最大重试次数: 0 无限重试
	times            int       // 重试次数
	mseconds         int64
	RetrySs          map[int]int64 // 重启策略
}

func NewRetryConfig(retryId string, retrySeconds int64, retryMaxTimes int, retrySs map[int]int64) *RetryConfig {
	return &RetryConfig{RetryId: retryId, RetrySeconds: retrySeconds, RetryMaxTimes: retryMaxTimes, RetrySs: retrySs}
}

// Reset 重置重启对象状态ready，刷新时间为当前时间
func (r *RetryConfig) Reset() {
	r.RetryStatus = RetryStatusReady
	r.RetryRefreshTime = time.Now()
}

func (r *RetryConfig) ListenProgressRunning(c *BaseCommand) {
	// 状态为exit标识progress退出
	if r.RetryStatus == RetryStatusExit {
		c.Stop()
		return
	}

	if r.RetryStatus != RetryStatusReady {
		r.RetryStatus = RetryStatusReady
		r.RetryRefreshTime = time.Now()
		r.times = 0
	}
}

func (r *RetryConfig) ListenProgressFinish() {
	//
	if r.RetryStatus == RetryStatusExit {
		return
	}
	r.RetryStatus = RetryStatusFinish
	r.RetryRefreshTime = time.Now()
}

func (r *RetryConfig) SetExit() {
	r.RetryStatus = RetryStatusExit
	r.RetryRefreshTime = time.Now()
}

func (r *RetryConfig) GetStatus() string {
	s := r.GetStatus2()
	if "exit" == s {
		return "exit"
	}
	if "restart" == s {
		r.times += 1
		if r.RetryMaxTimes > 0 && r.times >= r.RetryMaxTimes {
			return "exit"
		}
		// 否则根据策略进行时间设置
		if r.times > 1000 {
			r.times = 1
		}
		mod := r.times % len(r.RetrySs)
		r.mseconds = r.RetrySs[mod]
	}
	return s
}

func (r *RetryConfig) GetStatus2() string {
	if r.RetryStatus == RetryStatusReady {
		return "run"
	}

	if r.RetryStatus == RetryStatusExit {
		return "exit"
	}

	s := time.Now().Unix() - r.RetryRefreshTime.Unix()
	if s > r.RetrySeconds {
		return "restart"
	}

	seconds := time.Duration(r.RetrySeconds-s) * time.Second
	logger.Infof("服务(%s)重启倒计时(%v)秒, 重试次数=%v次", r.RetryId, seconds+time.Duration(r.mseconds)*time.Second, r.times+1)
	<-time.After(seconds + time.Duration(r.mseconds)*time.Second)

	// TODO 等待结束期间如果被设置为结束状态，则等待完成后立即结束程序。
	if r.RetryStatus == RetryStatusExit {
		return "exit"
	}

	return "restart"
}
