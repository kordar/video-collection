package parser

import (
	"github.com/q191201771/naza/pkg/nazalog"
	"testing"
)

func TestHlsdateParser_Parser(t *testing.T) {
	parser := NewHlsdateParser("/Users/mac/Movies/ddd/stream.m3u8")
	sts, err := parser.Parser("20230103095944", "20230103120044")
	nazalog.Println(err)
	if err == nil {
		for _, st := range sts {
			nazalog.Printf("%+v", st)
		}
	}
}
