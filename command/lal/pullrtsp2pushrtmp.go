package lal

import (
	base2 "github.com/kordar/video-collection/command/base"
	"github.com/kordar/video-collection/util"
	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/remux"
	"github.com/q191201771/lal/pkg/rtmp"
	"github.com/q191201771/lal/pkg/rtsp"
	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/spf13/cast"
	"time"
)

type PullRtsp2PushRtmpCommand struct {
	*BaseLalCommand
}

func NewPullRtsp2PushRtmpCommand(strategy *BaseLalCommand) *PullRtsp2PushRtmpCommand {
	return &PullRtsp2PushRtmpCommand{
		BaseLalCommand: strategy,
	}
}

func (p *PullRtsp2PushRtmpCommand) Execute() error {

	pushSession := rtmp.NewPushSession(func(option *rtmp.PushSessionOption) {
		pushTimeoutMs := cast.ToInt(p.Params["pushTimeoutMs"])
		writeAvTimeoutMs := cast.ToInt(p.Params["writeAvTimeoutMs"])
		option.PushTimeoutMs = util.GetValueInt(pushTimeoutMs, 5000)
		option.WriteAvTimeoutMs = util.GetValueInt(writeAvTimeoutMs, 5000)
	})

	err := pushSession.Push(p.Output)
	if err != nil {
		nazalog.Fatalf("[%s:%s] (PullRtsp2PushRtmpCommand) -> pushSession error = %+v", p.CommandID, p.CommandName, err)
	}
	defer pushSession.Dispose()

	remuxer := remux.NewAvPacket2RtmpRemuxer().WithOnRtmpMsg(func(msg base.RtmpMsg) {
		err = pushSession.Write(rtmp.Message2Chunks(msg.Payload, &msg.Header))
		if err != nil {
			nazalog.Fatalf("[%s:%s] (PullRtsp2PushRtmpCommand) -> remuxer error = %+v", p.CommandID, p.CommandName, err)
		}
	})
	pullSession := rtsp.NewPullSession(remuxer, func(option *rtsp.PullSessionOption) {
		pullTimeoutMs := cast.ToInt(p.Params["pullTimeoutMs"])
		option.PullTimeoutMs = util.GetValueInt(pullTimeoutMs, 5000)
		option.OverTcp = p.OverTcp != 0
	})

	err = pullSession.Pull(p.Input)
	if err != nil {
		nazalog.Fatalf("[%s:%s] PullRtsp2PushRtmpCommand -> pullSession error = %+v", p.CommandID, p.CommandName, err)
	}
	defer pullSession.Dispose()

	p.Callback.BeforeFunc(p.AbstractBaseCommand)
	p.ProgressRefreshTime = time.Now()

	go func() {
		wait := time.Duration(p.ProgressRate) * time.Second
		for {

			// 重试状态不等于ready状态，则关闭session
			if p.RetryConfig.RetryStatus != base2.RetryReady {
				_ = pushSession.Dispose()
				_ = pullSession.Dispose()
				return
			}

			// TODO 采集Progress最新刷新时间
			p.ProgressRefreshTime = time.Now()
			// 重试策略执行
			p.RetryConfig.ListenProgressRunning(p.AbstractBaseCommand)

			pullSession.UpdateStat(1)
			pullStat := pullSession.GetStat()
			pushSession.UpdateStat(1)
			pushStat := pushSession.GetStat()
			nazalog.Infof("stat. pull=%+v, push=%+v", pullStat, pushStat)
			time.Sleep(wait)
		}
	}()

	select {
	case err = <-pullSession.WaitChan():
		nazalog.Infof("< pullSession.Wait(). err=%+v", err)
	case err = <-pushSession.WaitChan():
		nazalog.Infof("< pushSession.Wait(). err=%+v", err)
	}

	p.Callback.AfterFunc(p.AbstractBaseCommand)
	/**
	 * progress 结束后，监听Progress结束尝试设置为重启状态
	 */
	p.RetryConfig.ListenProgressFinish()

	return err
}
