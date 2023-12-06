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

type Mapping struct {
	Code   int
	Start  int
	End    int
	Length int
}

type Almanac struct {
	Seeds        []int
	Soils        map[int]Mapping
	Fertilizers  map[int]Mapping
	Waters       map[int]Mapping
	Lights       map[int]Mapping
	Temperatures map[int]Mapping
	Humidities   map[int]Mapping
	Locations    map[int]Mapping
}

func getLowestLocationCode(lines []string) int {

	almanac := parseLines(lines)

	var locations []int

	for _, seed := range almanac.Seeds {
		soil := getCode(almanac.Soils, seed)
		fertilizer := getCode(almanac.Fertilizers, soil)
		water := getCode(almanac.Waters, fertilizer)
		light := getCode(almanac.Lights, water)
		temperature := getCode(almanac.Temperatures, light)
		humidity := getCode(almanac.Humidities, temperature)
		location := getCode(almanac.Locations, humidity)

		locations = append(locations, location)
	}

	lowestLocationCode := math.MaxInt

	for _, location := range locations {
		if location < lowestLocationCode {
			lowestLocationCode = location
		}
	}

	return lowestLocationCode
}

func getCode(mappings map[int]Mapping, value int) int {
	code := value

	for _, mapping := range mappings {
		if value >= mapping.Start && value <= mapping.End {
			difference := value - mapping.Start
			return mapping.Code + difference
		}
	}

	return code
}

func lineToMapping(line string) Mapping {
	fields := fieldsStringToNumber(strings.Fields(line))
	return Mapping{
		Code:   fields[0],
		Start:  fields[1],
		End:    fields[1] + fields[2] - 1,
		Length: fields[2],
	}
}

func parseLines(lines []string) Almanac {
	almanac := Almanac{
		Soils:        make(map[int]Mapping),
		Fertilizers:  make(map[int]Mapping),
		Waters:       make(map[int]Mapping),
		Lights:       make(map[int]Mapping),
		Temperatures: make(map[int]Mapping),
		Humidities:   make(map[int]Mapping),
		Locations:    make(map[int]Mapping),
	}

	section := "seeds"

	for lineIndex, line := range lines {

		if line == "" {
			continue
		}

		firstRune, _ := utf8.DecodeRuneInString(line)

		if lineIndex != 0 && unicode.IsLetter(firstRune) {
			section = strings.Fields(line)[0]
			continue
		}

		switch section {
		case "seeds":
			seeds := fieldsStringToNumber(strings.Fields(strings.TrimPrefix(line, "seeds:")))
			almanac.Seeds = seeds
		case "seed-to-soil":
			mapping := lineToMapping(line)
			almanac.Soils[mapping.Code] = mapping
		case "soil-to-fertilizer":
			mapping := lineToMapping(line)
			almanac.Fertilizers[mapping.Code] = mapping
		case "fertilizer-to-water":
			mapping := lineToMapping(line)
			almanac.Waters[mapping.Code] = mapping
		case "water-to-light":
			mapping := lineToMapping(line)
			almanac.Lights[mapping.Code] = mapping
		case "light-to-temperature":
			mapping := lineToMapping(line)
			almanac.Temperatures[mapping.Code] = mapping
		case "temperature-to-humidity":
			mapping := lineToMapping(line)
			almanac.Humidities[mapping.Code] = mapping
		case "humidity-to-location":
			mapping := lineToMapping(line)
			almanac.Locations[mapping.Code] = mapping
		}
	}

	return almanac
}

func main() {
	lines, err := readLinesFromFile("../../inputs/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	location := getLowestLocationCode(lines)

	fmt.Println("Location:", location)
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
