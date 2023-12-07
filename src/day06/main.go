package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Race struct {
	Time   int
	Record int
}

func main() {
	lines, err := readLinesFromFile("../../inputs/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	fmt.Println("Part 1:", part1(lines))
	fmt.Println("Part 2:", part2(lines))
}

func part1(lines []string) int {
	product := 1

	races := parseLines(lines)

	for _, race := range races {
		possibilities := calculateRacePossibilities(race)
		product = product * possibilities
	}

	return product
}

func part2(lines []string) int {
	product := 1

	races := parseLinesPart2(lines)

	for _, race := range races {
		possibilities := calculateRacePossibilities(race)
		product = product * possibilities
	}

	return product
}

func calculateRacePossibilities(race Race) int {
	possibilities, holdTime, distance := 0, 0, 0

	for distance >= 0 {
		distance = (holdTime * race.Time) - (holdTime * holdTime)

		if distance > race.Record {
			possibilities = possibilities + 1
		}

		holdTime = holdTime + 1
	}

	return possibilities
}

func parseLines(lines []string) []Race {
	var races []Race

	times := fieldsStringToNumber(strings.Fields(strings.TrimPrefix(lines[0], "Time:")))
	records := fieldsStringToNumber(strings.Fields(strings.TrimPrefix(lines[1], "Distance:")))

	for i, time := range times {
		race := Race{Time: time, Record: records[i]}
		races = append(races, race)
	}

	return races
}

func parseLinesPart2(lines []string) []Race {
	var races []Race

	time, _ := strconv.Atoi(strings.Join(strings.Fields(strings.TrimPrefix(lines[0], "Time:")), ""))
	record, _ := strconv.Atoi(strings.Join(strings.Fields(strings.TrimPrefix(lines[1], "Distance:")), ""))

	race := Race{Time: time, Record: record}
	races = append(races, race)

	return races
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

func fieldsStringToNumber(fields []string) []int {
	var intFields []int

	for _, field := range fields {
		intField, _ := strconv.Atoi(field)
		intFields = append(intFields, intField)
	}

	return intFields
}
