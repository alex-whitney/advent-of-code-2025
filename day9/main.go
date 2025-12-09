package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

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

func isLeft(p0 lib.Point[int], p1 lib.Point[int], p2 lib.Point[int]) int {
	return ((p1.Coordinates[0]-p0.Coordinates[0])*(p2.Coordinates[1]-p0.Coordinates[1]) -
		(p2.Coordinates[0]-p0.Coordinates[0])*(p1.Coordinates[1]-p0.Coordinates[1]))
}

func windingNumberTest(point lib.Point[int], polygon []lib.Point[int]) bool {
	// https://web.archive.org/web/20130126163405/http://geomalgorithms.com/a03-_inclusion.html

	counter := 0

	for i := range polygon {
		// iterate over each edge
		p1 := polygon[i]

		var p2 lib.Point[int]
		if i == len(polygon)-1 {
			p2 = polygon[0]
		} else {
			p2 = polygon[i+1]
		}

		var crossProduct int
		computedCrossProduct := false
		getCrossProduct := func(p0 lib.Point[int], p1 lib.Point[int], p2 lib.Point[int]) int {
			if !computedCrossProduct {
				crossProduct = isLeft(p0, p1, p2)
				computedCrossProduct = true

			}

			return crossProduct
		}

		if p1.Coordinates[1] <= point.Coordinates[1] {
			if p2.Coordinates[1] > point.Coordinates[1] {
				if getCrossProduct(p1, p2, point) > 0 {
					counter++
				}
			}
		} else {
			if p2.Coordinates[1] <= point.Coordinates[1] {
				if getCrossProduct(p1, p2, point) < 0 {
					counter--
				}
			}
		}

		// special case -- is this point _on_ the edge?
		// winding number treats the edges as outside
		if point.Coordinates[0] >= intgr.Min(p1.Coordinates[0], p2.Coordinates[0]) &&
			point.Coordinates[0] <= intgr.Max(p1.Coordinates[0], p2.Coordinates[0]) &&
			point.Coordinates[1] >= intgr.Min(p1.Coordinates[1], p2.Coordinates[1]) &&
			point.Coordinates[1] <= intgr.Max(p1.Coordinates[1], p2.Coordinates[1]) {
			if getCrossProduct(point, p1, p2) == 0 {
				return true
			}
		}
	}

	return counter != 0
}

func isRectInPolygon(p1 lib.Point[int], p2 lib.Point[int], polygon []lib.Point[int]) (result bool) {
	x1 := p1.Coordinates[0]
	x2 := p2.Coordinates[0]
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	y1 := p1.Coordinates[1]
	y2 := p2.Coordinates[1]
	if y2 < y1 {
		y1, y2 = y2, y1
	}

	for x := x1; x <= x2; x++ {
		if !windingNumberTest(lib.NewPoint([]int{x, y1}), polygon) {
			return false
		}
		if !windingNumberTest(lib.NewPoint([]int{x, y2}), polygon) {
			return false
		}
	}

	for y := y1; y <= y2; y++ {
		if !windingNumberTest(lib.NewPoint([]int{x1, y}), polygon) {
			return false
		}
		if !windingNumberTest(lib.NewPoint([]int{x2, y}), polygon) {
			return false
		}
	}

	return true
}

type solution struct {
	p1   lib.Point[int]
	p2   lib.Point[int]
	area int
}

func (d *Today) checkSolution(solutions []solution, worker int, workerCount int, validSolutions chan<- solution) {
	for i := worker; i < len(solutions); i += workerCount {
		if isRectInPolygon(solutions[i].p1, solutions[i].p2, d.points) {
			validSolutions <- solutions[i]
		} else {
			fmt.Printf("Checked solution %d of %d (%.2f%%)\n", i+1, len(solutions), float64(i+1)/float64(len(solutions))*100)
		}
	}
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

	// I could probably eagerly exit if I find a solution, rather than checking all possible
	// solutions, but that's a lot more bookkeeping
	numWorkers := 24
	var wg sync.WaitGroup
	results := make(chan solution, len(solutions))

	for i := range numWorkers {
		wg.Go(func() {
			d.checkSolution(solutions, i, numWorkers, results)
		})
	}

	wg.Wait()
	close(results)

	valid := []solution{}
	for s := range results {
		valid = append(valid, s)
	}

	sort.Slice(valid, func(i, j int) bool {
		return valid[i].area > valid[j].area
	})

	return strconv.Itoa(valid[0].area), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
