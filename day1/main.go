package main

import (
	"strconv"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	Instructions []int
}

func (d *Today) Init(input string) error {
	fileContents, err := lib.ReadStringFile(input)
	if err != nil {
		return err
	}

	d.Instructions = make([]int, len(fileContents))
	for i, row := range fileContents {
		d.Instructions[i], err = strconv.Atoi(row[1:])
		if err != nil {
			return err
		}

		if row[0] == 'L' {
			d.Instructions[i] = -1 * d.Instructions[i]
		}
	}

	return nil
}

func (d *Today) Part1() (string, error) {
	counter := 50
	result := 0

	for _, instruction := range d.Instructions {
		counter += instruction

		counter = counter % 100
		if counter < 0 {
			counter = counter + 100
		} else if counter == 0 {
			result++
		}
	}

	return strconv.Itoa(result), nil
}

func (d *Today) Part2() (string, error) {
	counter := 50
	result := 0

	for _, instruction := range d.Instructions {
		rotations := instruction / 100
		if rotations < 0 {
			result = result - rotations
		} else {
			result = result + rotations
		}
		instruction = instruction % 100

		initialPosition := counter
		counter += instruction

		if counter < 0 {
			counter += 100

			if initialPosition > 0 {
				result++
			}
		} else if counter > 99 {
			counter = counter - 100
			result++
		} else if counter == 0 {
			result++
		}
	}

	return strconv.Itoa(result), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
