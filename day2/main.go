package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Range struct {
	Start int
	End   int
}

type Today struct {
	Ranges []Range
}

func (d *Today) Init(input string) error {
	contents, err := lib.ReadDelimitedFile(input, ",")
	if err != nil {
		return err
	}

	d.Ranges = make([]Range, len(contents[0]))
	for i, rangeStr := range contents[0] {
		parts := strings.Split(rangeStr, "-")

		d.Ranges[i].Start, err = strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		d.Ranges[i].End, err = strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Today) Part1() (string, error) {
	count := 0
	for _, item := range d.Ranges {
		for i := item.Start; i <= item.End; i++ {
			str := strconv.Itoa(i)
			if len(str)%2 == 1 {
				continue
			}

			subLen := len(str) / 2
			if str[0:subLen] == str[subLen:] {
				fmt.Printf("In %d-%d, found %d\n", item.Start, item.End, i)
				count += i
			}
		}
	}

	return strconv.Itoa(count), nil
}

func (d *Today) Part2() (string, error) {
	count := 0
	for _, item := range d.Ranges {
		for i := item.Start; i <= item.End; i++ {
			str := strconv.Itoa(i)

			for substrLen := 1; substrLen <= len(str)/2; substrLen++ {
				if len(str)%substrLen != 0 {
					continue
				}

				repeatLen := len(str) / substrLen
				if strings.Repeat(str[0:substrLen], repeatLen) == str {
					fmt.Printf("In %d-%d, found %d\n", item.Start, item.End, i)
					count += i
					break
				}
			}
		}
	}

	return strconv.Itoa(count), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
