package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

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

func mergeDigitMaps(a map[int]rune, b map[int]rune) map[int]rune {
	merged := make(map[int]rune)

	for key, value := range a {
		merged[key] = value
	}

	for key, value := range b {
		merged[key] = value
	}

	return merged
}

func findPlainDigits(line string) map[int]rune {
	digits := make(map[int]rune)

	for i, char := range line {
		if unicode.IsDigit(char) {
			digits[i] = char
		}
	}

	return digits
}

func findSpelledDigits(line string) map[int]rune {
	digits := make(map[int]rune)

	digits[strings.Index(line, "one")] = '1'
	digits[strings.Index(line, "two")] = '2'
	digits[strings.Index(line, "three")] = '3'
	digits[strings.Index(line, "four")] = '4'
	digits[strings.Index(line, "five")] = '5'
	digits[strings.Index(line, "six")] = '6'
	digits[strings.Index(line, "seven")] = '7'
	digits[strings.Index(line, "eight")] = '8'
	digits[strings.Index(line, "nine")] = '9'

	digits[strings.LastIndex(line, "one")] = '1'
	digits[strings.LastIndex(line, "two")] = '2'
	digits[strings.LastIndex(line, "three")] = '3'
	digits[strings.LastIndex(line, "four")] = '4'
	digits[strings.LastIndex(line, "five")] = '5'
	digits[strings.LastIndex(line, "six")] = '6'
	digits[strings.LastIndex(line, "seven")] = '7'
	digits[strings.LastIndex(line, "eight")] = '8'
	digits[strings.LastIndex(line, "nine")] = '9'

	delete(digits, -1)

	return digits
}

func findAllDigits(line string) map[int]rune {
	plainDigits := findPlainDigits(line)
	spelledDigits := findSpelledDigits(line)

	merged := mergeDigitMaps(plainDigits, spelledDigits)

	return merged
}

func getFirstDigit(digits map[int]rune) rune {
	smallestIndex := math.MaxInt
	firstDigit := utf8.RuneError

	for index := range digits {
		if index < smallestIndex {
			smallestIndex = index
			firstDigit = digits[index]
		}
	}

	return firstDigit
}

func getLastDigit(digits map[int]rune) rune {
	biggestIndex := -1
	lastDigit := utf8.RuneError

	for index := range digits {
		if index > biggestIndex {
			biggestIndex = index
			lastDigit = digits[index]
		}
	}

	return lastDigit
}

func findNumberCode(line string) (int, error) {
	digits := findAllDigits(line)
	firstDigit := getFirstDigit(digits)
	lastDigit := getLastDigit(digits)

	var fullNumber string

	if firstDigit != utf8.RuneError {
		fullNumber = fullNumber + string(firstDigit)
	}

	if lastDigit != utf8.RuneError {
		fullNumber = fullNumber + string(lastDigit)
	}

	return strconv.Atoi(fullNumber)
}

func getSumOfCodes(lines []string) int {
	sum := 0

	for _, line := range lines {
		code, err := findNumberCode(line)

		if err != nil {
			continue
		}

		sum = sum + code
	}

	return sum
}

func main() {
	lines, err := readLinesFromFile("../../inputs/day01/input")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	sum := getSumOfCodes(lines)

	fmt.Println(sum)
}
