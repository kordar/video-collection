package video_collection

import (
	"errors"
	"fmt"
	logger "github.com/kordar/gologger"
	"github.com/xfrr/goffmpeg/transcoder"
	"io"
	"time"
)

type FFmpegCollection struct {
	trans    *transcoder.Transcoder                                                         // ffmpeg处理对象
	Running  func(value transcoder.Progress, cfg *Configuration, collect *FFmpegCollection) // 运行中回调
	Before   func(cfg *Configuration, collect *FFmpegCollection)                            // 执行前回调
	After    func(cfg *Configuration, collect *FFmpegCollection)                            // 执行后回调
	ExecPipe func(buff []byte, cfg *Configuration)
}

func (v *FFmpegCollection) Run(config *Configuration, retry Retry) error {
	if config.ProgressStatus == StartStatusRunning {
		message := fmt.Sprintf("[%s] the progress is running, start failed.", config.Name)
		return errors.New(message)
	}

	// 启动时重置重试状态为ready
	config.RetryStatus = RetryStatusReady
	return v.exec(config, retry)
}

func (v *FFmpegCollection) exec(config *Configuration, retry Retry) error {

	config.ProgressStatus = StartStatusReady

	if v.Running == nil {
		logger.Errorf("[%s] please set the run callback function.", config.Name)
		return errors.New("please set the run callback function")
	}

	v.trans = new(transcoder.Transcoder)
	err := v.trans.Initialize(config.FFmpegInputPath, "")
	if err != nil {
		logger.Errorf("[%s] %+v", config.Name, err)
		return err
	}

	v.trans.MediaFile().SetRawInputArgs(config.FFmpegRawInputArgs)
	v.trans.MediaFile().SetOutputPath(config.FFmpegOutputPath)
	v.trans.MediaFile().SetRawOutputArgs(config.FFmpegRawOutputArgs)

	logger.Debug(v.trans.GetCommand())
	logger.Infof("[%s] service startup: %s -> %s", config.Name, config.FFmpegInputPath, config.FFmpegOutputPath)

	if v.Before != nil {
		v.Before(config, v)
	}

	if config.OutputType == Image2Pipe {
		v.trans.MediaFile().SetOutputPath("")
		pip, err2 := v.trans.CreateOutputPipe("image2pipe")
		if err2 != nil {
			return err2
		}
		go func() {
			v.PipeRead(pip, config.FFmpegPipeBuffSize, config)
		}()
		config.ProgressStatus = StartStatusRunning

		done := v.trans.Run(false)
		doneErr := <-done

		if config.RetryStatus == RetryStatusReady && retry != nil {
			config.RetryTime = time.Now()
			retry.Execute(config, v)
		}

		config.ProgressStatus = StartStatusFinish
		return doneErr
	}

	// Run transcoder process without checking progress
	done := v.trans.Run(true)
	output := v.trans.Output()
	config.ProgressStatus = StartStatusRunning
	for progress := range output {
		// TODO 采集Progress最新刷新时间
		if progress.FramesProcessed != "" {
			logger.Infof("[%s] service startup success, progress = %+v", config.Name, progress)
			v.Running(progress, config, v)
		} else {
			logger.Warnf("[%s] service startup failed, progress = %+v", config.Name, progress)
		}
	}

	// This channel is used to wait for the transcoding process to end
	doneErr := <-done

	logger.Warnf("[%s] service finished！！！", config.Name)
	if v.After != nil {
		v.After(config, v)
	}

	if config.RetryStatus == RetryStatusReady && retry != nil {
		config.RetryTime = time.Now()
		retry.Execute(config, v)
	}

	config.ProgressStatus = StartStatusFinish
	return doneErr
}

func (v *FFmpegCollection) PipeRead(reader *io.PipeReader, buffSize int, configuration *Configuration) {
	if buffSize == 0 {
		buffSize = 128
	}
	buff := make([]byte, buffSize)
	logger.Infof("[%s] **************** Start Receiving ******************", configuration.Name)
	for {
		_, err := reader.Read(buff)
		if err != nil {
			logger.Errorf("[%s] PipeRead err: %v", configuration.Name, err)
			return
		}
		v.ExecPipe(buff, configuration)
	}
}

func (v *FFmpegCollection) Exit(cfg *Configuration) {
	_ = v.trans.Stop()
	cfg.RetryStatus = RetryStatusExit
	logger.Infof("exit the current progress")
}

func (v *FFmpegCollection) Reload(cfg *Configuration, retry Retry) error {
	_ = v.trans.Stop()
	if cfg.RetryStatus == RetryStatusExit {
		return nil
	}
	return v.exec(cfg, retry)
}
