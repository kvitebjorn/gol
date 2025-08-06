package main

import (
	"flag"
	"fmt"
	"os"
)

// InfiniteGrid represents a sparse, infinite board using a map.
type InfiniteGrid struct {
	Cells                          map[[2]int]Cell // key: [row, col], value: alive/dead
	minRow, minCol, maxRow, maxCol int
	boundsValid                    bool
}

func NewInfiniteGrid() InfiniteGrid {
	return InfiniteGrid{Cells: make(map[[2]int]Cell), boundsValid: false}
}

func (g *InfiniteGrid) At(row, col int) Cell {
	return g.Cells[[2]int{row, col}]
}

func (g *InfiniteGrid) Set(row, col int, val Cell) {
	key := [2]int{row, col}
	if val {
		if _, exists := g.Cells[key]; !exists {
			g.Cells[key] = true
			if g.boundsValid {
				if row < g.minRow {
					g.minRow = row
				}
				if row > g.maxRow {
					g.maxRow = row
				}
				if col < g.minCol {
					g.minCol = col
				}
				if col > g.maxCol {
					g.maxCol = col
				}
			}
		}
	} else {
		if _, exists := g.Cells[key]; exists {
			delete(g.Cells, key)
			// Removing a cell may invalidate bounds
			g.boundsValid = false
		}
	}
}

func (g *InfiniteGrid) Bounds() (minRow, minCol, maxRow, maxCol int) {
	if g.boundsValid {
		return g.minRow, g.minCol, g.maxRow, g.maxCol
	}
	first := true
	for k := range g.Cells {
		r, c := k[0], k[1]
		if first {
			minRow, maxRow = r, r
			minCol, maxCol = c, c
			first = false
		} else {
			if r < minRow {
				minRow = r
			}
			if r > maxRow {
				maxRow = r
			}
			if c < minCol {
				minCol = c
			}
			if c > maxCol {
				maxCol = c
			}
		}
	}
	if !first {
		g.minRow, g.minCol, g.maxRow, g.maxCol = minRow, minCol, maxRow, maxCol
		g.boundsValid = true
		return minRow, minCol, maxRow, maxCol
	}
	// No live cells
	g.boundsValid = true
	g.minRow, g.maxRow, g.minCol, g.maxCol = 0, 0, 0, 0
	return 0, 0, 0, 0
}

// DeepCopy returns a deep copy of the InfiniteGrid.
func (g InfiniteGrid) DeepCopy() InfiniteGrid {
	copy := NewInfiniteGrid()
	for k, v := range g.Cells {
		copy.Cells[k] = v
	}
	copy.minRow = g.minRow
	copy.maxRow = g.maxRow
	copy.minCol = g.minCol
	copy.maxCol = g.maxCol
	copy.boundsValid = g.boundsValid
	return copy
}

// AliveCells returns a slice of coordinates of all currently alive cells.
// This provides a fast sparse iteration path for rendering and other ops.
func (g *InfiniteGrid) AliveCells() [][2]int {
	out := make([][2]int, 0, len(g.Cells))
	for k := range g.Cells {
		out = append(out, k)
	}
	return out
}

// Cell represents a single cell in the Game of Life grid.
type Cell bool

// Game holds the state of the Game of Life, including double-buffered boards.
type Game struct {
	// BoardA and BoardB are used for double buffering whereby one Board is the current board,
	// and the other is the next generation. This allows for efficient updates without needing
	// to copy the entire grid each generation.
	BoardA InfiniteGrid
	BoardB InfiniteGrid
	// UseA indicates which board is currently active.
	UseA bool
	// Turn is the current turn number.
	Turn int
}

// Tick advances the game by one generation, applying the Game of Life rules.
func (g *Game) Tick() {
	// If BoardA is InfiniteGrid, use infinite tick logic
	g.TickInfinite()
}

// TickInfinite advances the game by one generation for InfiniteGrid.
func (g *Game) TickInfinite() {
	var src, dst *InfiniteGrid
	if g.UseA {
		src = &g.BoardA
		dst = &g.BoardB
	} else {
		src = &g.BoardB
		dst = &g.BoardA
	}

	// Clear destination
	dst.Cells = make(map[[2]int]Cell)
	neighborCounts := make(map[[2]int]int)

	// Count neighbors for all live cells and their neighbors
	for pos := range src.Cells {
		r, c := pos[0], pos[1]
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				if dr == 0 && dc == 0 {
					continue
				}
				npos := [2]int{r + dr, c + dc}
				neighborCounts[npos]++
			}
		}
	}
	// Apply rules
	for pos, count := range neighborCounts {
		alive := src.Cells[pos]
		var next Cell
		if alive {
			next = count == 2 || count == 3
		} else {
			next = count == 3
		}
		if next {
			dst.Cells[pos] = true
		}
	}
	g.UseA = !g.UseA
	g.Turn++
}

func main() {
	rleFile := flag.String("rle", "", "Path to RLE file to import as initial pattern")
	flag.Parse()

	var imported *InfiniteGrid
	if *rleFile != "" {
		f, err := os.Open(*rleFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open RLE file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		b, err := ImportRLE(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to import RLE: %v\n", err)
			os.Exit(1)
		}
		imported = &b
	}
	RunGUI(imported)
}
