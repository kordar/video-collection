package strategy

import (
	"fmt"
	"github.com/kordar/video-collection/retry"
	log "github.com/sirupsen/logrus"
	"github.com/xfrr/goffmpeg/models"
	"github.com/xfrr/goffmpeg/transcoder"
	"path"
)

type Strategy interface {
	Execute() error
	GetId() string
	GetBaseDir() string
	SetBaseDir(dir string)
	GetStatus() string // 获取策略状态 exit|restart|run
	Refresh()
	Stop()
}

type BaseStrategy struct {
	ID                    string // 配置标识符
	Input                 string
	OutputPrefix          string // 路径前缀
	Output                string
	trans                 *transcoder.Transcoder
	retryConfig           *retry.Config
	ProcessName           string
	RestartProcessSeconds int64
	ProgressCallback      func(models.Progress, *BaseStrategy)
	BeforeCallback        func(strategy *BaseStrategy)
	AfterCallback         func(strategy *BaseStrategy)
}

func NewBaseStrategy(ID string, input string, output string, prefix string, retrySeconds int64, maxRetryTimes int, restartProcessSeconds int64) *BaseStrategy {
	trans := new(transcoder.Transcoder)
	err := trans.Initialize(input, "")
	// Handle error...
	if err != nil {
		log.Fatalln(err)
	}
	return &BaseStrategy{
		ID:                    ID,
		trans:                 trans,
		Input:                 input,
		Output:                output,
		OutputPrefix:          prefix,
		ProcessName:           "保存流",
		RestartProcessSeconds: restartProcessSeconds,
		retryConfig: &retry.Config{
			Id:            ID,
			RetrySeconds:  retrySeconds,
			MaxRetryTimes: maxRetryTimes,
		},
	}
}

func (b *BaseStrategy) GetTrans() *transcoder.Transcoder {
	return b.trans
}

func (b *BaseStrategy) OutputDir() string {
	return path.Join(b.OutputPrefix, b.Output)
}

func (b *BaseStrategy) GetRetryConfig() *retry.Config {
	return b.retryConfig
}

func (b *BaseStrategy) GetId() string {
	return b.ID
}

func (b *BaseStrategy) GetBaseDir() string {
	return b.OutputPrefix
}

func (b *BaseStrategy) SetBaseDir(dir string) {
	b.OutputPrefix = dir
}

func (b *BaseStrategy) Execute() error {

	log.Debugln(b.GetTrans().GetCommand())
	log.Infoln(fmt.Sprintf("服务(%v)启动: %s -> %s", b.ID, b.Input, b.Output))

	if b.BeforeCallback != nil {
		b.BeforeCallback(b)
	}

	// Start transcoder process without checking progress
	done := b.GetTrans().Run(true)

	output := b.GetTrans().Output()
	for progress := range output {
		WaitRestartHandler.Set(b)
		if progress.FramesProcessed != "" {
			b.GetRetryConfig().Clear()
			log.Infof("服务(%s)-%s成功, Process = %+v", b.ID, b.ProcessName, progress)
			if b.ProgressCallback != nil {
				b.ProgressCallback(progress, b)
			}
		} else {
			log.Warningf("服务(%s)-%s失败, Process = %+v", b.ID, b.ProcessName, progress)
		}
	}

	// This channel is used to wait for the transcoding process to end
	err := <-done
	WaitRestartHandler.Clear(b)
	b.GetRetryConfig().End()

	if b.AfterCallback != nil {
		b.AfterCallback(b)
	}

	log.Warningln(fmt.Sprintf("服务(%v)结束！！！", b.ID))
	return err
}

func (b *BaseStrategy) GetStatus() string {
	return b.GetRetryConfig().GetStatus()
}

func (b *BaseStrategy) Refresh() {
	b.GetRetryConfig().Refresh()
}

func (b *BaseStrategy) Stop() {
	err := b.GetTrans().Stop()
	b.GetRetryConfig().Exit()
	if err != nil {
		log.Errorln(err)
		return
	}
}
