package game

import "maps"

// Cell represents a single cell in the Game of Life grid.
type Cell bool

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
	maps.Copy(copy.Cells, g.Cells)
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

func (g *InfiniteGrid) AliveCellsWithinBounds(minCol, minRow, maxCol, maxRow int) [][2]int {
	out := make([][2]int, 0)
	for k := range g.Cells {
		r, c := k[0], k[1]
		if r < minRow || r >= maxRow || c < minCol || c >= maxCol {
			continue
		}
		out = append(out, k)
	}
	return out
}
