package video_collection

import (
	logger "github.com/kordar/gologger"
	"github.com/xfrr/goffmpeg/transcoder"
)

type VideoFfmpeg struct {
	trans *transcoder.Transcoder // ffmpeg处理对象
}

func (v *VideoFfmpeg) Run(config Configuration) error {
	v.trans = new(transcoder.Transcoder)
	err := v.trans.Initialize(config.InputPath, "")
	if err != nil {
		logger.Warnf("Error initializing transcoder %s", err)
	}

	v.trans.MediaFile().SetRawInputArgs(config.RawInputArgs)
	v.trans.MediaFile().SetOutputPath(config.OutputPath)
	v.trans.MediaFile().SetRawOutputArgs(config.RawOutputArgs)

	logger.Debug(v.trans.GetCommand())
	logger.Infof("服务(%s:%s)启动: %s -> %s", b.CommandName, b.CommandID, b.Input, b.Output)

	//b.Callback.BeforeFunc(b)

	// Run transcoder process without checking progress
	done := v.trans.Run(true)
	//b.ProgressRefreshTime = time.Now()

	output := v.trans.Output()
	for progress := range output {
		// TODO 采集Progress最新刷新时间
		//b.ProgressRefreshTime = time.Now()
		// 重试策略执行
		//b.retry.ListenProgressRunning(b)
		if progress.FramesProcessed != "" {
			//logger.Infof("服务(%s:%s)成功, Process = %+v", b.CommandName, b.CommandID, progress)
			//b.Callback.RunningFunc(progress, b)
		} else {
			//logger.Warnf("服务(%s:%s)失败, Process = %+v", b.CommandName, b.CommandID, progress)
		}
	}

	// This channel is used to wait for the transcoding process to end
	err = <-done

	//logger.Warnf("服务(%s:%s)结束！！！", b.CommandName, b.CommandID)
	//b.Callback.AfterFunc(b)
	/**
	 * progress 结束后，监听Progress结束尝试设置为重启状态
	 */
	//b.retry.ListenProgressFinish()
	return err
}

func (v *VideoFfmpeg) Stop() {
	err := v.trans.Stop()
	//b.retry.SetExit()
	if err != nil {
		logger.Error(err)
		return
	}
}
