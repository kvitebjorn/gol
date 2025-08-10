#ifndef CUDART_LOADER_H
#define CUDART_LOADER_H

#ifdef __cplusplus
extern "C"
{
#endif

#include <stddef.h> // for size_t

#ifndef __CUDA_RUNTIME_API_H__
  typedef int cudaError_t;
  typedef int cudaMemcpyKind;
#endif

  int loadCUDARuntime(void);
  void unloadCUDARuntime(void);

  cudaError_t cudaMalloc_wrap(void **devPtr, size_t size);
  cudaError_t cudaFree_wrap(void *devPtr);
  cudaError_t cudaMemcpy_wrap(void *dst, const void *src, size_t count, cudaMemcpyKind kind);
  cudaError_t cudaDeviceSynchronize_wrap(void);
  const char *cudaGetErrorString_wrap(cudaError_t error);

#ifdef __cplusplus
}
#endif

#endif // CUDART_LOADER_H
