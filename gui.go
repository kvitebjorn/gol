package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
)

type C = layout.Context
type D = layout.Dimensions

var (
	game         Game         // Global game state
	initialBoard InfiniteGrid // For reset
)

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

var nextButton widget.Clickable
var resetButton widget.Clickable
var importButton widget.Clickable
var playPauseButton widget.Clickable

var explorerInstance *explorer.Explorer
var fileReadErr error
var fileDialogActive bool

var playing bool
var paused bool
var playStopCh chan struct{}

var (
	zoomLevel = 1.0 // multiplier on base cell size
	panX      = 0   // in cells
	panY      = 0   // in cells
)

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
	availableWidth := gtx.Constraints.Max.X - 2*gtx.Dp(unit.Dp(10)) // 10px margin on each side
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

	halfRows := rows / 2
	halfCols := cols / 2

	minRow = centerRow - halfRows
	maxRow = centerRow + halfRows
	minCol = centerCol - halfCols
	maxCol = centerCol + halfCols

	margin = 0
	width = cols*cellSize + 2*margin
	height = rows*cellSize + 2*margin

	return
}

func currentBoard(g *Game) *InfiniteGrid {
	if g.UseA {
		return &g.BoardA
	}
	return &g.BoardB
}

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

func draw(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()

	if explorerInstance == nil {
		explorerInstance = explorer.NewExplorer(w)
	}

	var tag = new(bool)
	event.Op(&ops, tag)

	// Event loop
	for {
		e := w.Event()
		explorerInstance.ListenEvents(e)
		switch evt := e.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, evt)

			// Key events
			for {
				ev, ok := gtx.Event(key.Filter{
					Optional: key.ModShift,
				})
				if !ok {
					break
				}
				if kev, ok := ev.(key.Event); ok {
					switch kev.Name {
					case key.NameUpArrow:
						panY -= 2
						w.Invalidate()
					case key.NameDownArrow:
						panY += 2
						w.Invalidate()
					case key.NameLeftArrow:
						panX -= 2
						w.Invalidate()
					case key.NameRightArrow:
						panX += 2
						w.Invalidate()
					case "+":
						zoomLevel *= 1.1
						if zoomLevel > 4 {
							zoomLevel = 4
						}
						w.Invalidate()
					case "-":
						zoomLevel *= 0.9
						if zoomLevel < 0.25 {
							zoomLevel = 0.25
						}
						w.Invalidate()
					}
				} else {
					break
				}
			}

			// Button events
			if playPauseButton.Clicked(gtx) {
				if !playing {
					playing = true
					paused = false
					playStopCh = make(chan struct{})
					go func(stopCh chan struct{}, win *app.Window) {
						for {
							select {
							case <-stopCh:
								return
							default:
								if !paused {
									game.Tick()
									win.Invalidate()
									select {
									case <-stopCh:
										return
									case <-time.After(100 * time.Millisecond):
									}
								} else {
									// If paused, poll for stop or unpause every 100ms
									select {
									case <-stopCh:
										return
									case <-time.After(100 * time.Millisecond):
									}
								}
							}
						}
					}(playStopCh, w)
				} else {
					paused = !paused
				}
			}
			if resetButton.Clicked(gtx) {
				stopPlayback()
				game = Game{
					BoardA: initialBoard.DeepCopy(),
					BoardB: initialBoard.DeepCopy(),
					UseA:   true,
					Turn:   1,
				}
				zoomLevel = 1.0
				panX = 0
				panY = 0
				w.Invalidate()
			}
			if nextButton.Clicked(gtx) && (!playing || paused) {
				game.Tick()
				w.Invalidate()
			}
			if importButton.Clicked(gtx) && !fileDialogActive {
				fileDialogActive = true
				go func() {
					r, err := explorerInstance.ChooseFile(".rle")
					if err != nil {
						fileReadErr = err
						fileDialogActive = false
						return
					}
					defer r.Close()
					b, err := ImportRLE(r)
					if err != nil {
						fileReadErr = err
					} else {
						initialBoard = b.DeepCopy()
						fileReadErr = nil
						resetButton.Click()
					}
					fileDialogActive = false
					w.Invalidate()
				}()
			}

			// Layout
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					gen := game.Turn
					label := material.H6(th, "Generation: "+itoa(gen))
					return layout.Center.Layout(gtx, label.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					label := material.Body1(th, fmt.Sprintf("Zoom: %.2fx  Pan: (%d,%d)", zoomLevel, panX, panY))
					return layout.Center.Layout(gtx, label.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					board := currentBoard(&game)
					minRow, minCol, maxRow, maxCol, cellSize, margin, width, height :=
						computeDynamicView(gtx, zoomLevel, panX, panY)
					size := gtx.Constraints.Max
					gtx.Constraints.Min = size
					return layout.Center.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Max.X = width
						gtx.Constraints.Max.Y = height
						for i := minRow; i < maxRow; i++ {
							for j := minCol; j < maxCol; j++ {
								x := margin + (j-minCol)*cellSize
								y := margin + (i-minRow)*cellSize
								col := color.NRGBA{R: 220, G: 220, B: 220, A: 255}
								if board.At(i, j) {
									col = color.NRGBA{R: 0, G: 200, B: 0, A: 255}
								}
								op := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
								paint.FillShape(gtx.Ops, col, clip.Rect{
									Min: image.Pt(0, 0),
									Max: image.Pt(cellSize, cellSize)}.Op())
								op.Pop()
							}
						}
						gridCol := color.NRGBA{R: 180, G: 180, B: 180, A: 255}
						for i := 0; i <= (maxRow - minRow); i++ {
							y := margin + i*cellSize
							op := op.Offset(image.Pt(margin, y)).Push(gtx.Ops)
							paint.FillShape(gtx.Ops, gridCol, clip.Rect{
								Min: image.Pt(0, 0),
								Max: image.Pt((maxCol-minCol)*cellSize, 1)}.Op())
							op.Pop()
						}
						for j := 0; j <= (maxCol - minCol); j++ {
							x := margin + j*cellSize
							op := op.Offset(image.Pt(x, margin)).Push(gtx.Ops)
							paint.FillShape(gtx.Ops, gridCol, clip.Rect{
								Min: image.Pt(0, 0),
								Max: image.Pt(1, (maxRow-minRow)*cellSize)}.Op())
							op.Pop()
						}
						return D{Size: image.Pt(width, height)}
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    unit.Dp(10),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(10),
						Right:  unit.Dp(10),
					}.Layout(gtx, func(gtx C) D {
						return layout.Flex{
							Axis:      layout.Horizontal,
							Spacing:   layout.SpaceSides,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								btn := material.Button(th, &nextButton, "Next")
								if playing && !paused {
									btn.Background = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
								}
								return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
							}),
							layout.Rigid(func(gtx C) D {
								btn := material.Button(th, &playPauseButton, func() string {
									if !playing {
										return "Play"
									} else if paused {
										return "Resume"
									} else {
										return "Pause"
									}
								}())
								return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
							}),
							layout.Rigid(func(gtx C) D {
								btn := material.Button(th, &resetButton, "Reset")
								return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
							}),
							layout.Rigid(func(gtx C) D {
								btn := material.Button(th, &importButton, "Import RLE")
								return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
							}),
						)
					})
				}),
			)
			if fileReadErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to import RLE: %v\n", fileReadErr)
			}
			evt.Frame(gtx.Ops)

		case app.DestroyEvent:
			return evt.Err
		}
	}
}

// RunGUI launches a minimal Game of Life GUI. Right arrow advances one tick.
func RunGUI(imported *InfiniteGrid) {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Game of Life"))
		w.Option(app.Maximized.Option())

		var ig InfiniteGrid
		if imported != nil {
			ig = imported.DeepCopy()
		} else {
			// Default: glider
			initial := [][2]int{
				{0, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2},
			}
			ig = NewInfiniteGrid()
			for _, p := range initial {
				ig.Set(p[0], p[1], true)
			}
		}
		initialBoard = ig.DeepCopy()
		game = Game{
			BoardA: initialBoard.DeepCopy(),
			BoardB: initialBoard.DeepCopy(),
			UseA:   true,
			Turn:   1,
		}
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
