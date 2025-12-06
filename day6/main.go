package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	values     [][]int
	operations []string

	// for part 2
	rawInput []string
}

func (d *Today) Init(input string) error {
	lines, err := lib.ReadStringFile(input)
	if err != nil {
		return err
	}

	d.rawInput = lines

	d.values = make([][]int, len(lines)-1)
	for i, line := range lines {
		re := regexp.MustCompile(`\s+`)
		tokens := re.Split(strings.TrimSpace(line), -1)

		if i < len(lines)-1 {
			d.values[i] = make([]int, len(tokens))
			for j, val := range tokens {
				d.values[i][j], err = strconv.Atoi(val)
				if err != nil {
					return err
				}
			}
		} else {
			d.operations = tokens
		}
	}

	return nil
}

func sum(values [][]int, operators []string) (int, error) {
	result := 0

	for i, op := range operators {
		rowResult := 0

		for j, val := range values[i] {
			if j == 0 {
				rowResult = val
			} else {
				if op == "+" {
					rowResult += val
				} else if op == "*" {
					rowResult *= val
				} else {
					return 0, errors.New("unknown operator: " + op)
				}
			}
		}

		result += rowResult
	}

	return result, nil
}

func (d *Today) Part1() (string, error) {
	values := lib.Transpose(d.values)
	result, err := sum(values, d.operations)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}

func (d *Today) Part2() (string, error) {
	// parse the file again. column starts are where the operator is located
	stringValues := d.rawInput[:len(d.rawInput)-1]
	operatorString := d.rawInput[len(d.rawInput)-1]

	columnStarts := []int{}
	for i, val := range operatorString {
		// already have parsed the operator list, so can ignore those here
		if val != ' ' {
			columnStarts = append(columnStarts, i)
		}
	}

	values := make([][]int, len(d.operations))
	for operatorIndex, columnStart := range columnStarts {
		column := columnStart
		for {
			stringVal := ""
			for row := range stringValues {
				stringVal = stringVal + string(stringValues[row][column])
			}

			stringVal = strings.TrimSpace(stringVal)
			if stringVal == "" {
				// end of the section
				break
			}

			val, err := strconv.Atoi(stringVal)
			if err != nil {
				return "", err
			}
			values[operatorIndex] = append(values[operatorIndex], val)

			column++
			if column >= len(operatorString) {
				// no more columns
				break
			}
		}
	}

	result, err := sum(values, d.operations)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
