package video_collection_test

import (
	"encoding/base64"
	logger "github.com/kordar/gologger"
	video_collection "github.com/kordar/video-collection"
	"github.com/kordar/video-collection/packet"
	"github.com/xfrr/goffmpeg/transcoder"
	"testing"
	"time"
)

var configuration = &video_collection.Configuration{
	Name:                "test",
	Input:               0,
	FFmpegInputPath:     "rtsp://admin:a1234567@192.168.10.67:554/h264/ch1/sub/av_stream",
	FFmpegRawInputArgs:  []string{"-re"},
	Output:              0,
	FFmpegOutputPath:    "-",
	FFmpegRawOutputArgs: []string{"-r", "1", "-q:v", "2"},
	OutputType:          video_collection.Image2Pipe,
	RetryTime:           time.Time{},
}

func TestFFmpegCollection_Run(t *testing.T) {
	var jpegutil = packet.NewJpegUtil(200 * 1024 * 8)
	collection := video_collection.FFmpegCollection{
		Running: func(value transcoder.Progress, cfg *video_collection.Configuration, collect *video_collection.FFmpegCollection) {
			logger.Info("---->>>>>>>>>>>>>%v----", cfg)
		},
		Before: nil,
		After: func(cfg *video_collection.Configuration, collect *video_collection.FFmpegCollection) {
			logger.Info("---->>>>>>>>>>>>>%v----", cfg)
		},
		ExecPipe: func(buf []byte, cfg *video_collection.Configuration) {
			jpegutil.ScanJpeg(buf, func(bytes []byte) {
				str := base64.StdEncoding.EncodeToString(bytes)
				logger.Infof("%v", str)
			})
		},
	}
	go func() {
		time.Sleep(20 * time.Second)
		collection.Exit(configuration)
	}()
	err := collection.Run(configuration, &video_collection.DefaultRetry{MaxTimes: 10, WaitSeconds: []int{}})
	logger.Errorf("=================%v", err)
}
