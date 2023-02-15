package command

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/xfrr/goffmpeg/models"
	"github.com/xfrr/goffmpeg/transcoder"
	"time"
)

type ProgressCallback struct {
	runningFunc func(models.Progress, *BaseCommand) // 运行中回调
	beforeFunc  func(strategy *BaseCommand)         // 执行前回调
	afterFunc   func(strategy *BaseCommand)         // 执行后回调
}

// RunningFunc 执行中回调
func (p *ProgressCallback) RunningFunc(progress models.Progress, strategy *BaseCommand) {
	if p.runningFunc != nil {
		p.runningFunc(progress, strategy)
	}
}

func (p *ProgressCallback) SetRunningFunc(runningFunc func(models.Progress, *BaseCommand)) {
	p.runningFunc = runningFunc
}

// BeforeFunc 执行前回调
func (p *ProgressCallback) BeforeFunc(strategy *BaseCommand) {
	if p.beforeFunc != nil {
		p.beforeFunc(strategy)
	}
}

func (p *ProgressCallback) SetBeforeFunc(beforeFunc func(strategy *BaseCommand)) {
	p.beforeFunc = beforeFunc
}

// AfterFunc 执行后回调
func (p *ProgressCallback) AfterFunc(strategy *BaseCommand) {
	if p.afterFunc != nil {
		p.afterFunc(strategy)
	}
}

func (p *ProgressCallback) SetAfterFunc(afterFunc func(strategy *BaseCommand)) {
	p.afterFunc = afterFunc
}

type ICommand interface {
	Execute() error
	GetId() string
	//GetBaseDir() string
	//SetBaseDir(dir string)
	GetStatus() string // 获取策略状态 exit|restart|run
	Refresh()
	Stop()
	GetProgressRefreshTime() time.Time
	GetProgressRestartSeconds() int64
}

// BaseCommand 公共命令处理器
type BaseCommand struct {
	CommandID   string                 // 命令Id
	CommandName string                 // 命令名称
	Input       string                 // 输入参数
	Output      string                 // 输出参数
	Params      map[string]interface{} // 扩展参数配置

	// Progress刷新时间
	ProgressRefreshTime    time.Time
	ProgressRestartSeconds int64
	//
	Callback *ProgressCallback

	trans *transcoder.Transcoder // ffmpeg处理对象
	retry *RetryConfig
}

func NewBaseCommand(commandID string, commandName string, input string, output string, params map[string]interface{}, retryConf *RetryConfig) *BaseCommand {
	trans := new(transcoder.Transcoder)
	err := trans.Initialize(input, "")
	// Handle error...
	if err != nil {
		log.Fatalln(err)
	}
	progressRestartSeconds := cast.ToInt64(params["progressRestartSeconds"])
	return &BaseCommand{CommandID: commandID,
		CommandName:            commandName,
		Input:                  input,
		Output:                 output,
		Params:                 params,
		ProgressRestartSeconds: progressRestartSeconds,
		retry:                  retryConf,
		trans:                  trans,
	}
}

// GetTrans 获取ffmpeg处理对象
func (b *BaseCommand) GetTrans() *transcoder.Transcoder {
	return b.trans
}

func (b *BaseCommand) Execute() error {
	log.Debugln(b.GetTrans().GetCommand())
	log.Infoln(fmt.Sprintf("服务(%s:%s)启动: %s -> %s", b.CommandName, b.CommandID, b.Input, b.Output))

	b.Callback.BeforeFunc(b)

	// Start transcoder process without checking progress
	done := b.GetTrans().Run(true)
	b.ProgressRefreshTime = time.Now()

	output := b.GetTrans().Output()
	for progress := range output {
		// TODO 采集Progress最新刷新时间
		b.ProgressRefreshTime = time.Now()
		// 重试策略执行
		b.retry.ListenProgressRunning(b)
		if progress.FramesProcessed != "" {
			log.Infof("服务(%s:%s)成功, Process = %+v", b.CommandName, b.CommandID, progress)
			b.Callback.RunningFunc(progress, b)
		} else {
			log.Warningf("服务(%s:%s)失败, Process = %+v", b.CommandName, b.CommandID, progress)
		}
	}

	// This channel is used to wait for the transcoding process to end
	err := <-done

	log.Warningln(fmt.Sprintf("服务(%s:%s)结束！！！", b.CommandName, b.CommandID))
	b.Callback.AfterFunc(b)
	/**
	 * progress 结束后，监听Progress结束尝试设置为重启状态
	 */
	b.retry.ListenProgressFinish()
	return err
}

func (b *BaseCommand) GetId() string {
	return b.CommandID
}

func (b *BaseCommand) Stop() {
	err := b.GetTrans().Stop()
	b.retry.SetExit()
	if err != nil {
		log.Errorln(err)
		return
	}
}

func (b *BaseCommand) GetStatus() string {
	return b.retry.GetStatus()
}

func (b *BaseCommand) Refresh() {
	b.retry.Reset()
}

func (b *BaseCommand) GetProgressRefreshTime() time.Time {
	return b.ProgressRefreshTime
}

func (b *BaseCommand) GetProgressRestartSeconds() int64 {
	return b.ProgressRestartSeconds
}
