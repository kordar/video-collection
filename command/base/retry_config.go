package base

import (
	"github.com/q191201771/naza/pkg/nazalog"
	"time"
)

// RetryConfig 重试策略
type RetryConfig struct {
	RetryId          string
	RetryStatus      ProgressState // 状态: 0 正常 1 完成
	RetryRefreshTime time.Time     // 刷新时间
	RetrySeconds     int64         // 重试间隔秒数
	RetryMaxTimes    int           // 最大重试次数: 0 无限重试
	times            int           // 重试次数
	mseconds         int64
	RetrySs          map[int]int64 // 重启策略
}

func NewRetryConfig(retryId string, retrySeconds int64, retryMaxTimes int, retrySs map[int]int64) *RetryConfig {
	return &RetryConfig{RetryId: retryId, RetrySeconds: retrySeconds, RetryMaxTimes: retryMaxTimes, RetrySs: retrySs}
}

// Reset 重置重启对象状态ready，刷新时间为当前时间
func (r *RetryConfig) Reset() {
	r.RetryStatus = RetryReady
	r.RetryRefreshTime = time.Now()
}

func (r *RetryConfig) ListenProgressRunning(c *AbstractBaseCommand) {
	// 状态为exit标识progress退出
	if r.RetryStatus == RetryExit {
		c.Stop()
		return
	}

	if r.RetryStatus == RetryReady {
		return
	}

	r.RetryStatus = RetryReady
	r.RetryRefreshTime = time.Now()
	r.times = 0
}

func (r *RetryConfig) ListenProgressFinish() {
	//
	if r.RetryStatus == RetryExit {
		return
	}
	r.RetryStatus = RetryRestart
	r.RetryRefreshTime = time.Now()
}

func (r *RetryConfig) SetExit() {
	r.RetryStatus = RetryExit
	r.RetryRefreshTime = time.Now()
}

func (r *RetryConfig) GetStatus() ProgressState {
	s := r.GetStatus2()
	if RetryExit == s {
		return RetryExit
	}
	if RetryRestart == s {
		r.times += 1
		if r.RetryMaxTimes > 0 && r.times >= r.RetryMaxTimes {
			return RetryExit
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

func (r *RetryConfig) GetStatus2() ProgressState {
	if r.RetryStatus == RetryReady {
		return RetryReady
	}

	if r.RetryStatus == RetryExit {
		return RetryExit
	}

	s := time.Now().Unix() - r.RetryRefreshTime.Unix()
	if s > r.RetrySeconds {
		return RetryRestart
	}

	seconds := time.Duration(r.RetrySeconds-s) * time.Second
	nazalog.Infof("服务(%s)重启倒计时(%v)秒, 重试次数=%v次", r.RetryId, seconds+time.Duration(r.mseconds)*time.Second, r.times+1)
	<-time.After(seconds + time.Duration(r.mseconds)*time.Second)

	// TODO 等待结束期间如果被设置为结束状态，则等待完成后立即结束程序。
	if r.RetryStatus == RetryExit {
		return RetryExit
	}

	return RetryRestart
}
