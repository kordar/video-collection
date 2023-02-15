package command

import (
	"github.com/kordar/video-collection/util"
	"path"
	"strconv"
	"time"
)

// HlsDatetimeCommand 通过时间配置切割
type HlsDatetimeCommand struct {
	OutPutDir   string
	HlsTime     int // 视频采集的间隔，单位秒
	HlsListSize int // 点播清单的最大分片数量
	*BaseCommand
}

func NewHlsDatetime(hlsTime int, hlsListSize int, outputDir string, strategy *BaseCommand) *HlsDatetimeCommand {
	return &HlsDatetimeCommand{
		HlsTime:     hlsTime,
		HlsListSize: hlsListSize,
		OutPutDir:   outputDir,
		BaseCommand: strategy,
	}
}

func (h *HlsDatetimeCommand) SetMediaFile() {
	output := h.OutPutDir
	util.CheckAndMkdir(output)
	timeout := strconv.FormatInt((5 * time.Second).Microseconds(), 10)
	h.GetTrans().MediaFile().SetRawInputArgs([]string{
		"-re",
		//"-buffer_size", "4086000",
		//"-listen_timeout", "5",
		"-rtsp_transport", "tcp",
		"-rtsp_flags", "prefer_tcp",
		//"-allowed_media_types", "video",
		//"-listen_timeout", timeout,
		//"-reorder_queue_size", "65535",
		//"-buffer_size", "65535",
		"-timeout", timeout,
	})
	//h.GetTrans().MediaFile().SetFrameRate(25)
	h.GetTrans().MediaFile().SetHlsListSize(h.HlsListSize)
	h.GetTrans().MediaFile().SetHlsSegmentDuration(h.HlsTime)
	//h.GetTrans().MediaFile().SetHlsPlaylistType("m3u8")
	h.GetTrans().MediaFile().SetHlsSegmentFilename(path.Join(output, "%Y", "%m", "%d", "segment_%Y%m%d%H%M%S.ts"))
	//h.GetTrans().MediaFile().SetHlsMasterPlaylistName(path.Join(output, "stream.m3u8"))
	h.GetTrans().MediaFile().SetVideoCodec("copy")
	h.GetTrans().MediaFile().SetAudioCodec("copy")
	h.GetTrans().MediaFile().SetOutputFormat("hls")
	h.GetTrans().MediaFile().SetOutputPath(path.Join(output, "stream.m3u8"))

	h.GetTrans().MediaFile().SetRawOutputArgs([]string{
		"-fflags", "flush_packets",
		"-max_delay", "1",
		"-an",
		//"-flags", "second_level_segment_index+second_level_segment_size+second_level_segment_duration",
		//"-fflags", "nobuffer",
		//"-pkt_size", "1300",
		//"-qscale:a", "4",
		"-strftime", "1",
		"-strftime_mkdir", "1",
		"-hls_flags", "omit_endlist+append_list",
	})
}

func (h *HlsDatetimeCommand) Execute() error {
	h.SetMediaFile()
	return h.BaseCommand.Execute()
}
