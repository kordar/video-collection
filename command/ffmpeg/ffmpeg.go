package ffmpeg

import (
	"github.com/kordar/video-collection/command/base"
	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/xfrr/goffmpeg/transcoder"
	"time"
)

type BaseFfmpegCommand struct {
	trans *transcoder.Transcoder // ffmpeg处理对象
	*base.AbstractBaseCommand
}

func NewBaseFfmpegCommand(commandID string, commandName string, input string,
	output string, params map[string]interface{}, retryConf *base.RetryConfig) *BaseFfmpegCommand {
	trans := new(transcoder.Transcoder)
	err := trans.Initialize(input, "")
	// Handle error...
	if err != nil {
		nazalog.Panicf("init trans err = %+v", err)
	}
	return &BaseFfmpegCommand{
		AbstractBaseCommand: base.NewAbstractBaseCommand(commandID, commandName, input, output, params, retryConf),
		trans:               trans,
	}
}

// GetTrans 获取ffmpeg处理对象
func (b *BaseFfmpegCommand) GetTrans() *transcoder.Transcoder {
	return b.trans
}

func (b *BaseFfmpegCommand) Execute() error {
	nazalog.Debug(b.GetTrans().GetCommand())
	nazalog.Infof("服务(%s:%s)启动: %s -> %s", b.CommandName, b.CommandID, b.Input, b.Output)

	/**
	 * progress 结束后，监听Progress结束尝试设置为重启状态
	 */
	defer func() {
		b.RetryConfig.ListenProgressFinish()
		nazalog.Warnf("服务(%s:%s)结束！！！", b.CommandName, b.CommandID)
		b.Callback.AfterFunc(b.AbstractBaseCommand)
	}()

	b.Callback.BeforeFunc(b.AbstractBaseCommand)

	// Start transcoder process without checking progress
	done := b.GetTrans().Run(true)
	b.ProgressRefreshTime = time.Now()

	output := b.GetTrans().Output()
	for progress := range output {
		// TODO 采集Progress最新刷新时间
		b.ProgressRefreshTime = time.Now()
		// 重试策略执行
		b.RetryConfig.ListenProgressRunning(b.AbstractBaseCommand)
		if progress.FramesProcessed != "" {
			nazalog.Infof("服务(%s:%s)成功, Process = %+v", b.CommandName, b.CommandID, progress)
			b.Callback.RunningFunc(progress, b.AbstractBaseCommand)
		} else {
			nazalog.Warnf("服务(%s:%s)失败, Process = %+v", b.CommandName, b.CommandID, progress)
		}
	}

	// This channel is used to wait for the transcoding process to end
	err := <-done

	return err
}

func (b *BaseFfmpegCommand) Stop() {
	err := b.GetTrans().Process().Process.Kill()
	b.RetryConfig.SetExit()
	if err != nil {
		nazalog.Errorf("[%s] ffmpeg Stop, err = %+v", b.CommandID, err)
		return
	}
}

func (b *BaseFfmpegCommand) JustRestart() {
	err := b.GetTrans().Stop()
	if err != nil {
		nazalog.Errorf("[%s] ffmpeg JustRestart, err = %+v", b.CommandID, err)
		return
	}
}
