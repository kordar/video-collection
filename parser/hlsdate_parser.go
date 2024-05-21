package parser

import (
	"errors"
	"github.com/etherlabsio/go-m3u8/m3u8"
	logger "github.com/kordar/gologger"
	"github.com/kordar/video-collection/util"
	"path"
	"time"
)

var (
	layout = "20060102150405"
)

type HlsdateParser struct {
	playlist *m3u8.Playlist
}

func NewHlsdateParser(m3u8path string) *HlsdateParser {
	playlist, err := m3u8.ReadFile(m3u8path)
	if err != nil {
		logger.Fatal(err)
	}
	return &HlsdateParser{
		playlist: playlist,
	}
}

func (p HlsdateParser) findSeg(target string) *m3u8.SegmentItem {
	var current *m3u8.SegmentItem
	for _, item := range p.playlist.Segments() {
		base := path.Base(item.Segment)
		if target < base[8:22] {
			if current == nil {
				return item
			}
			break
		}
		current = item
	}
	return current
}

func (p HlsdateParser) parserofend(item *m3u8.SegmentItem, s string, e string, n string, start string, end string, sts *[]util.St) bool {
	if n != e {
		return false
	}

	// 开始结束在同一区间
	ntime, _ := time.Parse(layout, n)
	var sss int64 = 0
	if e == s && start >= s {
		stime, _ := time.Parse(layout, start)
		sss = stime.Unix() - ntime.Unix()
	}

	if sss < 0 || sss > int64(item.Duration) {
		return true
	}

	etime, _ := time.Parse(layout, end)
	ss := etime.Unix() - ntime.Unix()
	if ss > int64(item.Duration) {
		ss = int64(item.Duration)
	}

	if ss > 0 {
		st := util.St{
			Filename: item.Segment,
			Ss:       util.ToTimeStr(sss),
			To:       util.ToTimeStr(ss - sss),
			Duration: int64(item.Duration),
		}
		*sts = append(*sts, st)
	}

	return true
}

func (p HlsdateParser) parserofstart(item *m3u8.SegmentItem, s string, n string, start string, sts *[]util.St) bool {
	if n == s && start >= s {
		starttime, _ := time.Parse(layout, start)
		ntime, _ := time.Parse(layout, n)
		ss := starttime.Unix() - ntime.Unix()
		if ss > 0 && ss <= int64(item.Duration) {
			st := util.St{
				Filename: item.Segment,
				Ss:       util.ToTimeStr(ss),
				To:       util.ToTimeStr(int64(item.Duration) - ss),
				Duration: int64(item.Duration),
			}
			*sts = append(*sts, st)
			return true
		}
	}
	return false
}

func (p HlsdateParser) Parser(start string, end string) ([]util.St, error) {

	if start >= end {
		return nil, errors.New("结束时间必须大于开始时间")
	}

	segofstart := p.findSeg(start)
	if segofstart == nil {
		return nil, errors.New("解析异常")
	}

	segofend := p.findSeg(end)
	if segofend == nil {
		return nil, errors.New("解析异常")
	}

	ss := path.Base(segofstart.Segment)[8:22]
	ee := path.Base(segofend.Segment)[8:22]

	logger.Info(ss, ee)

	sts := make([]util.St, 0)
	for _, item := range p.playlist.Segments() {
		nn := path.Base(item.Segment)[8:22]
		if nn >= ss && nn <= ee {

			if p.parserofend(item, ss, ee, nn, start, end, &sts) {
				break
			}

			if p.parserofstart(item, ss, nn, start, &sts) {
				continue
			}

			st := util.St{
				Filename: item.Segment,
				Ss:       util.ToTimeStr(0),
				To:       util.ToTimeStr(int64(item.Duration)),
				Duration: int64(item.Duration),
			}

			sts = append(sts, st)

		}
	}

	return sts, nil
}
