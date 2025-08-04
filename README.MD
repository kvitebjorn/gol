## Conway's Game of Life in Go

This project is a high-performance implementation of Conway's Game of Life in Go, featuring:

- Double-buffered board updates for efficiency
- 1D slice-based grid for cache-friendly access
- Automatic parallelization for large grids
- Test suite for common patterns (still lifes, oscillators, spaceships)

### Usage

Run the program:

```sh
go run main.go
```

You can modify the initial pattern in `main.go`.
TODO: accept patterns via cmd line

### Testing

Run all tests:

```sh
go test
```

### Features

- Efficient double-buffered updates (no per-tick allocations)
- Parallelized updates for large grids (uses all CPU cores)
- 1D slice for grid storage (improves cache locality)
- Comprehensive test suite for classic Game of Life patterns

### License

MIT License
