package gui

import (
	"gioui.org/app"
	"gioui.org/io/key"
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
