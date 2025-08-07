package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRLE_Block(t *testing.T) {
	block := [][]bool{
		{true, true},
		{true, true},
	}
	g := makeInfiniteGrid(block)
	var buf bytes.Buffer
	if err := ExportRLE(&buf, g); err != nil {
		t.Fatalf("ExportRLE failed: %v", err)
	}
	g2, err := ImportRLE(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ImportRLE failed: %v", err)
	}
	minRow, minCol, maxRow, maxCol := g.Bounds()
	if !gridsEqualRegion(g, g2, minRow, minCol, maxRow, maxCol) {
		t.Errorf("Block: exported/imported grid does not match original bounding box [%d:%d,%d:%d]", minRow, maxRow, minCol, maxCol)
		t.Logf("Original grid:")
		for i := minRow; i <= maxRow; i++ {
			var row string
			for j := minCol; j <= maxCol; j++ {
				if g.At(i, j) {
					row += "O"
				} else {
					row += "."
				}
			}
			t.Log(row)
		}
		t.Logf("Round-tripped grid:")
		for i := minRow; i <= maxRow; i++ {
			var row string
			for j := minCol; j <= maxCol; j++ {
				if g2.At(i, j) {
					row += "O"
				} else {
					row += "."
				}
			}
			t.Log(row)
		}
	}
}

func TestRLE_Blinker(t *testing.T) {
	blinker := [][]bool{
		{true, true, true},
	}
	g := makeInfiniteGrid(blinker)
	var buf bytes.Buffer
	if err := ExportRLE(&buf, g); err != nil {
		t.Fatalf("ExportRLE failed: %v", err)
	}
	g2, err := ImportRLE(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ImportRLE failed: %v", err)
	}
	minRow, minCol, maxRow, maxCol := g.Bounds()
	if !gridsEqualRegion(g, g2, minRow, minCol, maxRow, maxCol) {
		t.Errorf("Blinker: exported/imported grid does not match original bounding box [%d:%d,%d:%d]", minRow, maxRow, minCol, maxCol)
		t.Logf("Original grid:")
		for i := minRow; i <= maxRow; i++ {
			var row string
			for j := minCol; j <= maxCol; j++ {
				if g.At(i, j) {
					row += "O"
				} else {
					row += "."
				}
			}
			t.Log(row)
		}
		t.Logf("Round-tripped grid:")
		for i := minRow; i <= maxRow; i++ {
			var row string
			for j := minCol; j <= maxCol; j++ {
				if g2.At(i, j) {
					row += "O"
				} else {
					row += "."
				}
			}
			t.Log(row)
		}
	}
}

func TestRLE_Toad(t *testing.T) {
	toad := [][]bool{
		{false, true, true, true},
		{true, true, true, false},
	}
	g := makeInfiniteGrid(toad)
	var buf bytes.Buffer
	if err := ExportRLE(&buf, g); err != nil {
		t.Fatalf("ExportRLE failed: %v", err)
	}
	g2, err := ImportRLE(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ImportRLE failed: %v", err)
	}
	minRow, minCol, maxRow, maxCol := g.Bounds()
	if !gridsEqualRegion(g, g2, minRow, minCol, maxRow, maxCol) {
		t.Errorf("Toad: exported/imported grid does not match original bounding box [%d:%d,%d:%d]", minRow, maxRow, minCol, maxCol)
		t.Logf("Original grid:")
		for i := minRow; i <= maxRow; i++ {
			var row string
			for j := minCol; j <= maxCol; j++ {
				if g.At(i, j) {
					row += "O"
				} else {
					row += "."
				}
			}
			t.Log(row)
		}
		t.Logf("Round-tripped grid:")
		for i := minRow; i <= maxRow; i++ {
			var row string
			for j := minCol; j <= maxCol; j++ {
				if g2.At(i, j) {
					row += "O"
				} else {
					row += "."
				}
			}
			t.Log(row)
		}
	}
}

func TestRLE_Beacon(t *testing.T) {
	beacon := [][]bool{
		{true, true, false, false},
		{true, true, false, false},
		{false, false, true, true},
		{false, false, true, true},
	}
	rows := len(beacon)
	cols := len(beacon[0])
	g := makeInfiniteGrid(beacon)
	var buf bytes.Buffer
	if err := ExportRLE(&buf, g); err != nil {
		t.Fatalf("ExportRLE failed: %v", err)
	}
	g2, err := ImportRLE(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ImportRLE failed: %v", err)
	}
	if !gridsEqualRegion(g, g2, 0, 0, rows-1, cols-1) {
		t.Errorf("Beacon: exported/imported grid does not match original")
	}
}

func TestRLE_Glider(t *testing.T) {
	glider := [][]bool{
		{false, true, false},
		{false, false, true},
		{true, true, true},
	}
	rows := len(glider)
	cols := len(glider[0])
	g := makeInfiniteGrid(glider)
	var buf bytes.Buffer
	if err := ExportRLE(&buf, g); err != nil {
		t.Fatalf("ExportRLE failed: %v", err)
	}
	g2, err := ImportRLE(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ImportRLE failed: %v", err)
	}
	if !gridsEqualRegion(g, g2, 0, 0, rows-1, cols-1) {
		t.Errorf("Glider: exported/imported grid does not match original")
	}
}

func TestRLE_Diehard(t *testing.T) {
	diehard := [][]bool{
		{false, false, false, false, false, true, false},
		{true, true, false, false, false, false, false},
		{false, true, false, false, true, true, true},
	}
	g := makeInfiniteGrid(diehard)
	var buf bytes.Buffer
	if err := ExportRLE(&buf, g); err != nil {
		t.Fatalf("ExportRLE failed: %v", err)
	}
	g2, err := ImportRLE(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ImportRLE failed: %v", err)
	}
	minRow, minCol, maxRow, maxCol := g.Bounds()
	if !gridsEqualRegion(g, g2, minRow, minCol, maxRow, maxCol) {
		t.Errorf("Diehard: exported/imported grid does not match original bounding box [%d:%d,%d:%d]", minRow, maxRow, minCol, maxCol)
		t.Logf("Original grid:")
		for i := minRow; i <= maxRow; i++ {
			var row string
			for j := minCol; j <= maxCol; j++ {
				if g.At(i, j) {
					row += "O"
				} else {
					row += "."
				}
			}
			t.Log(row)
		}
		t.Logf("Round-tripped grid:")
		for i := minRow; i <= maxRow; i++ {
			var row string
			for j := minCol; j <= maxCol; j++ {
				if g2.At(i, j) {
					row += "O"
				} else {
					row += "."
				}
			}
			t.Log(row)
		}
	}
}
