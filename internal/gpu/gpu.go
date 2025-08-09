package gpu

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lgpu
#include "gpu.h"
*/
import "C"
import (
	"unsafe"

	"github.com/kvitebjorn/gol/internal/board"
)

// Initial demo only to get gpu stuff hooked up, not used
func Square(a []float32) {
	C.square((*C.float)(&a[0]), C.int(len(a)))
}

// Our Game of Life rules applied
func Tick(src, dst *board.InfiniteGrid) {
	sminR, sminC, smaxR, smaxC := src.Bounds()
	if len(src.Cells) == 0 {
		dst.Cells = make(map[[2]int]board.Cell)
		dst.BoundsValid = false
		return
	}

	pad := 1 // extra cells around the bounding box for births

	// Expand bounds by pad in each direction
	sminR -= pad
	sminC -= pad
	smaxR += pad
	smaxC += pad

	rows := smaxR - sminR + 1
	cols := smaxC - sminC + 1

	n := rows * cols
	srcFlat := make([]C.int, n)

	// Fill flat src array - easier to work with here, and probably faster than 2d
	for coord := range src.Cells {
		r := coord[0] - sminR // shift by padded min
		c := coord[1] - sminC
		if r < 0 || r >= rows || c < 0 || c >= cols {
			continue
		}
		srcFlat[r*cols+c] = 1
	}

	dstFlat := make([]C.int, n)

	// Do the thing on the GPU
	C.tick(
		(*C.int)(unsafe.Pointer(&srcFlat[0])),
		(*C.int)(unsafe.Pointer(&dstFlat[0])),
		C.int(rows),
		C.int(cols),
	)

	// Map the GPU memory back into our host board structure
	dst.Cells = make(map[[2]int]board.Cell)
	dst.BoundsValid = false

	for i := 0; i < n; i++ {
		if dstFlat[i] != 0 {
			r := i / cols
			c := i % cols
			gr := r + sminR
			gc := c + sminC
			dst.Set(gr, gc, true)
		}
	}

	if len(dst.Cells) == 0 {
		dst.BoundsValid = true
		dst.MinRow, dst.MaxRow, dst.MinCol, dst.MaxCol = 0, 0, 0, 0
	}
}

func GPUAvailable() bool {
	return C.gpu_available() > 0
}
