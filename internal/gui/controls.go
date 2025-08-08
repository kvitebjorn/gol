package gui

import (
	"time"

	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/kvitebjorn/gol/internal/game"
	"github.com/kvitebjorn/gol/internal/util"
)

var (
	nextButton      widget.Clickable
	playPauseButton widget.Clickable
	resetButton     widget.Clickable
	importButton    widget.Clickable
)

func LayoutControls(gtx layout.Context, th *material.Theme, w *app.Window) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(10),
		Bottom: unit.Dp(10),
		Left:   unit.Dp(10),
		Right:  unit.Dp(10),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceSides,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &nextButton, "Next")
				if playing && !paused {
					btn.Background = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
				}
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &resetButton, "Reset")
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &importButton, "Import")
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
			}),
		)
	})
}

func HandleControlClicks(gtx C, cache *viewCache, w *app.Window) {
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
							gameState.Tick()
							win.Invalidate()
							select {
							case <-stopCh:
								return
							case <-time.After(10 * time.Millisecond):
							}
						} else {
							select {
							case <-stopCh:
								return
							case <-time.After(100 * time.Millisecond):
							}
						}
					}
				}
			}(playStopCh, w)
			w.Invalidate()
		} else {
			paused = !paused
			w.Invalidate()
		}
	}
	if resetButton.Clicked(gtx) {
		stopPlayback()
		gameState = game.Game{
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
		gameState.Tick()
		w.Invalidate()
	}
	if importButton.Clicked(gtx) && !fileDialogActive {
		fileDialogActive = true
		go func(win *app.Window) {
			explorer := GetExplorerInstance(win)
			r, err := explorer.ChooseFile(".rle")
			if err != nil {
				fileReadErr = err
				fileDialogActive = false
				return
			}
			defer r.Close()
			b, err := util.ImportRLE(r)
			if err != nil {
				fileReadErr = err
			} else {
				initialBoard = b.DeepCopy()
				fileReadErr = nil
				stopPlayback()
				gameState = game.Game{
					BoardA: initialBoard.DeepCopy(),
					BoardB: initialBoard.DeepCopy(),
					UseA:   true,
					Turn:   1,
				}
				zoomLevel = 1.0
				panX = 0
				panY = 0
				cache.img = nil
				win.Invalidate()
			}
			fileDialogActive = false
		}(w)
	}
}
