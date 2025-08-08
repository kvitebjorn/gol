package gui

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
	"github.com/kvitebjorn/gol/internal/util"
)

type C = layout.Context
type D = layout.Dimensions

func runWindow(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()

	var cache viewCache

	if explorerInstance == nil {
		explorerInstance = explorer.NewExplorer(w)
	}

	for {
		e := w.Event()
		if explorerInstance != nil {
			explorerInstance.ListenEvents(e)
		}

		switch evt := e.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, evt)

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

			HandleControlClicks(gtx, w)

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					gen := gameState.Turn
					label := material.H6(th, "Generation: "+util.Itoa(gen))
					return layout.Center.Layout(gtx, label.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					label := material.Body1(th, fmt.Sprintf("Zoom: %.2fx  Pan: (%d,%d)", zoomLevel, panX, panY))
					return layout.Center.Layout(gtx, label.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return LayoutBoard(gtx, &cache, zoomLevel, panX, panY, w)
				}),
				layout.Rigid(func(gtx C) D {
					return LayoutControls(gtx, th, w)
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
