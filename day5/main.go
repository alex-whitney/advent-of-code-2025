package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	ingredients map[int]bool
	fresh       []lib.Pair[int, int]
}

func (d *Today) Init(input string) error {
	text, err := lib.ReadFile(input)
	if err != nil {
		return err
	}

	parts := strings.Split(text, "\n\n")

	ingredientRanges := strings.Split(parts[0], "\n")
	d.fresh = make([]lib.Pair[int, int], len(ingredientRanges))
	for i, ingredientRange := range ingredientRanges {
		parts := strings.Split(ingredientRange, "-")

		first, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		second, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		d.fresh[i] = lib.NewPair(first, second)
	}

	d.ingredients = make(map[int]bool)
	ingredients := strings.Split(parts[1], "\n")
	for _, ingredient := range ingredients {
		val, err := strconv.Atoi(ingredient)
		if err != nil {
			return err
		}

		d.ingredients[val] = true
	}

	return nil
}

func (d *Today) Part1() (string, error) {
	freshCount := 0

	for ingredient := range d.ingredients {
		usedCount := 0
		for _, freshRange := range d.fresh {
			if ingredient >= freshRange.Left && ingredient <= freshRange.Right {
				usedCount++
			}
		}

		if usedCount > 0 {
			freshCount++
		}
	}

	return strconv.Itoa(freshCount), nil
}

func mergeRanges(inRanges []lib.Pair[int, int]) []lib.Pair[int, int] {
	outRanges := make([]lib.Pair[int, int], 0)

	fmt.Printf("Merging ranges\n")

	for _, thisRange := range inRanges {
		// Case 1: This range doesn't overlap any existing ranges
		//   in this case, add it as a new range to the list
		isNewRange := true
		overlapsRange := make([]bool, len(outRanges))
		for i, otherRange := range outRanges {
			if thisRange.Right >= otherRange.Left && thisRange.Left <= otherRange.Right {
				isNewRange = false
				overlapsRange[i] = true
			}
		}
		if isNewRange {
			outRanges = append(outRanges, thisRange)
			continue
		}

		// Case 2: This range is fully overlapped by an existing range
		//   in this case, do nothing
		// Case 3: This range extends an existing range.
		//   in this case, extend the range in the list
		// These can be handled the same
		for i, otherRange := range outRanges {
			if !overlapsRange[i] {
				continue
			}

			if otherRange.Left > thisRange.Left {
				outRanges[i].Left = thisRange.Left
			}
			if otherRange.Right < thisRange.Right {
				outRanges[i].Right = thisRange.Right
			}
		}

		// Case 3b: This range extends an existing range into another existing range
		//   In this case, calling mergeRanges again should merge those ranges
		//   Would need to continue merging until the result is stable
	}

	return outRanges
}

func (d *Today) Part2() (string, error) {
	// keeping track of items in a map would be impractical
	// max value is 562 328 260 897 038

	ranges := d.fresh
	for {
		beforeSize := len(ranges)
		ranges = mergeRanges(ranges)

		if len(ranges) == beforeSize {
			break
		}
	}

	fmt.Printf("\ndone merging \n")

	freshCount := 0
	for _, ingredientRange := range ranges {
		freshCount += ingredientRange.Right - ingredientRange.Left + 1
	}

	return strconv.Itoa(freshCount), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
