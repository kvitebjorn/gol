package gui

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/kvitebjorn/gol/internal/util"
)

type C = layout.Context
type D = layout.Dimensions

func runWindow(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()

	var cache viewCache

	explorer := GetExplorerInstance(w)

	// Main render loop
	for {
		e := w.Event()
		if explorer != nil {
			explorer.ListenEvents(e)
		}

		switch evt := e.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, evt)

			HandleEvents(gtx, &cache, w)
			HandleControlClicks(gtx, &cache, w)

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
