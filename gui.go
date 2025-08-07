package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"sort"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
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

var boardClickable widget.Clickable
var boardTag = new(bool)

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

// Used to provide a fast path for rendering alive cells in the view.
type aliveProvider interface {
	AliveCells() [][2]int
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
	availableWidth := gtx.Constraints.Max.X - 2*gtx.Dp(unit.Dp(10))
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

func runWindow(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()

	if explorerInstance == nil {
		explorerInstance = explorer.NewExplorer(w)
	}

	var tag = new(bool)
	event.Op(&ops, tag)

	// simple persistent cache variables captured in closure
	var cache viewCache

	// Event loop
	for {
		e := w.Event()
		explorerInstance.ListenEvents(e)
		switch evt := e.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, evt)

			// Key events
			changed := false
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
						changed = true
					case key.NameDownArrow:
						panY += 2
						changed = true
					case key.NameLeftArrow:
						panX -= 2
						changed = true
					case key.NameRightArrow:
						panX += 2
						changed = true
					case "+":
						old := zoomLevel
						zoomLevel *= 1.1
						if zoomLevel > 4 {
							zoomLevel = 4
						}
						if zoomLevel != old {
							changed = true
						}
					case "-":
						old := zoomLevel
						zoomLevel *= 0.9
						if zoomLevel < 0.25 {
							zoomLevel = 0.25
						}
						if zoomLevel != old {
							changed = true
						}
					}
				} else {
					break
				}
			}
			if changed {
				// Invalidate both Gio and our cache (cache will be rebuilt lazily on next paint)
				w.Invalidate()
				cache.img = nil
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
					// new play started — ensure one frame
					w.Invalidate()
				} else {
					paused = !paused
					w.Invalidate()
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
				cache.img = nil
			}
			if nextButton.Clicked(gtx) && (!playing || paused) {
				game.Tick()
				w.Invalidate()
				cache.img = nil
			}
			if importButton.Clicked(gtx) && !fileDialogActive {
				fileDialogActive = true
				go func(win *app.Window) {
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
						win.Invalidate()
						cache.img = nil
					}
					fileDialogActive = false
				}(w)
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

					// If the computed width/height is zero or negative, skip heavy rendering.
					if width <= 0 || height <= 0 || cellSize <= 0 {
						return D{Size: image.Pt(0, 0)}
					}

					size := gtx.Constraints.Max
					gtx.Constraints.Min = size
					return layout.Center.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Max.X = width
						gtx.Constraints.Max.Y = height

						return boardClickable.Layout(gtx, func(gtx C) D {
							event.Op(gtx.Ops, boardTag)

							for {
								ev, ok := gtx.Event(pointer.Filter{
									Target: boardTag,
									Kinds:  pointer.Press,
								})
								if !ok {
									break
								}

								if x, ok := ev.(pointer.Event); ok {
									switch x.Kind {
									case pointer.Press:

										clickPos := x.Position

										stopPlayback()

										clickX, clickY := clickPos.X, clickPos.Y
										minRow, minCol, _, _, cellSize, _, _, _ := computeDynamicView(gtx, zoomLevel, panX, panY)

										cellCol := int(clickX) / cellSize
										cellRow := int(clickY) / cellSize

										row := minRow + cellRow
										col := minCol + cellCol

										board := currentBoard(&game)
										cur := board.At(row, col)
										board.Set(row, col, !cur)

										cache.img = nil
										w.Invalidate()
									}
								}
							}

							// If cache matches, paint cached image directly
							useCache := cache.img != nil &&
								cache.turn == game.Turn &&
								cache.panX == panX &&
								cache.panY == panY &&
								cache.zoom == zoomLevel &&
								cache.width == width &&
								cache.height == height &&
								cache.cellSize == cellSize

							if !useCache {
								// build new image cache for this view
								img := image.NewRGBA(image.Rect(0, 0, width, height))

								// fill background (light gray)
								bg := image.NewUniform(color.NRGBA{R: 220, G: 220, B: 220, A: 255})
								draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

								// Bucket cells by row so we can draw contiguous runs more efficiently.
								colsByRow := map[int][]int{}
								for _, p := range board.AliveCells() {
									r := p[0]
									c := p[1]
									if r < minRow || r >= maxRow || c < minCol || c >= maxCol {
										continue
									}
									colsByRow[r] = append(colsByRow[r], c)
								}

								fillCol := image.NewUniform(color.NRGBA{R: 0, G: 200, B: 0, A: 255})
								for r, cols := range colsByRow {
									if len(cols) == 0 {
										continue
									}
									sort.Ints(cols)
									start := cols[0]
									last := start
									y := margin + (r-minRow)*cellSize
									for i := 1; i < len(cols); i++ {
										if cols[i] == last || cols[i] == last+1 {
											last = cols[i]
											continue
										}

										// draw run start..last
										x := margin + (start-minCol)*cellSize
										wPixels := (last - start + 1) * cellSize
										rect := image.Rect(x, y, x+wPixels, y+cellSize)
										draw.Draw(img, rect, fillCol, image.Point{}, draw.Src)
										start = cols[i]
										last = cols[i]
									}

									// final run
									x := margin + (start-minCol)*cellSize
									wPixels := (last - start + 1) * cellSize
									rect := image.Rect(x, y, x+wPixels, y+cellSize)
									draw.Draw(img, rect, fillCol, image.Point{}, draw.Src)
								}

								// draw grid lines into image
								gridCol := image.NewUniform(color.NRGBA{R: 180, G: 180, B: 180, A: 255})

								// horizontal lines
								for i := 0; i <= (maxRow - minRow); i++ {
									y := margin + i*cellSize
									rect := image.Rect(margin, y, margin+(maxCol-minCol)*cellSize, y+1)
									draw.Draw(img, rect, gridCol, image.Point{}, draw.Src)
								}

								// vertical lines
								for j := 0; j <= (maxCol - minCol); j++ {
									x := margin + j*cellSize
									rect := image.Rect(x, margin, x+1, margin+(maxRow-minRow)*cellSize)
									draw.Draw(img, rect, gridCol, image.Point{}, draw.Src)
								}

								// store cache
								cache.img = img
								cache.turn = game.Turn
								cache.panX = panX
								cache.panY = panY
								cache.zoom = zoomLevel
								cache.width = width
								cache.height = height
								cache.cellSize = cellSize
							}

							// paint cached image (single image op)
							if cache.img != nil {
								paint.NewImageOp(cache.img).Add(gtx.Ops)
								paint.PaintOp{}.Add(gtx.Ops)
							}

							return D{Size: image.Pt(width, height)}
						})
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
		if err := runWindow(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
