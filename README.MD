## Conway's Game of Life in Go

This project is a high-performance implementation of Conway's Game of Life in Go, featuring:

- A sparse map-based "infinite" grid that can grow dynamically during runtime
- NVIDIA GPU support, falls back to CPU if unavailable!
- A minimal GUI
  - Zoom with -/+ or the mouse wheel
  - Pan with arrow keys or secondary button drag
  - Edit cells with a primary button click
  - 20~ fps on large patterns on a decent machine (CPU)
  - 40~ fps on large patterns on a decent machine (GPU)
- RLE support 
- Test suite for common patterns (still lifes, oscillators, spaceships)

![Screenshot](assets/gol.png)

## Automatic GPU detection
`gol` detects an NVIDIA GPU at runtime:

- If an NVIDIA driver + at least one CUDA-capable GPU is present, GPU mode is used.
- Otherwise, CPU mode is used.

```bash
go build
```

If you have CUDA installed, you need to export LD_LIBRARY_PATH before building

```
export LD_LIBRARY_PATH=/home/kyle/dev/github.com/kvitebjorn/gol/internal/gpu:$LD_LIBRARY_PATH
```

This requires the CUDA toolkit and dev libraries installed. See the `gpu` package README for further details.

### License

MIT License
