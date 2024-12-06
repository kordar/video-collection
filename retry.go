package video_collection

import (
	logger "github.com/kordar/gologger"
	"time"
)

type Retry interface {
	Execute(configuration *Configuration, collection Collection)
}

type DefaultRetry struct {
	current     int   // 当前重试次数
	MaxTimes    int   `json:"max_times"`    // 最大重试次数
	WaitSeconds []int `json:"wait_seconds"` // 重试等待时间
}

func (d *DefaultRetry) Execute(configuration *Configuration, collection Collection) {

	// 最大重试次数，0无限重试
	if d.MaxTimes != 0 && d.current >= d.MaxTimes {
		logger.Warnf("[%s] the maximum number of retries has been reached, and the process is exiting.", configuration.Name)
		d.current = 0
		return
	}

	// 等待重试睡眠时间
	if d.WaitSeconds != nil && len(d.WaitSeconds) > 0 {
		l := len(d.WaitSeconds)
		index := d.current % l
		second := d.WaitSeconds[index]
		time.Sleep(time.Duration(second) * time.Second)
	}

	d.current++
	configuration.RetryCount = d.current
	logger.Infof("[%s] current retry count: %d times", configuration.Name, d.current)
	_ = collection.Reload(configuration, d)
}
