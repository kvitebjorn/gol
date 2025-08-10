package gpu

// #cgo LDFLAGS: -ldl
// #include <stdlib.h>
// #include <dlfcn.h>
// #include <stdint.h>
//
// typedef int (*cuInit_t)(unsigned int);
// typedef int (*cuDeviceGetCount_t)(int*);
//
// // Return 1 if CUDA driver + at least one device is present, 0 otherwise.
// static int detect_cuda() {
//     void *h = dlopen("libcuda.so", RTLD_LAZY | RTLD_LOCAL);
//     if (!h) {
//         // Some systems have versioned library names (e.g. libcuda.so.1).
//         h = dlopen("libcuda.so.1", RTLD_LAZY | RTLD_LOCAL);
//         if (!h) {
//             return 0;
//         }
//     }
//
//     cuInit_t p_cuInit = (cuInit_t)dlsym(h, "cuInit");
//     cuDeviceGetCount_t p_cuDeviceGetCount = (cuDeviceGetCount_t)dlsym(h, "cuDeviceGetCount");
//
//     if (!p_cuInit || !p_cuDeviceGetCount) {
//         dlclose(h);
//         return 0;
//     }
//
//     if (p_cuInit(0) != 0) {
//         dlclose(h);
//         return 0;
//     }
//
//     int count = 0;
//     if (p_cuDeviceGetCount(&count) != 0) {
//         dlclose(h);
//         return 0;
//     }
//
//     dlclose(h);
//     return count > 0 ? 1 : 0;
// }
import "C"

// HasCUDA returns true if libcuda is present and reports at least one device.
// This does not require linking to any CUDA library at build time.
func HasCUDA() bool {
	ret := C.detect_cuda()
	return ret == 1
}
