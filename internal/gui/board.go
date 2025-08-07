package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/widget"
)

var boardClickable widget.Clickable
var boardTag = new(bool)

// GUI-level cache for the *current visible view* (simple whole-view cache).
// This is invalidated when the board changes generation, pan/zoom changes, or when import/reset happens.
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

func computeDynamicView(gtx layout.Context,
	zoom float64,
	panX, panY int) (
	minRow,
	minCol,
	maxRow,
	maxCol,
	cellSize,
	margin,
	width,
	height int) {
	availableWidth := gtx.Constraints.Max.X
	availableHeight := gtx.Constraints.Max.Y

	// Determine desired cell size
	cellSizeF := zoom * 20 // base size of 20px at 1.0 zoom
	cellSize = int(cellSizeF)
	if cellSize > 50 {
		cellSize = 50
	}
	if cellSize < 2 {
		cellSize = 2
	}

	// Compute how many cells fit in the view
	cols := availableWidth / cellSize
	rows := availableHeight / cellSize

	// PanX/Y are cell offsets — apply them to center
	centerRow := panY
	centerCol := panX

	minRow = centerRow - rows/2
	minCol = centerCol - cols/2
	maxRow = minRow + rows
	maxCol = minCol + cols

	margin = 0
	width = cols*cellSize + 2*margin
	height = rows*cellSize + 2*margin

	return
}
