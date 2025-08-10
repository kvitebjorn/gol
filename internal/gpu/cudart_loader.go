package gpu

/*
#cgo LDFLAGS: -ldl
#include "cudart_loader.h"
*/
import "C"
import "unsafe"

// Go wrappers around the C _wrap functions.

func CudaMalloc(devPtr *unsafe.Pointer, size uint64) C.cudaError_t {
	return C.cudaMalloc_wrap(devPtr, C.size_t(size))
}

func CudaFree(devPtr unsafe.Pointer) C.cudaError_t {
	return C.cudaFree_wrap(devPtr)
}

func CudaMemcpy(dst unsafe.Pointer, src unsafe.Pointer, count uint64, kind int) C.cudaError_t {
	return C.cudaMemcpy_wrap(dst, src, C.size_t(count), C.cudaMemcpyKind(kind))
}

func CudaDeviceSynchronize() C.cudaError_t {
	return C.cudaDeviceSynchronize_wrap()
}

func CudaGetErrorString(err C.cudaError_t) string {
	return C.GoString(C.cudaGetErrorString_wrap(err))
}
