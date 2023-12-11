package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type Direction string

const (
	North Direction = "north"
	East  Direction = "east"
	South Direction = "south"
	West  Direction = "west"
)

// type PipeType string

// const (
// 	Pipe   = "pipe"
// 	Bend   = "bend"
// 	Empty  = "empty"
// 	Unknow = "unknow"
// )

// type Tile struct {
// 	Symbol      string
// 	Type        PipeType
// 	Connections []Direction
// }

// var Ground = Tile{
// 	Symbol: ".",
// 	Type:   Empty,
// }

// var Start = Tile{
// 	Symbol:      "S",
// 	Type:        Unknow,
// 	Connections: []Direction{North, East, South, West},
// }

// var VPipe = Tile{
// 	Symbol:      "|",
// 	Type:        Pipe,
// 	Connections: []Direction{North, South},
// }

// var HPipe = Tile{
// 	Symbol:      "-",
// 	Type:        Pipe,
// 	Connections: []Direction{East, West},
// }

// var NEBend = Tile{
// 	Symbol:      "L",
// 	Type:        Bend,
// 	Connections: []Direction{North, East},
// }

// var NWBend = Tile{
// 	Symbol:      "J",
// 	Type:        Bend,
// 	Connections: []Direction{North, West},
// }

// var SWBend = Tile{
// 	Symbol:      "7",
// 	Type:        Bend,
// 	Connections: []Direction{South, West},
// }

// var SEBend = Tile{
// 	Symbol:      "F",
// 	Type:        Bend,
// 	Connections: []Direction{South, East},
// }

type Position struct {
	X, Y int
}

func main() {
	lines, err := readLinesFromFile("../../inputs/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	distance := part1(lines)
	fmt.Println("Part 1:", distance)
	area := part2(lines)
	fmt.Println("Part 2:", area)
}

func part1(lines []string) int {
	tiles := parseLines(lines)

	var startPosition Position

	for i := 0; i < len(tiles); i = i + 1 {
		for j := 0; j < len(tiles); j = j + 1 {
			if tiles[i][j] == "S" {
				startPosition = Position{X: j, Y: i}
			}
		}
	}

	var pathLenght int

	if found, sum, _, _ := walkPath(startPosition, North, true, tiles); found {
		pathLenght = sum
	} else if found, sum, _, _ := walkPath(startPosition, East, true, tiles); found {
		pathLenght = sum
	} else if found, sum, _, _ := walkPath(startPosition, South, true, tiles); found {
		pathLenght = sum
	} else if found, sum, _, _ := walkPath(startPosition, West, true, tiles); found {
		pathLenght = sum
	}

	farthestPoint := pathLenght / 2

	return farthestPoint
}

func part2(lines []string) int {
	tiles := parseLines(lines)

	var startPosition Position

	for i := 0; i < len(tiles); i = i + 1 {
		for j := 0; j < len(tiles); j = j + 1 {
			if tiles[i][j] == "S" {
				startPosition = Position{X: j, Y: i}
			}
		}
	}

	var path map[int]map[int]string

	if found, _, foundPath, lastDirection := walkPath(startPosition, North, true, tiles); found {
		path = foundPath
		if lastDirection == North {
			tiles[startPosition.Y][startPosition.X] = "|"
		}
		if lastDirection == East {
			tiles[startPosition.Y][startPosition.X] = "J"
		}
		if lastDirection == West {
			tiles[startPosition.Y][startPosition.X] = "L"
		}
	} else if found, _, foundPath, lastDirection := walkPath(startPosition, East, true, tiles); found {
		path = foundPath
		if lastDirection == East {
			tiles[startPosition.Y][startPosition.X] = "-"
		}
		if lastDirection == North {
			tiles[startPosition.Y][startPosition.X] = "F"
		}
		if lastDirection == South {
			tiles[startPosition.Y][startPosition.X] = "L"
		}
	} else if found, _, foundPath, lastDirection := walkPath(startPosition, South, true, tiles); found {
		path = foundPath
		if lastDirection == South {
			tiles[startPosition.Y][startPosition.X] = "|"
		}
		if lastDirection == East {
			tiles[startPosition.Y][startPosition.X] = "7"
		}
		if lastDirection == West {
			tiles[startPosition.Y][startPosition.X] = "F"
		}
	} else if found, _, foundPath, lastDirection := walkPath(startPosition, West, true, tiles); found {
		path = foundPath
		if lastDirection == West {
			tiles[startPosition.Y][startPosition.X] = "-"
		}
		if lastDirection == North {
			tiles[startPosition.Y][startPosition.X] = "J"
		}
		if lastDirection == South {
			tiles[startPosition.Y][startPosition.X] = "7"
		}
	}

	area := calculateArea(path, tiles)

	return area
}

func calculateArea(path map[int]map[int]string, tiles [][]string) int {
	area := 0

	height := len(tiles)

	for j := 0; j < height; j = j + 1 {
		count := countTilesInsidePath(path[j], tiles[j])
		area = area + count
	}

	return area
}

func countTilesInsidePath(path map[int]string, tiles []string) int {
	if len(path) == 0 {
		return 0
	}

	count := 0

	isInside := false
	previousCorner := ""
	for i := 0; i < len(tiles); i = i + 1 {
		symbol := tiles[i]
		_, isPath := path[i]

		if isPath && symbol == "|" {
			isInside = !isInside
			continue
		}

		if isPath && symbol == "-" {
			continue
		}

		if isPath {
			if previousCorner == "" {
				previousCorner = symbol
				continue
			}

			// FJ or L7
			if (previousCorner == "F" && symbol == "J") || (previousCorner == "L" && symbol == "7") {
				isInside = !isInside
				previousCorner = ""
				continue
			}

			// F7 or LJ
			if (previousCorner == "F" && symbol == "7") || (previousCorner == "L" && symbol == "J") {
				previousCorner = ""
				continue
			}
		}

		if !isPath && isInside {
			count = count + 1
		}
	}

	return count
}

func walkPath(current Position, goesTo Direction, start bool, tiles [][]string) (bool, int, map[int]map[int]string, Direction) {
	symbol := tiles[current.Y][current.X]

	switch symbol {
	case "-":
		if goesTo != East && goesTo != West {
			return false, 0, nil, goesTo
		}
	case "|":
		if goesTo != North && goesTo != South {
			return false, 0, nil, goesTo
		}
	case "L":
		if goesTo == South {
			goesTo = East
		} else if goesTo == West {
			goesTo = North
		} else {
			return false, 0, nil, goesTo
		}
	case "J":
		if goesTo == East {
			goesTo = North
		} else if goesTo == South {
			goesTo = West
		} else {
			return false, 0, nil, goesTo
		}
	case "7":
		if goesTo == East {
			goesTo = South
		} else if goesTo == North {
			goesTo = West
		} else {
			return false, 0, nil, goesTo
		}
	case "F":
		if goesTo == North {
			goesTo = East
		} else if goesTo == West {
			goesTo = South
		} else {
			return false, 0, nil, goesTo
		}
	case "S":
		if !start {
			path := make(map[int]map[int]string)
			path[current.Y] = make(map[int]string)
			path[current.Y][current.X] = symbol
			return true, 0, path, goesTo
		}
	default:
		return false, 0, nil, goesTo
	}

	var nextPosition Position
	switch goesTo {
	case East:
		nextPosition = Position{X: current.X + 1, Y: current.Y}
	case West:
		nextPosition = Position{X: current.X - 1, Y: current.Y}
	case North:
		nextPosition = Position{X: current.X, Y: current.Y - 1}
	case South:
		nextPosition = Position{X: current.X, Y: current.Y + 1}
	}

	found, sum, path, lastDirection := walkPath(nextPosition, goesTo, false, tiles)

	if found {
		if _, ok := path[current.Y]; !ok {
			path[current.Y] = make(map[int]string)
		}
		path[current.Y][current.X] = symbol
		return true, sum + 1, path, lastDirection
	}

	return false, 0, nil, goesTo
}

func parseLines(lines []string) [][]string {
	var tiles [][]string

	width := utf8.RuneCountInString(lines[0]) + 2
	tiles = append(tiles, strings.Split(strings.Repeat(".", width), ""))

	for _, line := range lines {
		line = "." + line + "."
		tiles = append(tiles, strings.Split(line, ""))
	}

	tiles = append(tiles, tiles[0])

	return tiles
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
