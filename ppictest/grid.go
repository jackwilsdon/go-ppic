package ppictest

import "fmt"

func Parse(source [8]string) (grid [8][8]bool) {
	for y, row := range source {
		if l := len(row); l != 8 {
			panic(fmt.Sprintf("len(grid[%d]) != 8 (got %d)", y, l))
		}

		for x, c := range row {
			if c != '#' && c != ' ' {
				panic(fmt.Sprintf("grid[%d][%d] is not '#' or ' ' (got %q)", y, x, c))
			}

			grid[y][x] = c == '#'
		}
	}

	return
}
