package main

import (
	"slices"
	"sort"
	"strconv"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	points []lib.Point[int]

	numConnections int
}

func (d *Today) Init(input string) error {
	tokens, err := lib.ReadDelimitedFile(input, ",")
	if err != nil {
		return err
	}

	d.points = make([]lib.Point[int], len(tokens))
	for i, row := range tokens {
		d.points[i] = lib.NewPoint([]int{0, 0, 0})

		for j := range []int{0, 1, 2} {
			d.points[i].Coordinates[j], err = strconv.Atoi(row[j])
			if err != nil {
				return err
			}
		}
	}

	if len(d.points) == 20 {
		d.numConnections = 10
	} else {
		d.numConnections = 1000
	}

	return nil
}

type Edge struct {
	From     int
	To       int
	Distance float64
}

type Set map[int]struct{}

func (d *Today) Part1() (string, error) {
	edges := []Edge{}

	for r := range d.points {
		for c := r + 1; c < len(d.points); c++ {
			distance, err := d.points[r].Distance(&d.points[c])
			if err != nil {
				return "", err
			}
			edges = append(edges, Edge{
				From:     r,
				To:       c,
				Distance: distance,
			})
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Distance < edges[j].Distance
	})

	// node index -> set
	nodeToSet := map[int]mapset.Set[int]{}
	sets := []mapset.Set[int]{}

	for i := 0; i < d.numConnections; i++ {
		edge := edges[i]

		if nodeToSet[edge.From] == nil && nodeToSet[edge.To] == nil {
			set := mapset.NewSet(edge.From, edge.To)
			nodeToSet[edge.From] = set
			nodeToSet[edge.To] = set
			sets = append(sets, set)
		} else if nodeToSet[edge.From] == nodeToSet[edge.To] {
			// nodes already connected & in the same set
		} else if nodeToSet[edge.From] == nil {
			nodeToSet[edge.To].Add(edge.From)
			nodeToSet[edge.From] = nodeToSet[edge.To]
		} else if nodeToSet[edge.To] == nil {
			nodeToSet[edge.From].Add(edge.To)
			nodeToSet[edge.To] = nodeToSet[edge.From]
		} else {
			// connecting two existing graphs
			graph1 := nodeToSet[edge.From]
			graph2 := nodeToSet[edge.To]

			graph1.Append(graph2.ToSlice()...)
			for val := range graph2.Iter() {
				nodeToSet[val] = graph1
			}

			sets = slices.DeleteFunc(sets, func(item mapset.Set[int]) bool {
				return item == graph2
			})
		}
	}

	sort.Slice(sets, func(i, j int) bool {
		return sets[i].Cardinality() > sets[j].Cardinality()
	})

	result := 1
	for _, set := range sets[0:3] {
		result *= set.Cardinality()
	}

	return strconv.Itoa(result), nil
}

func (d *Today) Part2() (string, error) {
	edges := []Edge{}
	var lastEdge Edge

	for r := range d.points {
		for c := r + 1; c < len(d.points); c++ {
			distance, err := d.points[r].Distance(&d.points[c])
			if err != nil {
				return "", err
			}
			edges = append(edges, Edge{
				From:     r,
				To:       c,
				Distance: distance,
			})
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Distance < edges[j].Distance
	})

	// node index -> set
	nodeToSet := map[int]mapset.Set[int]{}
	sets := []mapset.Set[int]{}

	for i := 0; ; i++ {
		edge := edges[i]

		if nodeToSet[edge.From] == nil && nodeToSet[edge.To] == nil {
			set := mapset.NewSet(edge.From, edge.To)
			nodeToSet[edge.From] = set
			nodeToSet[edge.To] = set
			sets = append(sets, set)
		} else if nodeToSet[edge.From] == nodeToSet[edge.To] {
			// nodes already connected & in the same set
		} else if nodeToSet[edge.From] == nil {
			nodeToSet[edge.To].Add(edge.From)
			nodeToSet[edge.From] = nodeToSet[edge.To]
		} else if nodeToSet[edge.To] == nil {
			nodeToSet[edge.From].Add(edge.To)
			nodeToSet[edge.To] = nodeToSet[edge.From]
		} else {
			// connecting two existing graphs
			graph1 := nodeToSet[edge.From]
			graph2 := nodeToSet[edge.To]

			graph1.Append(graph2.ToSlice()...)
			for val := range graph2.Iter() {
				nodeToSet[val] = graph1
			}

			sets = slices.DeleteFunc(sets, func(item mapset.Set[int]) bool {
				return item == graph2
			})
		}

		if len(nodeToSet) == len(d.points) && len(sets) == 1 {
			lastEdge = edge
			break
		}
	}

	p1 := d.points[lastEdge.From]
	p2 := d.points[lastEdge.To]
	return strconv.Itoa(p1.Coordinates[0] * p2.Coordinates[0]), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
