package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type Today struct {
	devices map[string]Device
}

type Device struct {
	Name    string
	Outputs []string
}

func (d *Today) Init(input string) error {
	lines, err := lib.ReadStringFile(input)
	if err != nil {
		return err
	}

	d.devices = make(map[string]Device)
	for _, line := range lines {
		parts := strings.Split(line, " ")

		device := Device{}
		device.Name = parts[0][0 : len(parts[0])-1]
		device.Outputs = parts[1:]

		d.devices[device.Name] = device
	}

	device := Device{
		Name: "out",
	}
	d.devices[device.Name] = device

	return nil
}

func (d *Today) dfs(start Device, dest string, path string) int {
	if strings.Contains(path, start.Name) {
		return 0
	}

	if start.Name == dest {
		return 1
	}

	counter := 0
	for _, next := range start.Outputs {
		counter += d.dfs(d.devices[next], dest, path+" "+start.Name)
	}

	return counter
}

func (d *Today) Part1() (string, error) {
	counter := d.dfs(d.devices["you"], "out", "")

	return strconv.Itoa(counter), nil
}

func (d *Today) dfsPart2(partialCounts map[string]int, start Device, dest string, ignore []string, path string) int {
	// keep track of path counts for nodes walked

	// explicitly avoid nodes that should not be in this segment
	if slices.Contains(ignore, start.Name) {
		return 0
	}
	if count, ok := partialCounts[start.Name]; ok {
		return count
	}

	if start.Name == dest {
		return 1
	}

	counter := 0
	for _, next := range start.Outputs {
		counter += d.dfsPart2(partialCounts, d.devices[next], dest, ignore, path+" "+start.Name)
	}

	partialCounts[start.Name] = counter
	return counter
}

func (d *Today) Part2() (string, error) {
	// Too slow and complicated to count up all paths in a single traversal
	// can instead count up path segments, and combine the counts at the end

	fftToDac := d.dfsPart2(map[string]int{}, d.devices["fft"], "dac", []string{"srv", "out"}, "")
	dacToFft := d.dfsPart2(map[string]int{}, d.devices["dac"], "fft", []string{"srv", "out"}, "")
	srvToDac := d.dfsPart2(map[string]int{}, d.devices["svr"], "dac", []string{"fft", "out"}, "")
	srvToFft := d.dfsPart2(map[string]int{}, d.devices["svr"], "fft", []string{"dac", "out"}, "")
	dacToOut := d.dfsPart2(map[string]int{}, d.devices["dac"], "out", []string{"fft", "srv"}, "")
	fftToOut := d.dfsPart2(map[string]int{}, d.devices["fft"], "out", []string{"dac", "srv"}, "")

	fmt.Printf("dacToFft %+v\n", dacToFft)
	fmt.Printf("srvToDac %+v\n", srvToDac)
	fmt.Printf("srvToFft %+v\n", srvToFft)
	fmt.Printf("dacToOut %+v\n", dacToOut)
	fmt.Printf("fftToOut %+v\n", fftToOut)

	results := (srvToDac * dacToFft * fftToOut) + (srvToFft * fftToDac * dacToOut)

	return strconv.Itoa(results), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
