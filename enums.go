package video_collection

type Input int

const (
	InputFilePath Input = iota
	InputRtsp
)

func (input Input) String() string {
	return [...]string{"FilePath", "Rtsp"}[input]
}

type Output int

const (
	OutputFilePath Output = iota
	OutputNone
)

func (output Output) String() string {
	return [...]string{"FilePath", "None"}[output]
}

type OutputFormat int

const (
	Image2Pipe OutputFormat = iota
)

func (output OutputFormat) String() string {
	return [...]string{"image2pipe"}[output]
}
