package video

import (
	"github.com/kordar/goutil"
	strategy2 "github.com/kordar/video-collection/strategy"
	"strings"
)

func GetStream(section string) strategy2.Strategy {
	t := goutil.GetSectionValue(section, "type")
	switch t {
	case "hls_datetime":
		return hlsDatetimeStream(section)
	default:
		return commonStream(section)
	}
}

func baseConfig(section string) *strategy2.BaseStrategy {
	id := goutil.GetSectionValue(section, "id")
	if id == "" {
		id = section
	}
	input := goutil.GetSectionValue(section, "input")
	outputPrefix := goutil.GetSectionValue(section, "output_prefix")
	output := goutil.GetSectionValue(section, "output")
	retrySeconds := goutil.GetSectionValueInt(section, "retry_seconds")
	if retrySeconds == 0 {
		retrySeconds = goutil.GetSectionValueInt("system", "ffmpeg_retry_seconds")
	}
	if retrySeconds < 3 {
		retrySeconds = 30
	}
	maxRetryTimes := goutil.GetSectionValueInt(section, "max_retry_times")

	// 重启process，pipe网络等原因假死处理
	restartProcessSeconds := goutil.GetSectionValueInt(section, "restart_process_seconds")
	if restartProcessSeconds == 0 {
		restartProcessSeconds = goutil.GetSectionValueInt("system", "restart_process_seconds")
	}
	return strategy2.NewBaseStrategy(id, input, output, outputPrefix, int64(retrySeconds), maxRetryTimes, int64(restartProcessSeconds))
}

func hlsDatetimeStream(section string) strategy2.Strategy {
	hlsTime := goutil.GetSectionValueInt(section, "hls_time")
	hlsListSize := goutil.GetSectionValueInt(section, "hls_list_size")
	base := baseConfig(section)
	return strategy2.NewHlsDatetime(hlsTime, hlsListSize, base)
}

func commonStream(section string) strategy2.Strategy {
	base := baseConfig(section)
	base.ProcessName = "推送流"
	// ProfileV      string // 编码质量，https://gist.github.com/jedfoster/3c0b396097783f5884fb
	// -threads 5 -re -rtsp_transport tcp
	input := goutil.GetSectionValue(section, "input_args")
	inputArgs := strings.Split(input, " ")

	// -c:v libx264 -qp 51  -profile:v high -threads 5 -preset:v ultrafast -level 4.1 -x264opts crf=10 -an -f flv
	output := goutil.GetSectionValue(section, "output_args")
	outputArgs := strings.Split(output, " ")

	return strategy2.NewCommonStrategy(inputArgs, outputArgs, base)
}
