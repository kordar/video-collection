package command

// CommonCommand 公共策略
type CommonCommand struct {
	RawInputArgs  []string
	RawOutputArgs []string
	*BaseCommand
}

func NewCommonCommand(rawInputArgs []string, rawOutputArgs []string, strategy *BaseCommand) *CommonCommand {
	return &CommonCommand{
		RawInputArgs:  rawInputArgs,
		RawOutputArgs: rawOutputArgs,
		BaseCommand:   strategy,
	}
}

func (r *CommonCommand) SetMediaFile() {
	r.GetTrans().MediaFile().SetRawInputArgs(r.RawInputArgs)
	r.GetTrans().MediaFile().SetOutputPath(r.Output)
	r.GetTrans().MediaFile().SetRawOutputArgs(r.RawOutputArgs)
}

func (r *CommonCommand) Execute() error {
	r.SetMediaFile()
	return r.BaseCommand.Execute()
}
