package gui

import "gioui.org/widget"

var (
	zoomLevel = 1.0 // multiplier on base cell size
	panX      = 0   // in cells
	panY      = 0   // in cells
)
var nextButton widget.Clickable
var resetButton widget.Clickable
var importButton widget.Clickable
var playPauseButton widget.Clickable
