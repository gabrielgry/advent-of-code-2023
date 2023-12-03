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
	cubeQuantity := map[string]int{"red": 12, "green": 13, "blue": 14}

	lines, err := readLinesFromFile("../../inputs/day02/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	sum := 0

	for _, line := range lines {
		game, isValidGame := validadeGame(line, cubeQuantity)

		if isValidGame {
			sum = sum + game
		}
	}

	fmt.Println(sum)
}

func validadeGame(line string, cubeQuantity map[string]int) (int, bool) {
	header, setsString, _ := strings.Cut(line, ":")
	game, err := strconv.Atoi(strings.TrimPrefix(header, "Game "))

	if err != nil {
		log.Fatal("Could not convert header value to int")
	}
	
	for _, setString := range strings.Split(setsString, ";") {
		for _, cubeString := range strings.Split(setString, ",") {
			cubeCountString, cubeColor, _ := strings.Cut(strings.TrimSpace(cubeString), " ")

			cubeCount, err := strconv.Atoi(cubeCountString)

			if err != nil {
				log.Fatal("Could not convert cubeCountString to int")
			}

			if cubeCount > cubeQuantity[cubeColor] {
				return game, false
			}
		}
	}

	return game, true
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
