package strategy

// CommonStrategy 公共策略
type CommonStrategy struct {
	RawInputArgs  []string
	RawOutputArgs []string
	*BaseStrategy
}

func NewCommonStrategy(rawInputArgs []string, rawOutputArgs []string, strategy *BaseStrategy) *CommonStrategy {
	return &CommonStrategy{
		RawInputArgs:  rawInputArgs,
		RawOutputArgs: rawOutputArgs,
		BaseStrategy:  strategy,
	}
}

func (r *CommonStrategy) SetMediaFile() {
	r.GetTrans().MediaFile().SetRawInputArgs(r.RawInputArgs)
	r.GetTrans().MediaFile().SetOutputPath(r.Output)
	r.GetTrans().MediaFile().SetRawOutputArgs(r.RawOutputArgs)
}

func (r *CommonStrategy) Execute() error {
	r.SetMediaFile()
	return r.BaseStrategy.Execute()
}
