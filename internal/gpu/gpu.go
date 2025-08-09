package gpu

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lgpu
#include "gpu.h"
*/
import "C"

func Square(a []float32) {
	C.square((*C.float)(&a[0]), C.int(len(a)))
}
