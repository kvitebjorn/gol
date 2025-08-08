package gui

import (
	"image"
	"image/color"
	"image/draw"
	"sort"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

func computeDynamicView(gtx layout.Context, zoom float64, panX, panY int) (
	minRow, minCol, maxRow, maxCol, cellSize, margin, width, height int) {
	availableWidth := gtx.Constraints.Max.X
	availableHeight := gtx.Constraints.Max.Y

	cellSizeF := zoom * 20
	cellSize = int(cellSizeF)
	if cellSize > 50 {
		cellSize = 50
	}
	if cellSize < 2 {
		cellSize = 2
	}

	cols := availableWidth / cellSize
	rows := availableHeight / cellSize

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

func LayoutBoard(
	gtx layout.Context,
	cache *viewCache,
	zoomLevel float64,
	panX, panY int,
	w *app.Window) layout.Dimensions {
	minRow, minCol, maxRow, maxCol, cellSize, margin, width, height :=
		computeDynamicView(gtx, zoomLevel, panX, panY)

	if width <= 0 || height <= 0 || cellSize <= 0 {
		return layout.Dimensions{Size: image.Pt(0, 0)}
	}

	size := gtx.Constraints.Max
	gtx.Constraints.Min = size

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = width
		gtx.Constraints.Max.Y = height

		return boardClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			event.Op(gtx.Ops, boardTag)

			HandleBoardEvents(gtx, cache, w)

			useCache := cache.img != nil &&
				cache.turn == gameState.Turn &&
				cache.panX == panX &&
				cache.panY == panY &&
				cache.zoom == zoomLevel &&
				cache.width == width &&
				cache.height == height &&
				cache.cellSize == cellSize

			if !useCache {
				img := image.NewRGBA(image.Rect(0, 0, width, height))

				bg := image.NewUniform(color.NRGBA{R: 220, G: 220, B: 220, A: 255})
				draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

				// Render an entire row at once
				// This is more efficient than rendering cell by cell!
				colsByRow := map[int][]int{}
				for _, p := range gameState.CurrentBoard().AliveCells() {
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

						x := margin + (start-minCol)*cellSize
						wPixels := (last - start + 1) * cellSize
						rect := image.Rect(x, y, x+wPixels, y+cellSize)
						draw.Draw(img, rect, fillCol, image.Point{}, draw.Src)
						start = cols[i]
						last = cols[i]
					}

					x := margin + (start-minCol)*cellSize
					wPixels := (last - start + 1) * cellSize
					rect := image.Rect(x, y, x+wPixels, y+cellSize)
					draw.Draw(img, rect, fillCol, image.Point{}, draw.Src)
				}

				gridCol := image.NewUniform(color.NRGBA{R: 180, G: 180, B: 180, A: 255})

				for i := 0; i <= (maxRow - minRow); i++ {
					y := margin + i*cellSize
					rect := image.Rect(margin, y, margin+(maxCol-minCol)*cellSize, y+1)
					draw.Draw(img, rect, gridCol, image.Point{}, draw.Src)
				}

				for j := 0; j <= (maxCol - minCol); j++ {
					x := margin + j*cellSize
					rect := image.Rect(x, margin, x+1, margin+(maxRow-minRow)*cellSize)
					draw.Draw(img, rect, gridCol, image.Point{}, draw.Src)
				}

				cache.img = img
				cache.turn = gameState.Turn
				cache.panX = panX
				cache.panY = panY
				cache.zoom = zoomLevel
				cache.width = width
				cache.height = height
				cache.cellSize = cellSize
			}

			if cache.img != nil {
				paint.NewImageOp(cache.img).Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
			}

			return layout.Dimensions{Size: image.Pt(width, height)}
		})
	})
}
