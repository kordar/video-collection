package base

// ProgressState 执行进度状态
type ProgressState int

// iota 初始化后会自动递增
const (
	MgrReady ProgressState = iota
	MgrRunning
	RetryReady
	RetryExit
	RetryRestart
)

const (
	MgrReadyLabel     string = "MgrReady"
	MgrRunningLabel   string = "MgrRunning"
	RetryReadyLabel   string = "running"
	RetryRestartLabel string = "restart"
	RetryExitLabel    string = "exit"
)

func (s ProgressState) String() string {
	switch s {
	case MgrReady:
		return MgrReadyLabel
	case MgrRunning:
		return MgrRunningLabel
	case RetryReady:
		return RetryReadyLabel
	case RetryExit:
		return RetryExitLabel
	case RetryRestart:
		return RetryRestartLabel
	default:
		return "Unknown"
	}
}
