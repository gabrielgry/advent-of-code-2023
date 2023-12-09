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
	winnings2 := part2(lines)
	fmt.Println("Part 2:", winnings2)
}

func part1(lines []string) int {
	var hands []Hand

	for _, line := range lines {
		hand := parseLine(line)
		hands = append(hands, hand)
	}

	sort.Slice(hands, func(i, j int) bool { return hands[i].Strength < hands[j].Strength })

	winnings := 0
	for index, hands := range hands {
		winnings = winnings + hands.Bid*(index+1)
	}

	return winnings
}

func part2(lines []string) int {
	var hands []Hand

	for _, line := range lines {
		hand := parseLineWithJoker(line)
		hands = append(hands, hand)
	}

	sort.Slice(hands, func(i, j int) bool { return hands[i].Strength < hands[j].Strength })

	winnings := 0
	for index, hands := range hands {
		winnings = winnings + hands.Bid*(index+1)
	}

	return winnings
}

func parseLineWithJoker(line string) Hand {
	hand := Hand{
		Labels: make(map[rune]int),
	}

	cardsString, bidString, _ := strings.Cut(line, " ")
	fmt.Println(cardsString)

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

	morthedLabels := morthJoker(hand.Labels)
	handType := getHandType(morthedLabels)
	hand.Type = handType

	strength := getHandStrengthWithJoker(hand)
	hand.Strength = strength

	return hand
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

func morthJoker(labels map[rune]int) map[rune]int {
	jokerCount, hasJoker := labels['J']

	if !hasJoker || jokerCount >= 5 {
		return labels
	}

	delete(labels, 'J')

	var sortedKeys []rune

	for key := range labels {
		sortedKeys = append(sortedKeys, key)
	}

	sort.SliceStable(sortedKeys, func(i, j int) bool {
		iStrength := getCardStrengthWithJoker(sortedKeys[i]) + (labels[sortedKeys[i]] * 100)
		jStrength := getCardStrengthWithJoker(sortedKeys[j]) + (labels[sortedKeys[j]] * 100)
		return iStrength > jStrength
	})

	labels[sortedKeys[0]] = labels[sortedKeys[0]] + jokerCount

	return labels
}

func getCardStrengthWithJoker(card rune) int {
	if unicode.IsDigit(card) {
		value, _ := strconv.Atoi(string(card))
		return value
	}

	switch card {
	case 'A':
		return 14
	case 'K':
		return 13
	case 'Q':
		return 12
	case 'J':
		return 1
	case 'T':
		return 10
	default:
		log.Fatal("Invalid card label")
	}

	return 0
}

func getCardStrength(card rune) int {
	if unicode.IsDigit(card) {
		value, _ := strconv.Atoi(string(card))
		return value
	}

	switch card {
	case 'A':
		return 14
	case 'K':
		return 13
	case 'Q':
		return 12
	case 'J':
		return 11
	case 'T':
		return 10
	default:
		log.Fatal("Invalid card label")
	}

	return 0
}

func getHandStrengthWithJoker(hand Hand) int64 {
	strength := strconv.Itoa(int(hand.Type))

	for _, card := range hand.Cards {
		cardStrength := getCardStrengthWithJoker(card)
		cardStrengthFormated := fmt.Sprintf("%0*d", 2, cardStrength)
		strength = strength + cardStrengthFormated
	}

	value, _ := strconv.ParseInt(strength, 10, 64)
	return value
}

func getHandStrength(hand Hand) int64 {
	strength := strconv.Itoa(int(hand.Type))

	for _, card := range hand.Cards {
		cardStrength := getCardStrength(card)
		cardStrengthFormated := fmt.Sprintf("%0*d", 2, cardStrength)
		strength = strength + cardStrengthFormated
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
