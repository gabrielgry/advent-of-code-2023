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
	lines, err := readLinesFromFile("../../inputs/day02/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	sum := 0

	for _, line := range lines {
		minCubeQuantity := getGameMinCubeQuantity(line)

		power := minCubeQuantity["red"] * minCubeQuantity["green"] * minCubeQuantity["blue"]

		sum = sum + power
	}

	fmt.Println(sum)
}

func getGameMinCubeQuantity(line string) map[string]int {
	_, setsString, _ := strings.Cut(line, ":")

	minCubeQuantity := map[string]int{"red": 0, "green": 0, "blue": 0}

	for _, setString := range strings.Split(setsString, ";") {
		for _, cubeString := range strings.Split(setString, ",") {
			cubeCountString, cubeColor, _ := strings.Cut(strings.TrimSpace(cubeString), " ")

			cubeCount, err := strconv.Atoi(cubeCountString)

			if err != nil {
				log.Fatal("Could not convert cubeCountString to int")
			}

			if cubeCount > minCubeQuantity[cubeColor] {
				minCubeQuantity[cubeColor] = cubeCount
			}
		}
	}

	return minCubeQuantity
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
