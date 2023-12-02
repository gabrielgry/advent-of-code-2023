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

func findFirstDigit(code string) rune {
	for _, char := range code {
		if unicode.IsDigit(char) {
			return char
		}
	}

	return utf8.RuneError
}

func findLastDigit(code string) rune {
	for i, w := len(code), 0; i >= 0; i = i - w {
		char, size := utf8.DecodeLastRuneInString(code[:i])

		if unicode.IsDigit(char) {
			return char
		}

		w = size
	}

	return utf8.RuneError
}

func findNumberCode(word string) (int, error) {
	firstDigit := findFirstDigit(word)
	lastDigit := findLastDigit(word)

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
	lines, err := readLinesFromFile("../../input/day01")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	sum := getSumOfCodes(lines)

	fmt.Println(sum)
}
