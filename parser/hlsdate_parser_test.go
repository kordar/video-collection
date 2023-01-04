package parser

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestHlsdateParser_Parser(t *testing.T) {
	parser := NewHlsdateParser("/Users/mac/Movies/ddd/stream.m3u8")
	sts, err := parser.Parser("20230103095944", "20230103120044")
	log.Println(err)
	if err == nil {
		for _, st := range sts {
			log.Printf("%+v", st)
		}
	}
}
