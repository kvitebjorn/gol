#include <stdio.h>
#include <cuda.h>

// Device (kernel)
__global__ void square_cuda(float *a, int N)
{
  int idx = blockIdx.x * blockDim.x + threadIdx.x;
  if (idx < N)
    a[idx] = a[idx] * a[idx];
}

extern "C"
{
  // Host driver - implements gpu.h
  void square(float *a, int N)
  {
    float *a_d;
    size_t size = N * sizeof(float);

    // Allocate memory on the GPU
    cudaMalloc((void **)&a_d, size);

    // Copy the input data from CPU memory to GPU memory
    cudaMemcpy(a_d, a, size, cudaMemcpyHostToDevice);

    // Launch the GPU kernel to do work
    int block_size = 4;
    int n_blocks = N / block_size + (N % block_size == 0 ? 0 : 1);
    square_cuda<<<n_blocks, block_size>>>(a_d, N);

    // Copy the result data from GPU memory back to our CPU memory
    cudaMemcpy(a, a_d, size, cudaMemcpyDeviceToHost);

    // Free the GPU memory
    cudaFree(a_d);
  }
}
