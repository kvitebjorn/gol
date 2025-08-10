// cudart_loader.c
#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include "cudart_loader.h"

#ifdef __cplusplus
extern "C"
{
#endif

  // Function pointer typedefs (match CUDA runtime API signatures)
  typedef int (*cudaMalloc_t)(void **, size_t);
  typedef int (*cudaFree_t)(void *);
  typedef int (*cudaMemcpy_t)(void *, const void *, size_t, int); // cudaMemcpyKind is int-like
  typedef int (*cudaDeviceSynchronize_t)(void);
  typedef const char *(*cudaGetErrorString_t)(int);

  static void *cudart_handle = NULL;
  static cudaMalloc_t p_cudaMalloc = NULL;
  static cudaFree_t p_cudaFree = NULL;
  static cudaMemcpy_t p_cudaMemcpy = NULL;
  static cudaDeviceSynchronize_t p_cudaDeviceSynchronize = NULL;
  static cudaGetErrorString_t p_cudaGetErrorString = NULL;

  // Try candidate library names commonly present on systems.
  // Return 1 on success (loaded), 0 on failure.
  static int try_open_cudart(void)
  {
    const char *names[] = {
        "libcudart.so",
        "libcudart.so.13.0",
        "libcudart.so.12.0",
        "libcudart.so.11.0",
        "libcudart.so.10.1",
        NULL};
    for (const char **p = names; *p; ++p)
    {
      void *h = dlopen(*p, RTLD_LAZY | RTLD_LOCAL);
      if (h)
      {
        cudart_handle = h;
        return 1;
      }
    }
    return 0;
  }

  int loadCUDARuntime(void)
  {
    if (cudart_handle)
      return 1; // already loaded

    if (!try_open_cudart())
    {
      // not present
      return 0;
    }

    p_cudaMalloc = (cudaMalloc_t)dlsym(cudart_handle, "cudaMalloc");
    p_cudaFree = (cudaFree_t)dlsym(cudart_handle, "cudaFree");
    p_cudaMemcpy = (cudaMemcpy_t)dlsym(cudart_handle, "cudaMemcpy");
    p_cudaDeviceSynchronize = (cudaDeviceSynchronize_t)dlsym(cudart_handle, "cudaDeviceSynchronize");
    p_cudaGetErrorString = (cudaGetErrorString_t)dlsym(cudart_handle, "cudaGetErrorString");

    if (!p_cudaMalloc || !p_cudaFree || !p_cudaMemcpy || !p_cudaDeviceSynchronize || !p_cudaGetErrorString)
    {
      dlclose(cudart_handle);
      cudart_handle = NULL;
      p_cudaMalloc = NULL;
      p_cudaFree = NULL;
      p_cudaMemcpy = NULL;
      p_cudaDeviceSynchronize = NULL;
      p_cudaGetErrorString = NULL;
      return 0;
    }

    return 1;
  }

  void unloadCUDARuntime(void)
  {
    if (cudart_handle)
    {
      dlclose(cudart_handle);
      cudart_handle = NULL;
    }
    p_cudaMalloc = NULL;
    p_cudaFree = NULL;
    p_cudaMemcpy = NULL;
    p_cudaDeviceSynchronize = NULL;
    p_cudaGetErrorString = NULL;
  }

  // The wrappers return non-zero if cuda runtime couldn't be loaded or call failed.
  // When runtime is not present they return '1' (non-zero) to indicate error.
  cudaError_t cudaMalloc_wrap(void **devPtr, size_t size)
  {
    if (!p_cudaMalloc)
    {
      if (!loadCUDARuntime())
        return 1;
    }
    return p_cudaMalloc(devPtr, size);
  }

  cudaError_t cudaFree_wrap(void *devPtr)
  {
    if (!p_cudaFree)
    {
      if (!loadCUDARuntime())
        return 1;
    }
    return p_cudaFree(devPtr);
  }

  cudaError_t cudaMemcpy_wrap(void *dst, const void *src, size_t count, cudaMemcpyKind kind)
  {
    if (!p_cudaMemcpy)
    {
      if (!loadCUDARuntime())
        return 1;
    }
    return p_cudaMemcpy(dst, src, count, (int)kind);
  }

  cudaError_t cudaDeviceSynchronize_wrap(void)
  {
    if (!p_cudaDeviceSynchronize)
    {
      if (!loadCUDARuntime())
        return 1;
    }
    return p_cudaDeviceSynchronize();
  }

  const char *cudaGetErrorString_wrap(cudaError_t error)
  {
    if (!p_cudaGetErrorString)
    {
      if (!loadCUDARuntime())
        return "CUDA runtime not loaded";
    }
    return p_cudaGetErrorString(error);
  }

#ifdef __cplusplus
} // extern "C"
#endif
