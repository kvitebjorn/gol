package main

// Helper to create an InfiniteGrid from [][]bool, with (0,0) at upper-left
func makeInfiniteGrid(pattern [][]bool) InfiniteGrid {
	grid := NewInfiniteGrid()
	for i := range pattern {
		for j := range pattern[i] {
			if pattern[i][j] {
				grid.Set(i, j, true)
			}
		}
	}
	return grid
}

// Helper to compare two InfiniteGrids in a given region
func gridsEqualRegion(a, b InfiniteGrid, minRow, minCol, maxRow, maxCol int) bool {
	for i := minRow; i <= maxRow; i++ {
		for j := minCol; j <= maxCol; j++ {
			if a.At(i, j) != b.At(i, j) {
				return false
			}
		}
	}
	return true
}
