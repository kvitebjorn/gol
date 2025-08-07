package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/kvitebjorn/gol/internal/game"
)

// ImportRLE parses an RLE file and returns an InfiniteGrid with the pattern.
func ImportRLE(r io.Reader) (game.InfiniteGrid, error) {
	scanner := bufio.NewScanner(r)
	var header string
	var rows, cols int
	var dataLines []string
	headerRe := regexp.MustCompile(`x *= *(\d+), *y *= *(\d+)`)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "+") || strings.HasPrefix(line, "#") {
			continue
		}
		if header == "" {
			header = line
			m := headerRe.FindStringSubmatch(header)
			if m == nil {
				return game.InfiniteGrid{}, errors.New("invalid RLE header")
			}
			cols, _ = strconv.Atoi(m[1])
			rows, _ = strconv.Atoi(m[2])
			continue
		}
		dataLines = append(dataLines, line)
	}
	if rows == 0 || cols == 0 {
		return game.InfiniteGrid{}, errors.New("missing RLE header")
	}
	ig := game.NewInfiniteGrid()
	x, y := 0, 0
	rle := strings.Join(dataLines, "")
	num := 0
	// Always place the pattern at (0,0)
parseLoop:
	for i := 0; i < len(rle); i++ {
		c := rle[i]
		switch {
		case c >= '0' && c <= '9':
			num = num*10 + int(c-'0')
		case c == 'b' || c == 'o':
			n := num
			if n == 0 {
				n = 1
			}
			for j := 0; j < n; j++ {
				if c == 'o' {
					ig.Set(y, x, true)
				}
				x++
			}
			num = 0
		case c == '$':
			n := num
			if n == 0 {
				n = 1
			}
			y += n
			x = 0
			num = 0
		case c == '!':
			break parseLoop
		}
	}
	return ig, nil
}

// ExportRLE writes the InfiniteGrid as an RLE pattern to the writer.
// The exported region is the bounding box of all live cells.
func ExportRLE(w io.Writer, g game.InfiniteGrid) error {
	minRow, minCol, maxRow, maxCol := g.Bounds()
	rows := maxRow - minRow + 1
	cols := maxCol - minCol + 1
	// Export the bounding box region, shifted to (0,0) in the RLE output
	_, err := fmt.Fprintf(w, "x = %d, y = %d\n", cols, rows)
	if err != nil {
		return err
	}
	for y := 0; y < rows; y++ {
		runChar := byte(0)
		runLen := 0
		for x := 0; x < cols; x++ {
			alive := g.At(minRow+y, minCol+x)
			c := byte('b')
			if alive {
				c = 'o'
			}
			if runLen == 0 {
				runChar = c
				runLen = 1
			} else if c == runChar {
				runLen++
			} else {
				if runLen == 1 {
					_, err = fmt.Fprintf(w, "%c", runChar)
				} else {
					_, err = fmt.Fprintf(w, "%d%c", runLen, runChar)
				}
				if err != nil {
					return err
				}
				runChar = c
				runLen = 1
			}
		}
		if runLen > 0 {
			if runLen == 1 {
				_, err = fmt.Fprintf(w, "%c", runChar)
			} else {
				_, err = fmt.Fprintf(w, "%d%c", runLen, runChar)
			}
			if err != nil {
				return err
			}
		}
		if y < rows-1 {
			_, err = fmt.Fprintf(w, "$\n")
			if err != nil {
				return err
			}
		}
	}
	_, err = fmt.Fprintf(w, "!\n")
	return err
}
