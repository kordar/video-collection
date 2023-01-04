# 摄像头数据采集、切片处理

## 基于HLS点播功能的数据采集


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