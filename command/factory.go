package command

import (
	"github.com/kordar/video-collection/command/base"
	"github.com/kordar/video-collection/command/ffmpeg"
	"github.com/kordar/video-collection/command/lal"
	"github.com/spf13/cast"
	"strings"
)

func GetStream(params map[string]interface{}, retrySs map[int]int64, callback *base.ProgressCallback) base.ICommand {
	t := cast.ToString(params["type"])
	switch t {
	case "hls_datetime":
		return hlsDatetimeStream(params, retrySs, callback)
	case "pullrtsp2pushrtmp":
		return pullRtsp2PushRtmp(params, retrySs, callback)
	default:
		return commonStream(params, retrySs, callback)
	}
}

func baseFfmpegConfig(params map[string]interface{}, retrySs map[int]int64, name string) *ffmpeg.BaseFfmpegCommand {
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
	//restartProcessSeconds := goutil.GetSectionValueInt(section, "restart_process_seconds")
	//if restartProcessSeconds == 0 {
	//	restartProcessSeconds = goutil.GetSectionValueInt("system", "restart_process_seconds")
	//}

	retryConfig := base.NewRetryConfig(id, int64(retrySeconds), retryMaxTimes, retrySs)
	return ffmpeg.NewBaseFfmpegCommand(id, name, input, output, params, retryConfig)
}

func hlsDatetimeStream(params map[string]interface{}, retrySs map[int]int64, callback *base.ProgressCallback) base.ICommand {
	hlsTime := cast.ToInt(params["hls_time"])
	hlsListSize := cast.ToInt(params["hls_list_size"])
	ffmpegConfig := baseFfmpegConfig(params, retrySs, "Hls采集")
	ffmpegConfig.Callback = callback
	outputDir := cast.ToString(params["output_dir"])
	return ffmpeg.NewFfmpegHlsDatetime(hlsTime, hlsListSize, outputDir, ffmpegConfig)
}

func commonStream(params map[string]interface{}, retrySs map[int]int64, callback *base.ProgressCallback) base.ICommand {
	ffmpegConfig := baseFfmpegConfig(params, retrySs, "推送流")
	ffmpegConfig.Callback = callback
	// ProfileV      string // 编码质量，https://gist.github.com/jedfoster/3c0b396097783f5884fb
	// -threads 5 -re -rtsp_transport tcp
	input := cast.ToString(params["input_args"])
	inputArgs := strings.Split(input, " ")

	// -c:v libx264 -qp 51  -profile:v high -threads 5 -preset:v ultrafast -level 4.1 -x264opts crf=10 -an -f flv
	output := cast.ToString(params["output_args"])
	outputArgs := strings.Split(output, " ")

	return ffmpeg.NewFfmpegCommonCommand(inputArgs, outputArgs, ffmpegConfig)
}

func baseLalConfig(params map[string]interface{}, retrySs map[int]int64, name string) *lal.BaseLalCommand {
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
	//restartProcessSeconds := goutil.GetSectionValueInt(section, "restart_process_seconds")
	//if restartProcessSeconds == 0 {
	//	restartProcessSeconds = goutil.GetSectionValueInt("system", "restart_process_seconds")
	//}

	retryConfig := base.NewRetryConfig(id, int64(retrySeconds), retryMaxTimes, retrySs)
	return lal.NewBaseLalCommand(id, name, input, output, params, retryConfig)
}

func pullRtsp2PushRtmp(params map[string]interface{}, retrySs map[int]int64, callback *base.ProgressCallback) base.ICommand {
	lalConfig := baseLalConfig(params, retrySs, "lal推送流客户端")
	lalConfig.Callback = callback
	return lal.NewPullRtsp2PushRtmpCommand(lalConfig)
}
