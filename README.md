# 摄像头数据采集、切片处理

## FFmpeg处理器

```go
var configuration = &video_collection.Configuration{
    Name:                "test",
    Input:               0,
    FFmpegInputPath:     "rtsp://admin:a1234567@192.168.10.67:554/h264/ch1/sub/av_stream",
    FFmpegRawInputArgs:  []string{"-re"},
    Output:              0,
    FFmpegOutputPath:    "-",
    FFmpegRawOutputArgs: []string{"-r", "1", "-q:v", "2"},
    OutputType:          video_collection.Image2Pipe,
    RetryTime:           nil,
}

var jpegutil = packet.NewJpegUtil(200 * 1024 * 8)
collection := video_collection.FFmpegCollection{
        Running: func(value transcoder.Progress, cfg *video_collection.Configuration, collect *video_collection.FFmpegCollection) {
            logger.Info("---->>>>>>>>>>>>>----", value)
        },
        Before: nil,
        After:  nil,
        ExecPipe: func(buf []byte, cfg *video_collection.Configuration) {
            jpegutil.ScanJpeg(buf, func(bytes []byte) {
            str := base64.StdEncoding.EncodeToString(bytes)
            logger.Infof("%v", str)
        })
    },
}

err := collection.Run(configuration, &video_collection.DefaultRetry{MaxTimes: 10, WaitSeconds: []int{}})
logger.Errorf("=================%v", err)
```