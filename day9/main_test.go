package main

import (
	"testing"

	"github.com/alex-whitney/advent-of-code-2025/lib"
	"github.com/stretchr/testify/assert"
)

func TestWindingNumberTest(t *testing.T) {
	polygon := []lib.Point[int]{
		lib.NewPoint([]int{0, 0}),
		lib.NewPoint([]int{10, 0}),
		lib.NewPoint([]int{10, 10}),
		lib.NewPoint([]int{0, 10}),
	}

	tests := []struct {
		x      int
		y      int
		result bool
	}{
		// vertices
		{
			x:      0,
			y:      0,
			result: true,
		}, {
			x:      10,
			y:      0,
			result: true,
		}, {
			x:      10,
			y:      10,
			result: true,
		}, {
			x:      0,
			y:      10,
			result: true,
		},

		// inside points
		{

			x:      5,
			y:      5,
			result: true,
		}, {
			x:      1,
			y:      1,
			result: true,
		}, {
			x:      9,
			y:      1,
			result: true,
		}, {
			x:      1,
			y:      9,
			result: true,
		},

		// perimeter
		{
			x:      5,
			y:      10,
			result: true,
		}, {
			x:      10,
			y:      5,
			result: true,
		}, {
			x:      0,
			y:      5,
			result: true,
		}, {
			x:      5,
			y:      0,
			result: true,
		},

		// outside
		{
			x:      -1,
			y:      -1,
			result: false,
		}, {
			x:      -1,
			y:      0,
			result: false,
		}, {
			x:      11,
			y:      0,
			result: false,
		}, {
			x:      0,
			y:      11,
			result: false,
		}, {
			x:      5,
			y:      11,
			result: false,
		}, {
			x:      11,
			y:      5,
			result: false,
		},
	}

	for _, test := range tests {
		testPoint := lib.NewPoint([]int{test.x, test.y})
		result := windingNumberTest(testPoint, polygon)
		assert.Equalf(t, test.result, result, testPoint.String())
	}
}
