# GPU dev notes
`gpu` assumes that the CUDA Toolkit is installed.
At the time of this writing, the latest is 13-0.

```
wget https://developer.download.nvidia.com/compute/cuda/13.0.0/local_installers/cuda-repo-debian12-13-0-local_13.0.0-580.65.06-1_amd64.deb
```
```
sudo dpkg -i cuda-repo-debian12-13-0-local_13.0.0-580.65.06-1_amd64.deb
```
```
sudo cp /var/cuda-repo-debian12-13-0-local/cuda-*-keyring.gpg /usr/share/keyrings/
```
```
sudo apt-get update
```
```
sudo apt-get -y install cuda-toolkit-13-0
```

Then follow the post-installation actions here: https://docs.nvidia.com/cuda/cuda-installation-guide-linux/index.html#post-installation-actions


# Building with the CUDA wrapper

## In this dir: 

Make sure `CUDA_LIB_PATH` is correct in `Makefile`, then run

```
make
```


## In the project dir:

```
export LD_LIBRARY_PATH=/usr/local/cuda-13.0/targets/x86_64-linux/lib:$LD_LIBRARY_PATH
```

```
export CGO_LDFLAGS="-L/usr/local/cuda-13.0/targets/x86_64-linux/lib -lcudart"
```

```
go build
```
