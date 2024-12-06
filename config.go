package video_collection

import (
	"strings"
	"time"
)

type Configuration struct {
	Name                string
	Input               Input
	Output              Output
	OutputType          OutputFormat // 输出类型
	RetryTime           time.Time    // 重试时间
	RetryCount          int          // 重试次数
	RetryStatus         Status       // 重试状态
	ProgressStatus      Status       // 当前进度状态
	Err                 error
	FFmpegInputPath     string
	FFmpegOutputPath    string
	FFmpegRawInputArgs  []string
	FFmpegRawOutputArgs []string
	FFmpegPipeBuffSize  int // 内存通道buff大小
}

type ConfigurationVO struct {
	Name                string       `json:"name"`
	Input               Input        `json:"input"`
	InputLabel          string       `json:"input_label"`
	Output              Output       `json:"output"`
	OutputLabel         string       `json:"output_label"`
	OutputType          OutputFormat `json:"output_type"`
	OutputTypeLabel     string       `json:"output_type_label"`
	RetryTime           string       `json:"retry_time"`
	RetryCount          int          `json:"retry_count"`
	RetryStatus         Status       `json:"retry_status"`
	RetryStatusLabel    string       `json:"retry_status_label"`
	ProgressStatus      Status       `json:"progress_status"`
	ProgressStatusLabel string       `json:"progress_status_label"`
	Err                 string       `json:"err"`
	FFmpegInputPath     string       `json:"ffmpeg_input_path"`
	FFmpegOutputPath    string       `json:"ffmpeg_output_path"`
	FFmpegRawInputArgs  string       `json:"ffmpeg_raw_input_args"`
	FFmpegRawOutputArgs string       `json:"ffmpeg_raw_output_args"`
	FFmpegPipeBuffSize  int          `json:"ffmpeg_pipe_buff_size"`
}

func (vo *ConfigurationVO) Load(configuration Configuration) {
	vo.Name = configuration.Name
	vo.Input = configuration.Input
	vo.InputLabel = configuration.Input.String()
	vo.Output = configuration.Output
	vo.OutputLabel = configuration.Output.String()
	vo.OutputType = configuration.OutputType
	vo.OutputTypeLabel = configuration.OutputType.String()
	vo.RetryTime = configuration.RetryTime.Format("2006-01-02 15:04:05")
	vo.RetryCount = configuration.RetryCount
	vo.RetryStatus = configuration.RetryStatus
	vo.RetryStatusLabel = configuration.RetryStatus.String()
	vo.ProgressStatus = configuration.ProgressStatus
	vo.ProgressStatusLabel = configuration.ProgressStatus.String()
	vo.Err = configuration.Err.Error()
	vo.FFmpegPipeBuffSize = configuration.FFmpegPipeBuffSize
	vo.FFmpegInputPath = configuration.FFmpegInputPath
	vo.FFmpegOutputPath = configuration.FFmpegOutputPath
	vo.FFmpegRawInputArgs = strings.Join(configuration.FFmpegRawInputArgs, " ")
	vo.FFmpegRawOutputArgs = strings.Join(configuration.FFmpegRawOutputArgs, " ")
}
