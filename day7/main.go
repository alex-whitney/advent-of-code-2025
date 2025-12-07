package main

import (
	"strconv"
	"strings"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	grid [][]rune
}

func (d *Today) Init(input string) error {
	lines, err := lib.ReadStringFile(input)
	if err != nil {
		return err
	}

	d.grid = make([][]rune, len(lines))
	for i, line := range lines {
		line = strings.Replace(line, "S", "|", 1)
		d.grid[i] = []rune(line)
	}

	return nil
}

func (d *Today) Part1() (string, error) {
	splitCount := 0

	resultGrid := make([][]rune, len(d.grid))
	resultGrid[0] = append([]rune{}, d.grid[0]...)

	for row := 1; row < len(d.grid); row++ {
		resultGrid[row] = append([]rune{}, d.grid[row]...)

		for col := range d.grid[row] {
			if d.grid[row][col] == '.' && resultGrid[row-1][col] == '|' {
				resultGrid[row][col] = '|'
			}

			if d.grid[row][col] == '^' && resultGrid[row-1][col] == '|' {
				splitCount++
				if col > 0 && resultGrid[row][col-1] == '.' {
					resultGrid[row][col-1] = '|'
				}
				if col < len(d.grid)-1 && resultGrid[row][col+1] == '.' {
					resultGrid[row][col+1] = '|'
				}
			}
		}
	}

	return strconv.Itoa(splitCount), nil
}

func (d *Today) Part2() (string, error) {
	// only thing that matters is the number of paths that end up at each point
	// that means we can propogate numbers down the array, adding values when the paths overlap
	// then just sum the paths at the end

	// initialize - all paths travel through starting point
	lastRow := make([]int, len(d.grid[0]))
	lastRow[strings.Index(string(d.grid[0]), "|")] = 1

	for row := 1; row < len(d.grid); row++ {
		thisRow := make([]int, len(d.grid[row]))

		for col := range d.grid[row] {
			if d.grid[row][col] == '.' {
				thisRow[col] += lastRow[col]
			}

			if d.grid[row][col] == '^' {
				if col > 0 {
					thisRow[col-1] += lastRow[col]
				}
				if col < len(d.grid)-1 {
					thisRow[col+1] += lastRow[col]
				}
			}
		}

		lastRow = thisRow
	}

	result := lib.Sum(lastRow)
	return strconv.Itoa(result), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
