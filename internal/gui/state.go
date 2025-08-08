package gui

import (
	"image"

	"gioui.org/widget"
	"gioui.org/x/explorer"
	"github.com/kvitebjorn/gol/internal/game"
)

type viewCache struct {
	img      *image.RGBA
	turn     int
	panX     int
	panY     int
	zoom     float64
	width    int
	height   int
	cellSize int
}

var (
	// Pan and zoom
	panX, panY int
	zoomLevel  float64 = 1.0

	// Playback state
	playing    bool
	paused     bool
	playStopCh chan struct{}

	// Game state
	gameState    game.Game
	initialBoard game.InfiniteGrid

	// Board clickable tag
	boardTag       = new(bool)
	boardClickable = widget.Clickable{}

	// File dialog related
	fileReadErr      error
	fileDialogActive bool

	// ONLY INITIALIZE THIS ONCE!!!
	explorerInstance *explorer.Explorer
)

func stopPlayback() {
	if playing {
		if playStopCh != nil {
			close(playStopCh)
			playStopCh = nil
		}
		playing = false
		paused = false
	}
}
