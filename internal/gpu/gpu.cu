#include <stdio.h>
#include <cuda.h>
#include "cudart_loader.h"

// Device (kernel)
__global__ void square_cuda(float *a, int N)
{
  int idx = blockIdx.x * blockDim.x + threadIdx.x;
  if (idx < N)
    a[idx] = a[idx] * a[idx];
}

// Device `tick` (kernel). src and dst are flat row-major int arrays (0 or 1).
__global__ void tick_cuda(const int *src, int *dst, int rows, int cols)
{
  int idx = blockIdx.x * blockDim.x + threadIdx.x;
  int total = rows * cols;
  if (idx >= total)
    return;

  int r = idx / cols;
  int c = idx % cols;

  int alive = 0;

  // iterate neighbors
  for (int dr = -1; dr <= 1; ++dr)
  {
    for (int dc = -1; dc <= 1; ++dc)
    {
      if (dr == 0 && dc == 0)
        continue;
      int rr = r + dr;
      int cc = c + dc;

      // bounds check: treat outside as dead
      if (rr < 0 || rr >= rows || cc < 0 || cc >= cols)
        continue;
      int nidx = rr * cols + cc;
      alive += src[nidx];
    }
  }

  int cur = src[idx];
  int next = 0;
  if (cur == 1)
  {
    // live cell: survives if 2 or 3 neighbors
    next = (alive == 2 || alive == 3) ? 1 : 0;
  }
  else
  {
    // dead cell: becomes alive if exactly 3 neighbors
    next = (alive == 3) ? 1 : 0;
  }
  dst[idx] = next;
}

extern "C"
{
  // Host driver - implements gpu.h `square`
  void square(float *a, int N)
  {
    float *a_d;
    size_t size = N * sizeof(float);

    // Allocate memory on the GPU
    cudaMalloc_wrap((void **)&a_d, size);

    // Copy the input data from CPU memory to GPU memory
    cudaMemcpy_wrap(a_d, a, size, cudaMemcpyHostToDevice);

    // Launch the GPU kernel to do work
    int block_size = 4;
    int n_blocks = N / block_size + (N % block_size == 0 ? 0 : 1);
    square_cuda<<<n_blocks, block_size>>>(a_d, N);

    // Copy the result data from GPU memory back to our CPU memory
    cudaMemcpy_wrap(a, a_d, size, cudaMemcpyDeviceToHost);

    // Free the GPU memory
    cudaFree_wrap(a_d);
  }

  // Host driver - implements gpu.h `tick`
  void tick(int *src, int *dst, int rows, int cols)
  {
    size_t n = (size_t)rows * (size_t)cols;
    if (n == 0)
      return;

    int *src_d = nullptr;
    int *dst_d = nullptr;
    size_t bytes = n * sizeof(int);

    cudaError_t err;

    // Allocate memory on the device for src & dst
    err = cudaMalloc_wrap((void **)&src_d, bytes);
    if (err != cudaSuccess)
    {
      fprintf(stderr, "cudaMalloc src failed: %s\n", cudaGetErrorString_wrap(err));
      return;
    }
    err = cudaMalloc_wrap((void **)&dst_d, bytes);
    if (err != cudaSuccess)
    {
      fprintf(stderr, "cudaMalloc dst failed: %s\n", cudaGetErrorString_wrap(err));
      cudaFree_wrap(src_d);
      return;
    }

    // Copy values to the device for src
    err = cudaMemcpy_wrap(src_d, src, bytes, cudaMemcpyHostToDevice);
    if (err != cudaSuccess)
    {
      fprintf(stderr, "cudaMemcpy to device failed: %s\n", cudaGetErrorString_wrap(err));
      cudaFree_wrap(src_d);
      cudaFree_wrap(dst_d);
      return;
    }

    // Launch the kernel to do the thing
    int block_size = 256;
    int n_blocks = (int)((n + block_size - 1) / block_size);

    tick_cuda<<<n_blocks, block_size>>>(src_d, dst_d, rows, cols);
    cudaDeviceSynchronize_wrap();

    // Copy the device dst to our host
    err = cudaMemcpy_wrap(dst, dst_d, bytes, cudaMemcpyDeviceToHost);
    if (err != cudaSuccess)
    {
      fprintf(stderr, "cudaMemcpy to host failed: %s\n", cudaGetErrorString_wrap(err));
    }

    cudaFree_wrap(src_d);
    cudaFree_wrap(dst_d);
  }
}
