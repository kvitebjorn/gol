# Name of the shared library
LIB_NAME = libgpu.so
CUDA_LIB_PATH := /usr/local/cuda-13.0/targets/x86_64-linux/lib

# Compiler flags
NVCC_FLAGS = --ptxas-options=-v --compiler-options '-fPIC'
GPP_FLAGS = -fPIC -shared

all: $(LIB_NAME)

$(LIB_NAME): gpu.cu gpu.h
	nvcc $(NVCC_FLAGS) -c gpu.cu -o gpu.o
	g++ $(GPP_FLAGS) gpu.o gpu.h -o $(LIB_NAME) -L$(CUDA_LIB_PATH) -lcudart

clean:
	rm -f *.so *.o *~
