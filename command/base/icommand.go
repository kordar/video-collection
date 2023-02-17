package base

import (
	"time"
)

type ICommand interface {
	Execute() error
	GetId() string
	JustRestart()
	GetStatus() ProgressState // 获取策略状态 exit|restart|run
	Refresh()
	Stop()
	GetProgressRefreshTime() time.Time // Progress运行会刷该值
	GetProgressRestartSeconds() int64  // 如果程序产生假死状态，使用该值进行关闭或重启
}
