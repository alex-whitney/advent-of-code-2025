package main

import (
	"math"
	"strconv"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	Banks [][]int
}

func (d *Today) Init(input string) error {
	lines, err := lib.ReadStringFile(input)
	if err != nil {
		return err
	}

	d.Banks = make([][]int, len(lines))
	for i, line := range lines {
		d.Banks[i], err = lib.ParseIntegerSlice(line, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Today) Part1() (string, error) {
	sum := 0
	for _, bank := range d.Banks {
		max := 0
		for i := 0; i < len(bank); i++ {
			for j := i + 1; j < len(bank); j++ {
				val := bank[i]*10 + bank[j]
				if val > max {
					max = val
				}
			}
		}
		sum += max
	}

	return strconv.Itoa(sum), nil
}

func (d *Today) Part2() (string, error) {
	sum := 0
	for _, bank := range d.Banks {
		value := 0

		// 100 choose 12 is pretty big
		// Can greedily pick the biggest number, don't need to look at all options

		lastDigitIndex := -1
		for digit := 0; digit < 12; digit++ {
			maxDigit := 0
			currentMaxDigitIndex := -1

			// the 12th digit can be the last item in the slice
			lastPossibleDigit := len(bank) - (11 - digit)
			for i := lastDigitIndex + 1; i < lastPossibleDigit; i++ {
				if bank[i] > maxDigit {
					currentMaxDigitIndex = i
					maxDigit = bank[i]
				}
			}

			lastDigitIndex = currentMaxDigitIndex
			value += maxDigit * int(math.Pow10(12-digit-1))
		}

		sum += value
	}

	return strconv.Itoa(sum), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
