//go:build linux
// +build linux

package memory

import (
	"fmt"
	"os"
	"unsafe"

	mmap "github.com/edsrzf/mmap-go"
)

const (
	FrameSize    = 1024      // Размер фрейма в 16-битных значениях
	TargetOffset = 0x2000000 // Физический адрес (пример для FPGA)
	PageSize     = 4096      // Размер страницы памяти
)

func ReadFrame(path string, freq float64, data *[]float64) error {
	if path == "" {
		path = "/dev/mem"
	}

	file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return fmt.Errorf("open %s failed: %w", path, err)
	}
	defer file.Close()

	// 2. Выравнивание параметров
	alignedOffset := TargetOffset & ^(uintptr(PageSize) - 1)
	offsetInPage := TargetOffset - alignedOffset
	length := FrameSize * 2 // 2 байта на значение

	// 3. Отображение памяти через mmap-go
	mem, err := mmap.MapRegion(
		file,
		int(offsetInPage)+length, // Общая длина
		mmap.RDONLY,
		mmap.ANON,
		int64(alignedOffset), // Выровненное смещение
	)
	if err != nil {
		return fmt.Errorf("mmap failed: %w", err)
	}
	defer mem.Unmap()

	values := make([]float64, FrameSize)
	for i := 0; i < FrameSize; i++ {
		raw := *(*uint16)(unsafe.Pointer(&mem[i*2]))
		values[i] = float64(raw)
	}

	*data = append(*data, values...)
	return nil
}
