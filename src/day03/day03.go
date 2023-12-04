package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type FoundNumber struct {
	Number int
	Line   int
	Start  int
	End    int
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

func hasSymbolOnRange(line string, start int, end int) bool {
	for index, char := range line {
		if index < start-1 || index > end+1 {
			continue
		}

		if char != '.' && !unicode.IsDigit(char) {
			return true
		}
	}

	return false
}

func extractPartNumbers(foundNumbers []FoundNumber, lines []string) []int {
	var partNumbers []int

	for _, foundNumber := range foundNumbers {
		isPartNumber := hasSymbolOnRange(lines[foundNumber.Line], foundNumber.Start, foundNumber.End)

		if !isPartNumber && foundNumber.Line > 0 {
			isPartNumber = hasSymbolOnRange(lines[foundNumber.Line-1], foundNumber.Start, foundNumber.End)
		}

		if !isPartNumber && foundNumber.Line < len(lines)-1 {
			isPartNumber = hasSymbolOnRange(lines[foundNumber.Line+1], foundNumber.Start, foundNumber.End)
		}

		if isPartNumber {
			partNumbers = append(partNumbers, foundNumber.Number)
		}
	}

	return partNumbers
}

func partNumbersSum(lines []string) int {
	sum := 0

	for lineIndex, line := range lines {
		foundNumbers := findNumbers(line, lineIndex)
		partNumbers := extractPartNumbers(foundNumbers, lines)
		fmt.Println(partNumbers)

		for _, partNumber := range partNumbers {
			sum = sum + partNumber
		}

	}

	return sum
}

func main() {
	lines, err := readLinesFromFile("../../inputs/day03/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	sum := partNumbersSum(lines)

	fmt.Println(sum)
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
