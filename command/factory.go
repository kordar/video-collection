package command

import (
	"github.com/kordar/goutil"
	"strings"
)

func GetStream(section string, params map[string]interface{}, retrySs map[int]int64, callback *ProgressCallback) ICommand {
	t := goutil.GetSectionValue(section, "type")
	switch t {
	case "hls_datetime":
		return hlsDatetimeStream(section, params, retrySs, callback)
	default:
		return commonStream(section, params, retrySs, callback)
	}
}

func baseConfig(section string, params map[string]interface{}, retrySs map[int]int64, name string) *BaseCommand {
	id := goutil.GetSectionValue(section, "id")
	if id == "" {
		id = section
	}
	input := goutil.GetSectionValue(section, "input")
	output := goutil.GetSectionValue(section, "output")
	retrySeconds := goutil.GetSectionValueInt(section, "retry_seconds")
	if retrySeconds == 0 {
		retrySeconds = goutil.GetSectionValueInt("system", "ffmpeg_retry_seconds")
	}
	if retrySeconds < 3 {
		retrySeconds = 30
	}
	retryMaxTimes := goutil.GetSectionValueInt(section, "retry_max_times")

	// 重启process，pipe网络等原因假死处理
	restartProcessSeconds := goutil.GetSectionValueInt(section, "restart_process_seconds")
	if restartProcessSeconds == 0 {
		restartProcessSeconds = goutil.GetSectionValueInt("system", "restart_process_seconds")
	}

	retryConfig := NewRetryConfig(id, int64(retrySeconds), retryMaxTimes, retrySs)
	return NewBaseCommand(id, name, input, output, params, retryConfig)
}

func hlsDatetimeStream(section string, params map[string]interface{}, retrySs map[int]int64, callback *ProgressCallback) ICommand {
	hlsTime := goutil.GetSectionValueInt(section, "hls_time")
	hlsListSize := goutil.GetSectionValueInt(section, "hls_list_size")
	base := baseConfig(section, params, retrySs, "Hls采集")
	base.Callback = callback
	outputDir := goutil.GetSectionValue(section, "output_dir")
	return NewHlsDatetime(hlsTime, hlsListSize, outputDir, base)
}

func commonStream(section string, params map[string]interface{}, retrySs map[int]int64, callback *ProgressCallback) ICommand {
	base := baseConfig(section, params, retrySs, "推送流")
	base.Callback = callback
	// ProfileV      string // 编码质量，https://gist.github.com/jedfoster/3c0b396097783f5884fb
	// -threads 5 -re -rtsp_transport tcp
	input := goutil.GetSectionValue(section, "input_args")
	inputArgs := strings.Split(input, " ")

	// -c:v libx264 -qp 51  -profile:v high -threads 5 -preset:v ultrafast -level 4.1 -x264opts crf=10 -an -f flv
	output := goutil.GetSectionValue(section, "output_args")
	outputArgs := strings.Split(output, " ")

	return NewCommonCommand(inputArgs, outputArgs, base)
}
