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

### Usage

```sh
go install github.com/kvitebjorn/gol@latest
```

### License

MIT License
