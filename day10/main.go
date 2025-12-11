package main

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alex-whitney/advent-of-code-2025/lib"
)

type LightState []bool

func (s LightState) String() string {
	ret := "["
	for _, v := range s {
		if v {
			ret = ret + "#"
		} else {
			ret = ret + "."
		}
	}
	return ret + "]"
}

type Button []int
type JoltageRequirements []int

func (s JoltageRequirements) equals(other JoltageRequirements) bool {
	if len(s) != len(other) {
		return false
	}

	for i := range s {
		if s[i] != other[i] {
			return false
		}
	}

	return true
}

func (s JoltageRequirements) isValid(other JoltageRequirements) bool {
	if len(s) != len(other) {
		return false
	}

	for i := range s {
		if s[i] > other[i] {
			return false
		}
	}

	return true
}

type Machine struct {
	indicatorLights     LightState
	wiringSchematics    []Button
	joltageRequirements JoltageRequirements
}

type Today struct {
	machines []Machine
}

func (d *Today) Init(input string) error {
	lines, err := lib.ReadStringFile(input)
	if err != nil {
		return err
	}

	d.machines = make([]Machine, len(lines))
	for machineNumber, line := range lines {
		machine := Machine{}
		parts := strings.Split(line, " ")

		machine.indicatorLights = make([]bool, len(parts[0])-2)
		for i := 1; i < len(parts[0])-1; i++ {
			machine.indicatorLights[i-1] = parts[0][i] == '#'
		}

		p := parts[len(parts)-1]
		machine.joltageRequirements, err = lib.ParseIntegerSlice(p[1:len(p)-1], ",")
		if err != nil {
			return err
		}

		machine.wiringSchematics = make([]Button, len(parts)-2)
		for i := 1; i < len(parts)-1; i++ {
			p = parts[i]
			machine.wiringSchematics[i-1], err = lib.ParseIntegerSlice(p[1:len(p)-1], ",")
			if err != nil {
				return err
			}
		}

		d.machines[machineNumber] = machine
	}

	return nil
}

func toggle(state LightState, button Button) LightState {
	ret := make([]bool, len(state))
	copy(ret, state)

	for _, i := range button {
		ret[i] = !state[i]
	}

	return ret
}

type searchState struct {
	state         LightState
	buttonPresses []Button
}

func findShortestPresses(machine Machine) ([]Button, error) {
	initialState := make(LightState, len(machine.indicatorLights))

	queue := []searchState{}

	visited := map[string]struct{}{
		initialState.String(): {},
	}
	for _, button := range machine.wiringSchematics {
		s := toggle(initialState, button)
		queue = append(queue, searchState{
			state:         s,
			buttonPresses: []Button{button},
		})
	}

	targetString := machine.indicatorLights.String()

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		if _, ok := visited[next.state.String()]; ok {
			continue
		}

		if next.state.String() == targetString {
			return next.buttonPresses, nil
		}

		visited[next.state.String()] = struct{}{}

		for _, button := range machine.wiringSchematics {
			s := toggle(next.state, button)
			queue = append(queue, searchState{
				state:         s,
				buttonPresses: append(next.buttonPresses, button),
			})
		}
	}

	return []Button{}, errors.New("couldn't find a solution")
}

func (d *Today) Part1() (string, error) {
	// uhh..
	// graph, with each node being a distinct configuration of lights,
	// and each (directed) edge being a button press
	//
	// that means each graph has 2^n elements
	// and 2^n * m edges
	//
	// can BFS from all off -> desired state
	//
	// Scanning the input, that seems like it should be fine. I'm assuming the joltages
	// are costs or something and we'll want to have the graph anyway

	totalButtonPresses := 0

	for _, machine := range d.machines {
		buttonPresses, err := findShortestPresses(machine)
		if err != nil {
			return "", err
		}

		totalButtonPresses += len(buttonPresses)
	}

	return strconv.Itoa(totalButtonPresses), nil
}

func determineMaxPresses(machine Machine, currentJoltage JoltageRequirements, button Button) int {
	min := math.MaxInt64
	for _, i := range button {
		val := machine.joltageRequirements[i] - currentJoltage[i]

		if val < min {
			min = val
		}
	}

	if min < 0 {
		fmt.Printf("Invalid state - %v %v\n", currentJoltage, button)
		min = 0
	}

	return min
}

func determinePossiblePresses(machine Machine, currentJoltage JoltageRequirements, button Button, otherButtons []Button) []int {
	max := determineMaxPresses(machine, currentJoltage, button)
	min := 0

	for _, i := range button {
		val := machine.joltageRequirements[i] - currentJoltage[i]
		maxPressedFromOtherbuttons := 0

		for _, b := range otherButtons {
			if slices.Contains(b, i) {
				m := determineMaxPresses(machine, currentJoltage, b)
				maxPressedFromOtherbuttons += m
			}
		}

		rem := val - maxPressedFromOtherbuttons
		if rem < 0 {
			rem = 0
		}

		if rem > min {
			min = rem
		}

		//fmt.Printf("   tgt=%d: curr=%d; rem=%d\n", machine.joltageRequirements[i], currentJoltage[i], rem)
	}

	// I'm pretty sure this full range doesn't need to be iterated over
	// probably just consider every combination of max presses for
	// overlapping buttons
	var possibilities []int
	if max >= min {
		possibilities = make([]int, max-min+1)
		for i := max; i >= min; i-- {
			possibilities[max-i] = i
		}
	}

	if len(possibilities) > 0 {
		//fmt.Printf("tgt=%v: curr=%v b=%v max=%d min=%d otherb=%v => pos=%v\n", machine.joltageRequirements, currentJoltage, button, max, min, otherButtons, possibilities)
	}

	return possibilities
}

func pressButton(currentJoltage JoltageRequirements, button Button, times int) JoltageRequirements {
	ret := make(JoltageRequirements, len(currentJoltage))
	copy(ret, currentJoltage)

	for _, i := range button {
		ret[i] += times
	}

	return ret
}

type Solution struct {
	minResult int
}

func explore(machine Machine, currentJoltage JoltageRequirements, buttons []Button, counter int, solution *Solution) bool {
	if currentJoltage.equals(machine.joltageRequirements) {
		if counter < solution.minResult {
			solution.minResult = counter
		}
		return true
	}
	if len(buttons) == 0 {
		return false
	}
	if !currentJoltage.isValid(machine.joltageRequirements) {
		fmt.Printf("whoops: tgt=%v, curr=%v, b=%v\n", machine.joltageRequirements, currentJoltage, buttons)
		return false
	}
	if solution.minResult < counter {
		return true
	}

	button := buttons[0]
	otherButtons := buttons[1:]

	possibilities := determinePossiblePresses(machine, currentJoltage, button, otherButtons)

	for _, i := range possibilities {
		buttonCount := counter + i
		joltage := pressButton(currentJoltage, button, i)

		hasSolution := explore(machine, joltage, otherButtons, buttonCount, solution)
		if hasSolution {
			// bail early because any further iterations will result in more button presses
			// -- turns out this doesn't actually hold to be true
			// if len(otherButtons) == 0 || len(otherButtons[0]) < len(button) {
			//	  return solutions, true
			// }
		}
	}

	return solution.minResult < math.MaxInt64
}

func (d *Today) solve(worker int, workerCount int, solution chan<- int) {
	for i := worker; i < len(d.machines); i += workerCount {
		t := time.Now()

		machine := d.machines[i]

		sortedButtons := make([]Button, len(machine.wiringSchematics))
		copy(sortedButtons, machine.wiringSchematics)

		// desc
		sort.Slice(sortedButtons, func(i, j int) bool {
			return len(sortedButtons[i]) > len(sortedButtons[j])
		})

		soln := &Solution{
			minResult: math.MaxInt64,
		}
		hasSolution := explore(machine, make([]int, len(machine.joltageRequirements)), sortedButtons, 0, soln)
		if !hasSolution {
			panic(fmt.Sprintf("didn't find solution for machine %d", i))
		}

		fmt.Printf("completed machine %d: %d in %dms\n", i, soln.minResult, time.Since(t).Milliseconds())
		solution <- soln.minResult
	}
}

func (d *Today) Part2() (string, error) {
	// don't know if exactly the same approach will work - the graph is finite but counters are
	// pretty high, so it would take a very long time to find the solution with a graph traversal
	//
	// pretty sure this is a pretty straightforward linear algebra problem, which I don't remember
	// how to solve. good luck to me
	//
	// this took over an hour to run on my computer, and most of that time was waiting for the
	// solver to get through a single solution. most finish very quickly, but one took forever.

	numWorkers := 24
	var wg sync.WaitGroup
	results := make(chan int, len(d.machines))

	for i := range numWorkers {
		wg.Go(func() {
			d.solve(i, numWorkers, results)
		})
	}

	wg.Wait()
	close(results)

	counter := 0
	for val := range results {
		counter += val
	}

	return strconv.Itoa(counter), nil
}

func main() {
	day := &Today{}
	lib.Run(day)
}
