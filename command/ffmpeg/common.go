package ffmpeg

// CommonCommand 公共策略
type CommonCommand struct {
	RawInputArgs  []string
	RawOutputArgs []string
	*BaseFfmpegCommand
}

func NewFfmpegCommonCommand(rawInputArgs []string, rawOutputArgs []string, strategy *BaseFfmpegCommand) *CommonCommand {
	return &CommonCommand{
		RawInputArgs:      rawInputArgs,
		RawOutputArgs:     rawOutputArgs,
		BaseFfmpegCommand: strategy,
	}
}

func (r *CommonCommand) SetMediaFile() {
	r.GetTrans().MediaFile().SetRawInputArgs(r.RawInputArgs)
	// TODO pipe输出，自己进行维护
	if r.Output == "-" {
		r.GetTrans().MediaFile().SetOutputPath("")
	}
	r.GetTrans().MediaFile().SetRawOutputArgs(r.RawOutputArgs)
}

func (r *CommonCommand) Execute() error {
	r.SetMediaFile()
	return r.BaseFfmpegCommand.Execute()
}
