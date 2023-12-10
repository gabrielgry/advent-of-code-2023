package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"unicode/utf8"
)

type Directions struct {
	Left  string
	Right string
}

type Locations map[string]Directions

func main() {
	lines, err := readLinesFromFile("../../inputs/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	steps := part1(lines)
	fmt.Println("Part 1:", steps)
	steps2 := part2(lines)
	fmt.Println("Part 2:", steps2)
}

func part1(lines []string) int {
	instructions, locations := parseLines(lines)

	steps := 0

	found := false
	instructionIndex := 0
	currentLocation := "AAA"

	for !found {
		instruction := instructions[instructionIndex]

		if currentLocation == "ZZZ" {
			found = true
			break
		}

		if instruction == "L" {
			currentLocation = locations[currentLocation].Left
		} else {
			currentLocation = locations[currentLocation].Right
		}

		if instructionIndex >= len(instructions)-1 {
			instructionIndex = 0
		} else {
			instructionIndex = instructionIndex + 1
		}

		steps = steps + 1
	}

	return steps
}

func part2(lines []string) *big.Int {
	instructions, locations := parseLines(lines)

	var currentLocations []string
	for location := range locations {
		if lastRune, _ := utf8.DecodeLastRuneInString(location); lastRune == 'A' {
			currentLocations = append(currentLocations, location)
		}
	}

	loopSteps := make([]*big.Int, len(currentLocations))

	for index, currentLocation := range currentLocations {
		steps := findLoopSteps(currentLocation, instructions, locations)
		loopSteps[index] = new(big.Int).SetInt64(int64(steps))
	}

	minSteps := findLCMOfArray(loopSteps)

	return minSteps
}

func findLCMOfArray(numbers []*big.Int) *big.Int {
	if len(numbers) == 0 {
		return nil
	}

	lcm := new(big.Int).Set(numbers[0])

	for _, num := range numbers[1:] {
		lcm = findLCM(lcm, num)
	}

	return lcm
}

func findLCM(a, b *big.Int) *big.Int {
	// Use GCD to find LCM: LCM(a, b) = |a * b| / GCD(a, b)
	var gcd big.Int
	gcd.GCD(nil, nil, a, b)

	lcm := new(big.Int).Abs(a)
	lcm.Mul(lcm, b)
	lcm.Div(lcm, &gcd)

	return lcm
}

func findLoopSteps(startLocation string, instructions []string, locations Locations) int {

	instructionIndex := 0
	currentLocation := startLocation

	previousLoopSteps := -1
	loopSteps := 0

	previousSteps := 0
	steps := 0

	for previousLoopSteps != loopSteps {
		instruction := instructions[instructionIndex]

		if lastRune, _ := utf8.DecodeLastRuneInString(currentLocation); lastRune == 'Z' {
			previousLoopSteps = loopSteps
			loopSteps = steps - previousSteps
			previousSteps = steps
		}

		if instruction == "L" {
			currentLocation = locations[currentLocation].Left
		} else {
			currentLocation = locations[currentLocation].Right
		}

		if instructionIndex >= len(instructions)-1 {
			instructionIndex = 0
		} else {
			instructionIndex = instructionIndex + 1
		}

		steps = steps + 1
	}

	return loopSteps
}

func parseLines(lines []string) ([]string, Locations) {
	instructions := strings.Split(lines[0], "")

	locations := make(Locations)

	for index, line := range lines {
		if index < 2 {
			continue
		}

		location, directions, _ := strings.Cut(line, "=")
		fields := strings.Fields(directions)

		location = strings.TrimSpace(location)
		left := strings.Trim(fields[0], "(,)")
		right := strings.Trim(fields[1], "(,)")

		locations[location] = Directions{Left: left, Right: right}
	}

	return instructions, locations
}

func readLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
