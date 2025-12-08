package lib

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type Point[T Number] struct {
	Coordinates []T
}

func NewPoint[T Number](coordinates []T) Point[T] {
	return Point[T]{
		Coordinates: coordinates,
	}
}

func (p *Point[T]) String() string {
	s := make([]string, len(p.Coordinates))

	for i, c := range p.Coordinates {
		s[i] = fmt.Sprintf("%v", c)
	}

	return "(" + strings.Join(s, ",") + ")"
}

func (p *Point[T]) Distance(p2 *Point[T]) (float64, error) {
	if len(p.Coordinates) != len(p2.Coordinates) {
		return 0, errors.New("point dimensions are not equal")
	}

	var total float64
	for i := range p.Coordinates {
		total += math.Pow(float64(p.Coordinates[i]-p2.Coordinates[i]), 2)
	}

	return math.Sqrt(total), nil
}
