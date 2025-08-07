package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	rleFile := flag.String("rle", "", "Path to RLE file to import as initial pattern")
	flag.Parse()

	var imported *InfiniteGrid
	if *rleFile != "" {
		f, err := os.Open(*rleFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open RLE file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		b, err := ImportRLE(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to import RLE: %v\n", err)
			os.Exit(1)
		}
		imported = &b
	}
	RunGUI(imported)
}
