//go:build linux
// +build linux

package memory

import (
	"encoding/binary"
	"fmt"
	"os"

	mmap "github.com/edsrzf/mmap-go"
)

const (
	FrameSize    = 1024
	TargetOffset = 0x2000000
	PageSize     = 4096
)

func ReadFrame(path string, freq float64, data *[]float64) error {
	if path == "" {
		path = "/dev/mem"
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return fmt.Errorf("open %s failed: %w", path, err)
	}
	defer file.Close()

	alignedOffset := TargetOffset & ^(PageSize - 1)
	offsetInPage := TargetOffset - alignedOffset
	length := FrameSize * 2 // 2 bytes per uint16

	mem, err := mmap.MapRegion(file, offsetInPage+length, mmap.RDONLY, 0, int64(alignedOffset))
	if err != nil {
		return fmt.Errorf("mmap failed: %w", err)
	}
	defer mem.Unmap()

	values := make([]float64, FrameSize)
	for i := 0; i < FrameSize; i++ {
		raw := binary.LittleEndian.Uint16(mem[offsetInPage+i*2 : offsetInPage+i*2+2])
		values[i] = float64(raw)
	}

	*data = append(*data, values...)
	return nil
}
