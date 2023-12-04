package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type FoundNumber struct {
	Number int
	Line   int
	Start  int
	End    int
}

type GearConnection struct {
	Number int
	X      int
	Y      int
}

type Gear struct {
	X       int
	Y       int
	Numbers []int
}

func findNumbers(line string, lineIndex int) []FoundNumber {
	var foundNumbers []FoundNumber

	numberString := ""
	startIndex := -1
	endIndex := -1

	for index, char := range line {
		lineCharCount := utf8.RuneCountInString(line)

		if unicode.IsDigit(char) {
			if numberString == "" {
				startIndex = index
			}

			numberString = numberString + string(char)

			if index >= lineCharCount-1 {
				endIndex = index
				number, _ := strconv.Atoi(numberString)
				foundNumber := FoundNumber{Number: number, Line: lineIndex, Start: startIndex, End: endIndex}
				foundNumbers = append(foundNumbers, foundNumber)
				numberString = ""
				startIndex = -1
				endIndex = -1
			}

			continue
		} else if numberString != "" {
			endIndex = index - 1
			number, _ := strconv.Atoi(numberString)
			foundNumber := FoundNumber{Number: number, Line: lineIndex, Start: startIndex, End: endIndex}
			foundNumbers = append(foundNumbers, foundNumber)
			numberString = ""
			startIndex = -1
			endIndex = -1
		}
	}

	return foundNumbers
}

func hasSymbolOnRange(line string, lineIndex int, foundNumber FoundNumber) (bool, []GearConnection) {
	var gears []GearConnection

	for index, char := range line {
		if index < foundNumber.Start-1 || index > foundNumber.End+1 || unicode.IsDigit(char) || char == '.' {
			continue
		}

		if char == '*' {
			gears = append(gears, GearConnection{Number: foundNumber.Number, X: index, Y: lineIndex})
		}

		return true, gears
	}

	return false, gears
}

func extractPartNumbers(foundNumbers []FoundNumber, lines []string) ([]int, []GearConnection) {
	var partNumbers []int
	var gears []GearConnection

	for _, foundNumber := range foundNumbers {
		isPartNumber, foundGears := hasSymbolOnRange(lines[foundNumber.Line], foundNumber.Line, foundNumber)

		if !isPartNumber && foundNumber.Line > 0 {
			isPartNumber, foundGears = hasSymbolOnRange(lines[foundNumber.Line-1], foundNumber.Line-1, foundNumber)
		}

		if !isPartNumber && foundNumber.Line < len(lines)-1 {
			isPartNumber, foundGears = hasSymbolOnRange(lines[foundNumber.Line+1], foundNumber.Line+1, foundNumber)
		}

		if isPartNumber {
			partNumbers = append(partNumbers, foundNumber.Number)
			gears = append(gears, foundGears...)
		}
	}

	return partNumbers, gears
}

func mergeGearConnections(gearConnections []GearConnection) []Gear {
	var gears []Gear

	for _, gearConnection := range gearConnections {
		updatedGear := false

		for i, gear := range gears {
			if gear.X == gearConnection.X && gear.Y == gearConnection.Y {
				gears[i].Numbers = append(gear.Numbers, gearConnection.Number)
				updatedGear = true
			}
		}

		if !updatedGear {
			var gear = Gear{
				X:       gearConnection.X,
				Y:       gearConnection.Y,
				Numbers: []int{gearConnection.Number},
			}

			gears = append(gears, gear)
		}

	}

	return gears
}

func getGearRatioSum(gears []Gear) int {
	sum := 0

	for _, gear := range gears {
		if len(gear.Numbers) == 2 {
			sum = sum + (gear.Numbers[0] * gear.Numbers[1])
		}
	}

	return sum
}

func partNumbersSum(lines []string) (int, int) {
	sum := 0
	var gearConnections []GearConnection
	possibleGears := 0

	for lineIndex, line := range lines {
		foundNumbers := findNumbers(line, lineIndex)
		partNumbers, foundGearConnections := extractPartNumbers(foundNumbers, lines)

		gearConnections = append(gearConnections, foundGearConnections...)

		for _, partNumber := range partNumbers {
			sum = sum + partNumber
		}

		possibleGears = possibleGears + strings.Count(line, "*")
	}

	fmt.Print(possibleGears)

	gears := mergeGearConnections(gearConnections)

	gearRatioSum := getGearRatioSum(gears)

	return sum, gearRatioSum
}

func main() {
	lines, err := readLinesFromFile("../../inputs/day03/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	sum, gearRatioSum := partNumbersSum(lines)

	fmt.Printf("Parts sum: %d\n", sum)
	fmt.Printf("Gear ration sum: %d\n", gearRatioSum)
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
