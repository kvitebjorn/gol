CUDA_PATH ?= /usr/local/cuda-13.0
CUDA_INC := $(CUDA_PATH)/include
CUDA_LIB := $(CUDA_PATH)/lib64

NVCC_FLAGS = --ptxas-options=-v --compiler-options '-fPIC' -I$(CUDA_INC)
CFLAGS = -fPIC -I$(CUDA_INC) -ldl
GPP_FLAGS = -fPIC -shared

CU_SRC = gpu.cu
C_SRC = cudart_loader.c
OBJS = gpu.o cudart_loader.o
LIB_NAME = libgpu.so

all: $(LIB_NAME)

gpu.o: $(CU_SRC) gpu.h cudart_loader.h
	nvcc $(NVCC_FLAGS) -c $(CU_SRC) -o gpu.o

cudart_loader.o: $(C_SRC) cudart_loader.h
	gcc $(CFLAGS) -c $(C_SRC) -o cudart_loader.o

$(LIB_NAME): $(OBJS)
	nvcc -shared -o $(LIB_NAME) $(OBJS) -ldl

clean:
	rm -f *.o $(LIB_NAME)
