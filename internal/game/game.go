package game

import (
	"github.com/kvitebjorn/gol/internal/board"
	"github.com/kvitebjorn/gol/internal/gpu"
)

// Game holds the state of the Game of Life, including double-buffered boards.
type Game struct {
	// BoardA and BoardB are used for double buffering whereby one Board is the current board,
	// and the other is the next generation. This allows for efficient updates without needing
	// to copy the entire grid each generation.
	BoardA board.InfiniteGrid
	BoardB board.InfiniteGrid

	// UseA indicates which board is currently active.
	UseA bool

	// Turn is the current turn number.
	Turn int
}

func (g *Game) CurrentBoard() *board.InfiniteGrid {
	if g.UseA {
		return &g.BoardA
	}
	return &g.BoardB
}

// Tick advances the game by one generation, applying the Game of Life rules.
func (g *Game) Tick() {
	// If BoardA is InfiniteGrid, use infinite tick logic
	// keeping this here in case we implement a bitboard and want to toggle board implementations
	g.TickInfinite()
}

// TickInfinite advances the game by one generation for InfiniteGrid.
func (g *Game) TickInfinite() {
	var src, dst *board.InfiniteGrid
	if g.UseA {
		src = &g.BoardA
		dst = &g.BoardB
	} else {
		src = &g.BoardB
		dst = &g.BoardA
	}

	if UseGpu {
		g.TickGpu(src, dst)
	} else {
		g.TickCpu(src, dst)
	}

	g.UseA = !g.UseA
	g.Turn++
}

func (g *Game) TickGpu(src, dst *board.InfiniteGrid) {
	gpu.Tick(src, dst)
}

func (g *Game) TickCpu(src, dst *board.InfiniteGrid) {
	// Clear destination
	dst.Cells = make(map[[2]int]board.Cell)
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
		var next board.Cell
		if alive {
			next = count == 2 || count == 3
		} else {
			next = count == 3
		}
		if next {
			dst.Cells[pos] = true
		}
	}
}
