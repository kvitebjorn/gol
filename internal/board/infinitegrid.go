package board

import "maps"

// Cell represents a single cell in the Game of Life grid.
type Cell bool

// InfiniteGrid represents a sparse, infinite board using a map.
type InfiniteGrid struct {
	Cells                          map[[2]int]Cell // key: [row, col], value: alive/dead
	MinRow, MinCol, MaxRow, MaxCol int
	BoundsValid                    bool
}

func NewInfiniteGrid() InfiniteGrid {
	return InfiniteGrid{Cells: make(map[[2]int]Cell), BoundsValid: false}
}

func (g *InfiniteGrid) At(row, col int) Cell {
	return g.Cells[[2]int{row, col}]
}

func (g *InfiniteGrid) Set(row, col int, val Cell) {
	key := [2]int{row, col}
	if val {
		if _, exists := g.Cells[key]; !exists {
			g.Cells[key] = true
			if g.BoundsValid {
				if row < g.MinRow {
					g.MinRow = row
				}
				if row > g.MaxRow {
					g.MaxRow = row
				}
				if col < g.MinCol {
					g.MinCol = col
				}
				if col > g.MaxCol {
					g.MaxCol = col
				}
			}
		}
	} else {
		if _, exists := g.Cells[key]; exists {
			delete(g.Cells, key)
			// Removing a cell may invalidate bounds
			g.BoundsValid = false
		}
	}
}

func (g *InfiniteGrid) Bounds() (minRow, minCol, maxRow, maxCol int) {
	if g.BoundsValid {
		return g.MinRow, g.MinCol, g.MaxRow, g.MaxCol
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
		g.MinRow, g.MinCol, g.MaxRow, g.MaxCol = minRow, minCol, maxRow, maxCol
		g.BoundsValid = true
		return minRow, minCol, maxRow, maxCol
	}
	// No live cells
	g.BoundsValid = true
	g.MinRow, g.MaxRow, g.MinCol, g.MaxCol = 0, 0, 0, 0
	return 0, 0, 0, 0
}

// DeepCopy returns a deep copy of the InfiniteGrid.
func (g InfiniteGrid) DeepCopy() InfiniteGrid {
	copy := NewInfiniteGrid()
	maps.Copy(copy.Cells, g.Cells)
	copy.MinRow = g.MinRow
	copy.MaxRow = g.MaxRow
	copy.MinCol = g.MinCol
	copy.MaxCol = g.MaxCol
	copy.BoundsValid = g.BoundsValid
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
