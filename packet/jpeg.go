package packet

import (
	"bytes"
	"encoding/binary"
)

type Jpeg struct {
	SOI  []byte // D8          文件头
	EOI  []byte // D9          文件尾
	SOF0 []byte // C0          帧开始（标准 JPEG）
	SOF1 []byte // C1          同上
	DHT  []byte // C4          定义 Huffman 表（霍夫曼表）
	SOS  []byte // DA          扫描行开始
	DQT  []byte // DB          定义量化表
	DRI  []byte // DD          定义重新开始间隔
	APP0 []byte // E0          定义交换格式和图像识别信息
	COM  []byte // FE          注释
	pk   int64
}

func (j *Jpeg) Scan(data []byte) {
	for _, b2 := range data {
		switch b2 {
		case 0xD8:
			j.writeSOI(data)
			break
		case 0xD9:
			j.writeEOI(data)
			return
		case 0xC0:
			j.writeSOF0(data)
			break
		default:
			j.pk++
		}
	}
}

func (j *Jpeg) writeSOI(data []byte) {
	j.SOI = []byte{data[j.pk-1], data[j.pk]}
	j.pk++
}

func (j *Jpeg) writeEOI(data []byte) {
	j.SOI = []byte{data[j.pk-1], data[j.pk]}
}

func (j *Jpeg) writeSOF0(data []byte) {
	j.SOI = []byte{data[j.pk-1], data[j.pk]}
	j.pk++
	d := 0
	bytebuf := bytes.NewBuffer([]byte{data[j.pk], data[j.pk+1]})
	binary.Write(bytebuf, binary.BigEndian, d)

}
