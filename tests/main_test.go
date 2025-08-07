package main

import (
	"testing"

	"github.com/kvitebjorn/gol/internal/game"
)

func TestBlinkerOscillator(t *testing.T) {
	// Blinker (period 2 oscillator)
	start := [][]bool{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, true, true, true, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
	}
	boardA := makeInfiniteGrid(start)
	boardB := game.NewInfiniteGrid()
	game := game.Game{BoardA: boardA, BoardB: boardB, UseA: true, Turn: 1}

	// Save initial state
	initial := boardA.DeepCopy()

	game.Tick()
	mid := *game.CurrentBoard()

	game.Tick()
	end := *game.CurrentBoard()

	if gridsEqualRegion(initial, mid, 0, 0, 4, 4) {
		t.Errorf("Blinker should change after one tick.")
	}
	if !gridsEqualRegion(initial, end, 0, 0, 4, 4) {
		t.Errorf("Blinker should return to original state after two ticks.")
	}
}

func TestToadOscillator(t *testing.T) {
	// Toad (period 2 oscillator)
	start := [][]bool{
		{false, false, false, false, false, false},
		{false, false, false, false, false, false},
		{false, false, true, true, true, false},
		{false, true, true, true, false, false},
		{false, false, false, false, false, false},
		{false, false, false, false, false, false},
	}
	boardA := makeInfiniteGrid(start)
	boardB := game.NewInfiniteGrid()
	game := game.Game{BoardA: boardA, BoardB: boardB, UseA: true, Turn: 1}
	initial := boardA.DeepCopy()
	game.Tick()
	mid := *game.CurrentBoard()
	game.Tick()
	end := *game.CurrentBoard()
	if gridsEqualRegion(initial, mid, 0, 0, 5, 5) {
		t.Errorf("Toad should change after one tick.")
	}
	if !gridsEqualRegion(initial, end, 0, 0, 5, 5) {
		t.Errorf("Toad should return to original state after two ticks.")
	}
}

func TestGliderSpaceship(t *testing.T) {
	start := [][]bool{
		{false, false, false, false, false},
		{false, true, false, false, false},
		{false, false, true, false, false},
		{true, true, true, false, false},
		{false, false, false, false, false},
	}
	boardA := makeInfiniteGrid(start)
	boardB := game.NewInfiniteGrid()
	game := game.Game{BoardA: boardA, BoardB: boardB, UseA: true, Turn: 1}
	seen := make(map[string]bool)
	for i := 0; i < 4; i++ {
		cur := *game.CurrentBoard()
		key := ""
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				if cur.At(r, c) {
					key += "1"
				} else {
					key += "0"
				}
			}
		}
		if seen[key] {
			t.Errorf("Glider repeated a state at tick %d, which should not happen in first 4 generations", i)
		}
		seen[key] = true
		game.Tick()
	}
}

func TestDiehardPattern(t *testing.T) {
	width := 20
	height := 8
	start := make([][]bool, height)
	for i := range start {
		start[i] = make([]bool, width)
	}
	row := 2
	col := 6
	start[row+0][col+6] = true // (2,12)
	start[row+1][col+0] = true // (3,6)
	start[row+1][col+1] = true // (3,7)
	start[row+2][col+1] = true // (4,7)
	start[row+2][col+5] = true // (4,11)
	start[row+2][col+6] = true // (4,12)
	start[row+2][col+7] = true // (4,13)
	boardA := makeInfiniteGrid(start)
	boardB := game.NewInfiniteGrid()
	game := game.Game{BoardA: boardA, BoardB: boardB, UseA: true, Turn: 1}
	maxTicks := 200
	allDead := false
	for i := 0; i < maxTicks; i++ {
		cur := *game.CurrentBoard()
		alive := false
		for r := 0; r < height; r++ {
			for c := 0; c < width; c++ {
				if cur.At(r, c) {
					alive = true
					break
				}
			}
			if alive {
				break
			}
		}
		if !alive {
			allDead = true
			t.Logf("All cells dead at tick %d", i)
			break
		}
		game.Tick()
	}
	if !allDead {
		t.Errorf("Diehard pattern should eventually disappear (all cells dead after %d ticks)", maxTicks)
	}
}
