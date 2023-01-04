package util

import (
	"fmt"
	"math"
)

func ToTimeStr(seconds int64) string {
	if seconds == 0 {
		return "00:00:00"
	}
	hh := 0
	if seconds >= 3600 {
		hh = int(math.Floor(float64(seconds) / float64(3600)))
		if hh == 24 {
			hh = 0
		}
	}
	mm := 0
	seconds = seconds % 3600
	if seconds >= 60 {
		mm = int(math.Floor(float64(seconds) / float64(60)))
	}
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, seconds)
}

type St struct {
	Filename string
	Ss       string
	To       string
	Duration int64
}
