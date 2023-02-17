package base

import "github.com/xfrr/goffmpeg/models"

type ProgressCallback struct {
	runningFunc func(models.Progress, *AbstractBaseCommand) // 运行中回调
	beforeFunc  func(strategy *AbstractBaseCommand)         // 执行前回调
	afterFunc   func(strategy *AbstractBaseCommand)         // 执行后回调
}

// RunningFunc 执行中回调
func (p *ProgressCallback) RunningFunc(progress models.Progress, strategy *AbstractBaseCommand) {
	if p.runningFunc != nil {
		p.runningFunc(progress, strategy)
	}
}

func (p *ProgressCallback) SetRunningFunc(runningFunc func(models.Progress, *AbstractBaseCommand)) {
	p.runningFunc = runningFunc
}

// BeforeFunc 执行前回调
func (p *ProgressCallback) BeforeFunc(strategy *AbstractBaseCommand) {
	if p.beforeFunc != nil {
		p.beforeFunc(strategy)
	}
}

func (p *ProgressCallback) SetBeforeFunc(beforeFunc func(strategy *AbstractBaseCommand)) {
	p.beforeFunc = beforeFunc
}

// AfterFunc 执行后回调
func (p *ProgressCallback) AfterFunc(strategy *AbstractBaseCommand) {
	if p.afterFunc != nil {
		p.afterFunc(strategy)
	}
}

func (p *ProgressCallback) SetAfterFunc(afterFunc func(strategy *AbstractBaseCommand)) {
	p.afterFunc = afterFunc
}
