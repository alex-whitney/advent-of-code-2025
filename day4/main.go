package main

import (
	"strconv"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	Paper [][]bool

	RowCount int
	ColCount int
}

func (d *Today) Init(input string) error {
	lines, err := lib.ReadStringFile(input)
	if err != nil {
		return err
	}

	d.RowCount = len(lines)
	d.ColCount = len(lines[0])
	d.Paper = make([][]bool, d.RowCount)
	for row, line := range lines {
		d.Paper[row] = make([]bool, d.ColCount)
		for col, val := range line {
			d.Paper[row][col] = val == '@'
		}
	}

	return nil
}

func (d *Today) Part1() (string, error) {
	accessibleCount := 0
	for row := 0; row < d.RowCount; row++ {
		for col := 0; col < d.ColCount; col++ {
			if !d.Paper[row][col] {
				continue
			}

			count := 0
			for offsetRow := -1; offsetRow <= 1; offsetRow++ {
				for offsetCol := -1; offsetCol <= 1; offsetCol++ {
					if offsetCol == 0 && offsetRow == 0 {
						continue
					}
					if row+offsetRow < 0 || row+offsetRow >= d.RowCount {
						continue
					}
					if col+offsetCol < 0 || col+offsetCol >= d.ColCount {
						continue
					}

					if d.Paper[row+offsetRow][col+offsetCol] {
						count++
					}
				}
			}

			if count < 4 {
				accessibleCount++
			}
		}
	}

	return strconv.Itoa(accessibleCount), nil
}

func (d *Today) countAndRemove(paper [][]bool) (int, [][]bool) {
	result := make([][]bool, d.ColCount)

	accessibleCount := 0
	for row := 0; row < d.RowCount; row++ {
		result[row] = make([]bool, d.ColCount)

		for col := 0; col < d.ColCount; col++ {
			if !paper[row][col] {
				continue
			}

			count := 0
			for offsetRow := -1; offsetRow <= 1; offsetRow++ {
				for offsetCol := -1; offsetCol <= 1; offsetCol++ {
					if offsetCol == 0 && offsetRow == 0 {
						continue
					}
					if row+offsetRow < 0 || row+offsetRow >= d.RowCount {
						continue
					}
					if col+offsetCol < 0 || col+offsetCol >= d.ColCount {
						continue
					}

					if paper[row+offsetRow][col+offsetCol] {
						count++
					}
				}
			}

			if count < 4 {
				accessibleCount++
			} else {
				result[row][col] = true
			}
		}
	}

	return accessibleCount, result

}

func (d *Today) Part2() (string, error) {
	total := 0
	currentPaper := d.Paper

	for {
		removed := 0
		removed, currentPaper = d.countAndRemove(currentPaper)

		if removed == 0 {
			break
		}

		total += removed
	}

	return strconv.Itoa(total), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
