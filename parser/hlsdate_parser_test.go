package parser

import (
	logger "github.com/kordar/gologger"
	"testing"
)

func TestHlsdateParser_Parser(t *testing.T) {
	parser := NewHlsdateParser("/Users/mac/Movies/ddd/stream.m3u8")
	sts, err := parser.Parser("20230103095944", "20230103120044")
	logger.Warn(err)
	if err == nil {
		for _, st := range sts {
			logger.Infof("%+v", st)
		}
	}
}
