package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
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

// For grid memoization (minRow, minCol, maxRow, maxCol, height, width)
var dimensionCache map[[3]int][6]int = make(map[[3]int][6]int)

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

// Move button state to package scope so it persists across frames
var nextButton widget.Clickable
var resetButton widget.Clickable
var importButton widget.Clickable
var playPauseButton widget.Clickable
var stopButton widget.Clickable

var explorerInstance *explorer.Explorer
var fileReadErr error
var fileDialogActive bool

var playing bool
var paused bool
var playStopCh chan struct{}

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

	for {
		e := w.Event()
		explorerInstance.ListenEvents(e)
		switch evt := e.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, evt)
			gtx.Execute(op.InvalidateCmd{})

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
									case <-time.After(25 * time.Millisecond):
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
			if stopButton.Clicked(gtx) {
				stopPlayback()
			}
			if resetButton.Clicked(gtx) {
				stopPlayback()
				game = Game{
					BoardA: initialBoard.DeepCopy(),
					BoardB: initialBoard.DeepCopy(),
					UseA:   true,
					Turn:   1,
				}
			}
			if nextButton.Clicked(gtx) && (!playing || paused) {
				game.Tick()
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
				}()
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					gen := game.Turn
					label := material.H6(th, "Generation: "+itoa(gen))
					return layout.Center.Layout(gtx, label.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					// TODO: pan/zoom... we need to support the infiniteness of our InfiniteGrid :P
					board := currentBoard(&game)

					viewSize := 40
					cellSize := 15
					margin := 8

					var height, width, minRow, minCol, maxRow, maxCol = 0, 0, 0, 0, 0, 0
					var dimensions, ok = dimensionCache[[3]int{viewSize, cellSize, margin}]
					if !ok {
						minRow, minCol := -viewSize/2, -viewSize/2
						maxRow, maxCol := viewSize/2, viewSize/2
						width = (maxCol-minCol)*cellSize + 2*margin
						height = (maxRow-minRow)*cellSize + 2*margin
						dimensionCache[[3]int{viewSize, cellSize, margin}] = [6]int{minRow, minCol, maxRow, maxCol, height, width}
					} else {
						minRow = dimensions[0]
						minCol = dimensions[1]
						maxRow = dimensions[2]
						maxCol = dimensions[3]
						height = dimensions[4]
						width = dimensions[5]
					}

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
								paint.FillShape(gtx.Ops, col, clip.Rect{Min: image.Pt(0, 0), Max: image.Pt(cellSize, cellSize)}.Op())
								op.Pop()
							}
						}
						gridCol := color.NRGBA{R: 180, G: 180, B: 180, A: 255}
						for i := 0; i <= (maxRow - minRow); i++ {
							y := margin + i*cellSize
							op := op.Offset(image.Pt(margin, y)).Push(gtx.Ops)
							paint.FillShape(gtx.Ops, gridCol, clip.Rect{Min: image.Pt(0, 0), Max: image.Pt((maxCol-minCol)*cellSize, 1)}.Op())
							op.Pop()
						}
						for j := 0; j <= (maxCol - minCol); j++ {
							x := margin + j*cellSize
							op := op.Offset(image.Pt(x, margin)).Push(gtx.Ops)
							paint.FillShape(gtx.Ops, gridCol, clip.Rect{Min: image.Pt(0, 0), Max: image.Pt(1, (maxRow-minRow)*cellSize)}.Op())
							op.Pop()
						}
						return D{Size: image.Pt(width, height)}
					})
				}),
				layout.Rigid(
					func(gtx C) D {
						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(5),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx,
							func(gtx C) D {
								btn := material.Button(th, &nextButton, "Next")
								if playing && !paused {
									btn.Background = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
								}
								return btn.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						margins := layout.Inset{
							Top:    unit.Dp(5),
							Bottom: unit.Dp(5),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx,
							func(gtx C) D {
								btn := material.Button(th, &playPauseButton, func() string {
									if !playing {
										return "Play"
									} else if paused {
										return "Resume"
									} else {
										return "Pause"
									}
								}())
								return btn.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						margins := layout.Inset{
							Top:    unit.Dp(5),
							Bottom: unit.Dp(5),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx,
							func(gtx C) D {
								btn := material.Button(th, &stopButton, "Stop")
								return btn.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						margins := layout.Inset{
							Top:    unit.Dp(5),
							Bottom: unit.Dp(5),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx,
							func(gtx C) D {
								btn := material.Button(th, &resetButton, "Reset")
								return btn.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						margins := layout.Inset{
							Top:    unit.Dp(5),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx,
							func(gtx C) D {
								btn := material.Button(th, &importButton, "Import RLE")
								return btn.Layout(gtx)
							},
						)
					},
				),
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
