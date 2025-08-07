package gui

import (
	"log"
	"os"

	"gioui.org/app"
	"github.com/kvitebjorn/gol/internal/game"
)

var (
	gameState    game.Game         // Global game state
	initialBoard game.InfiniteGrid // For reset
)

// RunGUI launches a minimal Game of Life GUI. Right arrow advances one tick.
func RunGUI(imported *game.InfiniteGrid) {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Game of Life"))
		w.Option(app.Maximized.Option())

		var ig game.InfiniteGrid
		if imported != nil {
			ig = imported.DeepCopy()
		} else {
			// Default: glider
			initial := [][2]int{
				{0, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2},
			}
			ig = game.NewInfiniteGrid()
			for _, p := range initial {
				ig.Set(p[0], p[1], true)
			}
		}
		initialBoard = ig.DeepCopy()
		gameState = game.Game{
			BoardA: initialBoard.DeepCopy(),
			BoardB: initialBoard.DeepCopy(),
			UseA:   true,
			Turn:   1,
		}
		if err := runWindow(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
