package video

import (
	"github.com/kordar/goutil"
	strategy2 "videosys/strategy"
)

func GetStream(section string) strategy2.Strategy {
	t := goutil.GetSectionValue(section, "type")
	switch t {
	case "hls_datetime":
		return hlsDatetimeStream(section)
	default:
		return hlsDatetimeStream(section)
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
	return strategy2.NewBaseStrategy(id, input, output, outputPrefix, int64(retrySeconds), maxRetryTimes)
}

func hlsDatetimeStream(section string) strategy2.Strategy {
	hlsTime := goutil.GetSectionValueInt(section, "hls_time")
	hlsListSize := goutil.GetSectionValueInt(section, "hls_list_size")
	base := baseConfig(section)
	return strategy2.NewHlsDatetime(hlsTime, hlsListSize, base)
}
