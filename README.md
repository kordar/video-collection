# 摄像头数据采集、切片处理

```go
import log "github.com/sirupsen/logrus"

log.SetLevel(log.DebugLevel)
log.SetFormatter(&nested.Formatter{
TimestampFormat: "2006-01-02 15:04:05",
})
goutil.ConfigInit("./conf.ini")
stream := video2.GetStream("demo")
basePath := goutil.GetSystemValue("output_base_dir")
manager := video2.NewStreamManager(basePath, 20)
manager.Add(stream)
manager.Run()
```

## 基于HLS点播功能的处理器


```ini
[demo-hls]
id=abc
input=rtsp://admin:a1234567@192.168.10.63:554/h264/ch1/main/av_stream
;input=/Users/mac/Movies/demo.flv
; output_prefix=/Users/mac/Movies
type=hls_datetime
output=ddd
retry_seconds=5
; 点播列表长度，配合hls_time 保存最近几日的点播数据
;hls_list_size=4320
hls_list_size=4320

; 分片时间，单位秒
hls_time=60
```

## 通用处理器


```ini
[common]
id=common
input=rtsp://admin:a1234567@192.168.10.56:554/h264/ch1/main/av_stream
output=rtmp://127.0.0.1:1985/myapp/56
input_args=-re -rtsp_transport tcp
output_args=-c:v copy -f flv
;output_args=-c:v libx264 -qp 51 -profile:v high -preset:v ultrafast -level 4.1 -x264opts crf=10 -an -f flv
retry_seconds=1
```