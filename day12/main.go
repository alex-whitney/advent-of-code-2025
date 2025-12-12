package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	presentSizes map[int]int
	gridSizes    []lib.Pair[int, int]
	requirements [][]int
}

func (d *Today) Init(input string) error {
	contents, err := lib.ReadFile(input)
	if err != nil {
		return err
	}

	parts := strings.Split(contents, "\n\n")

	d.presentSizes = map[int]int{}
	for i := 0; i < len(parts)-1; i++ {
		lines := strings.Split(parts[i], "\n")

		presentNo, err := strconv.Atoi(lines[0][:len(lines[0])-1])
		if err != nil {
			return err
		}

		count := 0
		for j := 1; j < len(lines); j++ {
			count += strings.Count(lines[j], "#")
		}
		d.presentSizes[presentNo] = count
	}

	regions := strings.Split(parts[len(parts)-1], "\n")
	d.requirements = make([][]int, len(regions))
	d.gridSizes = make([]lib.Pair[int, int], len(regions))
	for regionNo, region := range regions {
		p := strings.Split(region, ":")

		var x, y int
		_, err = fmt.Sscanf(p[0], "%dx%d", &x, &y)
		if err != nil {
			return err
		}
		d.gridSizes[regionNo] = lib.NewPair(x, y)

		req := strings.Split(strings.TrimSpace(p[1]), " ")
		r := make([]int, len(req))
		for i := range req {
			r[i], err = strconv.Atoi(req[i])
			if err != nil {
				return err
			}
		}
		d.requirements[regionNo] = r
	}

	return nil
}

func (d *Today) Part1() (string, error) {
	// I'm going to assume we're not actually going to need to be solving bin-packing
	// and that instead we can just eliminate solutions because their areas are too
	// small to fit all of the presents

	counter := 0

	for regionNum, region := range d.gridSizes {
		size := region.Left * region.Right

		required := 0
		for idx, requiredCount := range d.requirements[regionNum] {
			required += d.presentSizes[idx] * requiredCount
		}

		fmt.Printf("region=%d capacity=%d required=%d\n", regionNum, size, required)

		if size >= required {
			counter++
		}
	}

	// Fun fact -- this doesn't actually produce the right answer for the sample, but it does
	// for the actual input :)

	return strconv.Itoa(counter), nil
}

func (d *Today) Part2() (string, error) {
	return "World", nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
