package packet

type JpegUtil struct {
	parentByte byte
	img        []byte
	run        bool
}

func NewJpegUtil(byteLen int) *JpegUtil {
	return &JpegUtil{0, make([]byte, 0, byteLen), false}
}

// ScanJpeg 从字节流中扫描出jpeg数据包，详情查看：
func (u *JpegUtil) ScanJpeg(buffs []byte, f func([]byte)) {
	for _, buf := range buffs {
		if buf == 0xD8 && u.parentByte == 0xFF {
			u.run = true
			u.img = append(u.img, 0xFF)
		} else if buf == 0xD9 && u.parentByte == 0xFF {
			u.run = false
			u.img = append(u.img, 0xD9)
			f(u.img)
			u.img = u.img[:0]
		}

		if u.run {
			u.img = append(u.img, buf)
		}

		u.parentByte = buf
	}
}
