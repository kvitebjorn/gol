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

- If an NVIDIA driver + CUDA runtime + at least one CUDA-capable GPU is present, GPU mode is used.
  - You MUST have CUDA installed!!!
- Otherwise, CPU mode is used.

## Build
In order to build this project, you currently *must* install the CUDA Toolkit and build the `.so`. 
Follow the instructions in `internal/gpu`.

I'm working on making this easier, especially for machines that don't have an NVIDIA GPU or CUDA!

Otherwise, grab a release, it will fall back to the CPU mode :)

### License

MIT License
