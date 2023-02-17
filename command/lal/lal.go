package lal

import (
	base2 "github.com/kordar/video-collection/command/base"
	"github.com/kordar/video-collection/util"
	"github.com/spf13/cast"
)

type BaseLalCommand struct {
	OverTcp      int   // 0 非tcp 1 Tcp
	ProgressRate int64 // Progress打印频率，秒
	*base2.AbstractBaseCommand
}

func NewBaseLalCommand(commandID string, commandName string, input string,
	output string, params map[string]interface{}, retryConf *base2.RetryConfig) *BaseLalCommand {
	progressRate := cast.ToInt64(params["progressRate"])
	return &BaseLalCommand{
		OverTcp:             cast.ToInt(params["overTcp"]),
		ProgressRate:        util.GetValueInt64(progressRate, 1),
		AbstractBaseCommand: base2.NewAbstractBaseCommand(commandID, commandName, input, output, params, retryConf),
	}
}

func (b *BaseLalCommand) Stop() {
	b.RetryConfig.SetExit()
}

func (b *BaseLalCommand) JustRestart() {
	b.RetryConfig.ListenProgressFinish()
}
