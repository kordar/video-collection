# 摄像头数据采集、切片处理

```go
package main

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/kordar/goutil"
	video2 "github.com/kordar/video-collection/command"
	log "github.com/sirupsen/logrus"
	"github.com/xfrr/goffmpeg/models"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	goutil.ConfigInit("./conf.ini")
	manager := video2.NewStreamManager(20)
	params := map[string]interface{}{
		"progressRestartSeconds": 600,
	}
	ss := map[int]int64{1: 5, 2: 5}
	callback := &video2.ProgressCallback{}
	callback.SetBeforeFunc(func(strategy *video2.BaseCommand) {
		log.Println("this is before func,", strategy.GetId())
	})
	callback.SetAfterFunc(func(strategy *video2.BaseCommand) {
		log.Println("after func,", strategy.GetId())
	})
	callback.SetRunningFunc(func(progress models.Progress, command *video2.BaseCommand) {
		log.Println(fmt.Sprintf("%+v", progress))
	})
	manager.Run()
	manager.Add(video2.GetStream("video56", params, ss, callback))
	manager.Add(video2.GetStream("video64", params, ss, callback))
	manager.Add(video2.GetStream("video63", params, ss, callback))
	manager.StartCheckDeath("@every 20s")
	done := make(chan int)
	<-done
	//time.Sleep(10000 * time.Second)
}
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