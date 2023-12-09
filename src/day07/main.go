package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type HandType int

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type Hand struct {
	Cards    []rune
	Labels   map[rune]int
	Type     HandType
	Strength int64
	Bid      int
}

func main() {
	lines, err := readLinesFromFile("../../inputs/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	winnings := part1(lines)
	fmt.Println("Part 1:", winnings)
}

func part1(lines []string) int {
	hands := parseLines(lines)

	sort.Slice(hands, func(i, j int) bool { return hands[i].Strength < hands[j].Strength })

	winnings := 0
	for index, hands := range hands {
		winnings = winnings + hands.Bid*(index+1)
	}

	return winnings
}

func parseLines(lines []string) []Hand {
	var hands []Hand

	for _, line := range lines {
		hand := parseLine(line)
		hands = append(hands, hand)
	}

	return hands
}

func parseLine(line string) Hand {
	hand := Hand{
		Labels: make(map[rune]int),
	}

	cardsString, bidString, _ := strings.Cut(line, " ")

	bid, _ := strconv.Atoi(bidString)
	hand.Bid = bid

	for _, card := range cardsString {
		hand.Cards = append(hand.Cards, card)

		if labelCount, ok := hand.Labels[card]; ok {
			hand.Labels[card] = labelCount + 1
		} else {
			hand.Labels[card] = 1
		}
	}

	handType := getHandType(hand.Labels)
	hand.Type = handType

	strength := getHandStrength(hand)
	hand.Strength = strength

	return hand
}


func getHandStrength(hand Hand) int64 {
	strength := strconv.Itoa(int(hand.Type))

	for _, card := range hand.Cards {
		if unicode.IsDigit(card) {
			strength = strength + "0" + string(card)
			continue
		}

		cardValue := 0
		switch card {
		case 'A':
			cardValue = 14
		case 'K':
			cardValue = 13
		case 'Q':
			cardValue = 12
		case 'J':
			cardValue = 11
		case 'T':
			cardValue = 10
		default:
			log.Fatal("Invalid card label")
		}

		strength = strength + strconv.Itoa(cardValue)
	}

	value, _ := strconv.ParseInt(strength, 10, 64)
	return value
}

func getHandType(labels map[rune]int) HandType {

	singleCount, pairCount, trioCount := 0, 0, 0
	for _, count := range labels {
		switch count {
		case 5:
			return FiveOfAKind
		case 4:
			return FourOfAKind
		case 3:
			trioCount = 1
		case 2:
			pairCount = pairCount + 1
		case 1:
			singleCount = singleCount + 1
		default:
			log.Fatal("Invalid label count")
		}

		if trioCount == 1 && pairCount == 1 {
			return FullHouse
		} else if trioCount == 1 && singleCount >= 1 {
			return ThreeOfAKind
		} else if pairCount == 2 {
			return TwoPair
		} else if pairCount == 1 && singleCount >= 3 {
			return OnePair
		}
	}

	return HighCard
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
