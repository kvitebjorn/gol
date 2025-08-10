package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kvitebjorn/gol/internal/board"
	"github.com/kvitebjorn/gol/internal/game"
	"github.com/kvitebjorn/gol/internal/gpu"
	"github.com/kvitebjorn/gol/internal/gui"
	"github.com/kvitebjorn/gol/internal/util"
)

func main() {
	game.UseGpu = gpu.HasCUDA()

	rleFile := flag.String("rle", "", "Path to RLE file to import as initial pattern")
	flag.Parse()

	var imported *board.InfiniteGrid
	if *rleFile != "" {
		f, err := os.Open(*rleFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open RLE file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		b, err := util.ImportRLE(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to import RLE: %v\n", err)
			os.Exit(1)
		}
		imported = &b
	}
	gui.RunGUI(imported)
}
