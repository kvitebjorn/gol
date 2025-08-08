package gui

import (
	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
)

func HandleEvents(gtx C, cache *viewCache, w *app.Window) {
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
				panY -= 4
				changed = true
			case key.NameDownArrow:
				panY += 4
				changed = true
			case key.NameLeftArrow:
				panX -= 4
				changed = true
			case key.NameRightArrow:
				panX += 4
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
				if zoomLevel < 0.1 {
					zoomLevel = 0.1
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
		w.Invalidate()
		cache.img = nil
	}
}

func HandleBoardEvents(gtx C, cache *viewCache, w *app.Window) {
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:  boardTag,
			Kinds:   pointer.Press | pointer.Release | pointer.Drag | pointer.Scroll,
			ScrollY: pointer.ScrollRange{Min: -1, Max: 1},
		})
		if !ok {
			break
		}

		if x, ok := ev.(pointer.Event); ok {
			switch x.Kind {
			case pointer.Press:
				if x.Buttons == pointer.ButtonSecondary {
					isDragging = true
					dragStart = x.Position
					break
				}

				clickPos := x.Position

				toggleCell := true
				if playing && !paused {
					toggleCell = false
				}
				stopPlayback()

				if !toggleCell {
					break
				}

				clickX, clickY := clickPos.X, clickPos.Y
				minRow, minCol, _, _, cellSize, _, _, _ :=
					computeDynamicView(gtx, zoomLevel, panX, panY)

				cellCol := int(clickX) / cellSize
				cellRow := int(clickY) / cellSize

				row := minRow + cellRow
				col := minCol + cellCol

				cur := gameState.CurrentBoard().At(row, col)
				gameState.CurrentBoard().Set(row, col, !cur)

				cache.img = nil
				w.Invalidate()

			case pointer.Drag:
				if isDragging && x.Buttons == pointer.ButtonSecondary {
					dx := x.Position.X - dragStart.X
					dy := x.Position.Y - dragStart.Y

					panX -= int(dx) / 2
					panY -= int(dy) / 2

					dragStart = x.Position

					cache.img = nil
					w.Invalidate()
				}

			case pointer.Scroll:
				dir := x.Scroll.Y
				if dir == -1 {

				}
			}
		}
	}
}
