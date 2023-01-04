package retry

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

var ss = map[int]int64{1: 5, 2: 5, 3: 10, 4: 10, 5: 30, 6: 30, 7: 60, 8: 60, 9: 300, 10: 300, 11: 1800, 12: 1800}

type Config struct {
	Id            string
	Status        int       // 状态: 0 正常 1 完成
	RefreshTime   time.Time // 刷新时间
	RetrySeconds  int64     // 重试间隔秒数
	MaxRetryTimes int       // 最大重试次数: 0 无线重试
	times         int       // 重试次数
	mseconds      int64
}

func (c *Config) Refresh() {
	c.Status = 0
	c.RefreshTime = time.Now()
}

func (c *Config) Clear() {
	if c.Status != 0 {
		c.Status = 0
		c.RefreshTime = time.Now()
		c.times = 0
	}
}

func (c *Config) End() {
	c.Status = 1
	c.RefreshTime = time.Now()
}

func (c *Config) GetStatus() string {
	s := c.GetStatus2()
	if "restart" == s {
		c.times += 1
		if c.MaxRetryTimes > 0 && c.times > c.MaxRetryTimes {
			return "exit"
		}
		// 否则根据策略进行时间设置
		if c.times > 1000 {
			c.times = 1
		}
		mod := c.times % 12
		c.mseconds = ss[mod]
	}
	return s
}

func (c *Config) GetStatus2() string {
	if c.Status == 0 {
		return "run"
	}

	s := time.Now().Unix() - c.RefreshTime.Unix()
	if s > c.RetrySeconds {
		return "restart"
	}
	seconds := time.Duration(c.RetrySeconds-s) * time.Second
	log.Infoln(fmt.Sprintf("服务(%s)重启倒计时(%v)秒, 重试次数=%v次", c.Id, seconds+time.Duration(c.mseconds)*time.Second, c.times+1))
	<-time.After(seconds + time.Duration(c.mseconds)*time.Second)
	return "restart"
}
