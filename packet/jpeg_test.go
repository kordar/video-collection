package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"testing"
)

func TestJpeg(t *testing.T) {

	var pi int
	b := []byte{0x08}
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.LittleEndian, &pi)
	if err != nil {
		log.Fatalln("binary.Read failed:", err)
	}
	fmt.Println(pi)
}
