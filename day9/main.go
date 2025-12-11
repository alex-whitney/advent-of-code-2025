package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/alex-whitney/advent-of-code-2025/lib"
	"github.com/thomaso-mirodin/intmath/intgr"
)

type Today struct {
	points []lib.Point[int]
}

func (d *Today) Init(input string) error {
	lines, err := lib.ReadStringFile(input)
	if err != nil {
		return err
	}

	d.points = make([]lib.Point[int], len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ",")

		x, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		y, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		d.points[i] = lib.NewPoint([]int{x, y})
	}

	return nil
}

func (d *Today) Part1() (string, error) {
	maxArea := 0

	for i := range d.points {
		for j := i + 1; j < len(d.points); j++ {
			p1 := d.points[i]
			p2 := d.points[j]

			area := (intgr.Abs(p1.Coordinates[0]-p2.Coordinates[0]) + 1) *
				(intgr.Abs(p1.Coordinates[1]-p2.Coordinates[1]) + 1)

			if maxArea < area {
				maxArea = area
			}
		}
	}

	return fmt.Sprintf("%d", maxArea), nil
}

func isRectInPolygon(p1 lib.Point[int], p2 lib.Point[int], polygon []lib.Point[int]) (result bool) {
	x1, y1 := p1.Coordinates[0], p1.Coordinates[1]
	x2, y2 := p2.Coordinates[0], p2.Coordinates[1]

	// bounding box crosses an edge?
	// all edges of the polygon are vertical or horizontal, so a bounding box check is sufficient
	for idx := range polygon {
		x3, y3 := polygon[idx].Coordinates[0], polygon[idx].Coordinates[1]

		var x4, y4 int
		if idx < len(polygon)-1 {
			x4, y4 = polygon[idx+1].Coordinates[0], polygon[idx+1].Coordinates[1]
		} else {
			x4, y4 = polygon[0].Coordinates[0], polygon[0].Coordinates[1]
		}

		if intgr.Max(x1, x2) <= intgr.Min(x3, x4) || intgr.Min(x1, x2) >= intgr.Max(x3, x4) ||
			intgr.Max(y1, y2) <= intgr.Min(y3, y4) || intgr.Min(y1, y2) >= intgr.Max(y3, y4) {
			continue
		}

		return false
	}

	return true
}

type solution struct {
	p1   lib.Point[int]
	p2   lib.Point[int]
	area int
}

func (d *Today) Part2() (string, error) {
	// start the same as part 1
	// grab all possible rectangles, sort by area, then filter for those whose perimeters are
	// contained in the polygon

	solutions := []solution{}
	for i := range d.points {
		for j := i + 1; j < len(d.points); j++ {
			p1 := d.points[i]
			p2 := d.points[j]

			area := (intgr.Abs(p1.Coordinates[0]-p2.Coordinates[0]) + 1) *
				(intgr.Abs(p1.Coordinates[1]-p2.Coordinates[1]) + 1)

			solutions = append(solutions, solution{
				p1,
				p2,
				area,
			})
		}
	}

	sort.Slice(solutions, func(i, j int) bool {
		return solutions[i].area > solutions[j].area
	})

	for _, s := range solutions {
		if isRectInPolygon(s.p1, s.p2, d.points) {
			return strconv.Itoa(s.area), nil
		}
	}

	return "", errors.New("no solution found")
}

func main() {
	day := &Today{}
	lib.Run(day)
}
