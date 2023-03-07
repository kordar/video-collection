package packet

type JpegUtil struct {
	pbuf byte
	img  []byte
	run  bool
}

func NewJpegUtil(byteLen int) *JpegUtil {
	return &JpegUtil{0, make([]byte, 0, byteLen), false}
}

// ScanJpeg 从字节流中扫描出jpeg数据包，详情查看：
func (u *JpegUtil) ScanJpeg(buffs []byte, f func([]byte)) {
	for _, buff := range buffs {

		// 扫描到 0xFF 判断标识字段
		if u.pbuf == 0xFF {

			// 当前buf=0xD8,设置扫描开始
			if buff == 0xD8 {
				u.run = true
				u.img = append(u.img, 0xFF) // 修复标识位
			} else if buff == 0xD9 { // 查找结束位置
				u.pbuf = 0
				u.img = append(u.img, 0xD9)
				f(u.img)          // 图片回调
				u.img = u.img[:0] // 清空buff
				u.run = false
				continue
			}

		}

		if u.run {
			u.img = append(u.img, buff)
		}

		u.pbuf = buff

	}
}
