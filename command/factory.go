package command

import (
	"github.com/spf13/cast"
	"strings"
)

/**
 * type: 处理器类型
 * id: 处理器唯一id
 * input: 输入流
 * output: 输出流
 * retry_seconds: 重试秒数
 * ffmpeg_retry_seconds: ffmpeg重试秒数
 * retry_max_times: 重试最大次数 0 无限重试
 *
 */

func GetStream(params map[string]interface{},
	retrySs map[int]int64,
	callback *ProgressCallback,
) ICommand {
	t := cast.ToString(params["type"])
	switch t {
	case "hls_datetime":
		return hlsDatetimeStream(params, retrySs, callback)
	default:
		return commonStream(params, retrySs, callback)
	}
}

func baseConfig(params map[string]interface{}, retrySs map[int]int64, name string) *BaseCommand {
	id := cast.ToString(params["id"])
	input := cast.ToString(params["input"])
	output := cast.ToString(params["output"])
	retrySeconds := cast.ToInt(params["retry_seconds"])
	if retrySeconds == 0 {
		retrySeconds = cast.ToInt(params["ffmpeg_retry_seconds"])
	}
	if retrySeconds < 3 {
		retrySeconds = 30
	}
	retryMaxTimes := cast.ToInt(params["retry_max_times"])

	// 重启process，pipe网络等原因假死处理
	//restartProcessSeconds := cast.ToInt(params["restart_process_seconds"])
	//if restartProcessSeconds == 0 {
	//	restartProcessSeconds = cast.ToInt(params["restart_process_seconds"])
	//}

	retryConfig := NewRetryConfig(id, int64(retrySeconds), retryMaxTimes, retrySs)
	return NewBaseCommand(id, name, input, output, params, retryConfig)
}

func hlsDatetimeStream(params map[string]interface{}, retrySs map[int]int64, callback *ProgressCallback) ICommand {
	hlsTime := cast.ToInt(params["hls_time"])
	hlsListSize := cast.ToInt(params["hls_list_size"])
	base := baseConfig(params, retrySs, "Hls采集")
	base.Callback = callback
	outputDir := cast.ToString(params["output_dir"])
	return NewHlsDatetime(hlsTime, hlsListSize, outputDir, base)
}

func commonStream(params map[string]interface{}, retrySs map[int]int64, callback *ProgressCallback) ICommand {
	base := baseConfig(params, retrySs, "推送流")
	base.Callback = callback
	// ProfileV      string // 编码质量，https://gist.github.com/jedfoster/3c0b396097783f5884fb
	// -threads 5 -re -rtsp_transport tcp
	input := cast.ToString(params["input_args"])
	inputArgs := strings.Split(input, " ")

	// -c:v libx264 -qp 51  -profile:v high -threads 5 -preset:v ultrafast -level 4.1 -x264opts crf=10 -an -f flv
	output := cast.ToString(params["output_args"])
	outputArgs := strings.Split(output, " ")

	return NewCommonCommand(inputArgs, outputArgs, base)
}
