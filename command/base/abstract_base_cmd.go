package base

import (
	"github.com/spf13/cast"
	"time"
)

// AbstractBaseCommand 抽象处理器
type AbstractBaseCommand struct {
	CommandID   string                 // 命令Id
	CommandName string                 // 命令名称
	Input       string                 // 输入参数
	Output      string                 // 输出参数
	Params      map[string]interface{} // 扩展参数配置
	// Progress刷新时间
	ProgressRefreshTime    time.Time
	ProgressRestartSeconds int64
	//
	Callback    *ProgressCallback
	RetryConfig *RetryConfig
}

func NewAbstractBaseCommand(commandID string, commandName string, input string, output string, params map[string]interface{}, retryConf *RetryConfig) *AbstractBaseCommand {
	progressRestartSeconds := cast.ToInt64(params["progressRestartSeconds"])
	retryConf.RetryStatus = RetryReady
	return &AbstractBaseCommand{CommandID: commandID,
		CommandName:            commandName,
		Input:                  input,
		Output:                 output,
		Params:                 params,
		ProgressRestartSeconds: progressRestartSeconds,
		RetryConfig:            retryConf,
	}
}

func (b *AbstractBaseCommand) Execute() error {
	return nil
}

func (b *AbstractBaseCommand) GetId() string {
	return b.CommandID
}

func (b *AbstractBaseCommand) Stop() {

}

func (b *AbstractBaseCommand) JustRestart() {

}

func (b *AbstractBaseCommand) GetStatus() ProgressState {
	return b.RetryConfig.GetStatus()
}

func (b *AbstractBaseCommand) Refresh() {
	b.RetryConfig.Reset()
}

// GetProgressRefreshTime 最新刷新时间，配合进行超时检测
func (b *AbstractBaseCommand) GetProgressRefreshTime() time.Time {
	return b.ProgressRefreshTime
}

// GetProgressRestartSeconds 获取执行脚本timeout时间
func (b *AbstractBaseCommand) GetProgressRestartSeconds() int64 {
	return b.ProgressRestartSeconds
}
