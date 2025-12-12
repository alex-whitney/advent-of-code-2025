package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPart1(t *testing.T) {
	d := &Today{}
	err := d.Init("sample.txt")
	require.NoError(t, err)

	result, err := d.Part1()
	require.NoError(t, err)
	assert.Equal(t, "2", result)
}

func TestPart2(t *testing.T) {
	d := &Today{}
	err := d.Init("sample.txt")
	require.NoError(t, err)

	result, err := d.Part2()
	require.NoError(t, err)
	assert.Equal(t, "World", result)
}
