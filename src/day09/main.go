package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines, err := readLinesFromFile("../../inputs/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	sum := part1(lines)
	fmt.Println("Part 1:", sum)
	sum2 := part2(lines)
	fmt.Println("Part 2:", sum2)
}

func part1(lines []string) int {
	sum := 0

	for _, line := range lines {
		values := fieldsStringToNumber(strings.Fields(line))
		nextValues := predictNextValue(values)
		sum = sum + nextValues[len(nextValues)-1]
	}

	return sum
}

func part2(lines []string) int {
	sum := 0

	for _, line := range lines {
		values := fieldsStringToNumber(strings.Fields(line))
		reversedValues := reverseSlice(values)
		nextValues := predictNextValue(reversedValues)
		sum = sum + nextValues[len(nextValues)-1]
	}

	return sum
}

func reverseSlice(values []int) []int {
	length := len(values)

	reversedValues := make([]int, length)

	for i := 0; i < length; i = i + 1 {
		reversedValues[length-i-1] = values[i]
	}

	return reversedValues
}

func predictNextValue(values []int) []int {
	isAllZeros := true

	for i := 1; i < len(values); i = i + 1 {
		if values[i] != 0 {
			isAllZeros = false
			break
		}
	}

	if isAllZeros {
		values = append(values, 0)
		return values
	}

	var differences []int

	for i := 1; i < len(values); i = i + 1 {
		difference := values[i] - values[i-1]
		differences = append(differences, difference)
	}

	nextValues := predictNextValue(differences)

	nexValue := values[len(values)-1] + nextValues[len(nextValues)-1]
	values = append(values, nexValue)

	return values
}

func fieldsStringToNumber(fields []string) []int {
	var intFields []int

	for _, field := range fields {
		intField, _ := strconv.Atoi(field)
		intFields = append(intFields, intField)
	}

	return intFields
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
