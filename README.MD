## Conway's Game of Life in Go

This project is a high-performance implementation of Conway's Game of Life in Go, featuring:

- Double-buffered board updates for efficiency
- An "infinite" grid that can grow dynamically during runtime
- A minimal GUI
- Import/Export RLE support
- Test suite for common patterns (still lifes, oscillators, spaceships)

### Usage

Run the program:

```sh
go run main.go
```

### Testing

Run all tests:

```sh
go test
```

### Features

- Efficient double-buffered updates (no per-tick allocations)
- An "infinite" grid that can grow dynamically during runtime
- A minimal GUI
- Import/Export RLE support
- Comprehensive test suite for classic Game of Life patterns

### License

MIT License
