package video_collection

type Configuration struct {
	Input         Input
	InputPath     string
	RawInputArgs  []string
	Output        Output
	OutputPath    string
	RawOutputArgs []string
	OutputType    OutputFormat
}
