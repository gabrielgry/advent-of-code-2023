package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	Id             int
	WinningNumbers []int
	OwnNumbers     []int
	WinningCount   int
	Points         int
}

type CardStack struct {
	Card     Card
	Quantity int
}

type CardPool map[int]CardStack

func countCards(cardPool CardPool) int {
	count := 0

	for _, cardStack := range cardPool {
		count = count + cardStack.Quantity
	}

	return count
}

func processCardPool(cardPool CardPool) CardPool {
	maxIdCard := len(cardPool)

	for idCard := 1; idCard <= maxIdCard; idCard = idCard + 1 {
		startClonningId := idCard + 1
		endClonningId := idCard + cardPool[idCard].Card.WinningCount

		if startClonningId > maxIdCard {
			break
		}

		for cloneIdCard := startClonningId; (cloneIdCard <= endClonningId) && (cloneIdCard <= maxIdCard); cloneIdCard = cloneIdCard + 1 {
			if clonedCardStack, ok := cardPool[cloneIdCard]; ok {

				cloneQuantity := cardPool[idCard].Quantity
				clonedCardStack.Quantity = clonedCardStack.Quantity + cloneQuantity

				cardPool[cloneIdCard] = clonedCardStack
			}
		}
	}

	return cardPool
}

func createCardPool(cards []Card) CardPool {
	cardPool := make(CardPool)

	for _, card := range cards {
		cardPool[card.Id] = CardStack{Card: card, Quantity: 1}
	}

	return cardPool
}

func getTotalPoints(cards []Card) int {
	sum := 0

	for _, card := range cards {
		sum = sum + card.Points
	}

	return sum
}

func countCardPoints(winningNumbers []int, ownNumbers []int) (int, int) {
	count := 0
	points := 0

	for _, winningNumber := range winningNumbers {
		for _, ownNumber := range ownNumbers {
			if winningNumber == ownNumber {
				if points == 0 {
					points = 1
				} else {
					points = points + points
				}

				count = count + 1

				break
			}
		}
	}

	return count, points
}

func fieldsStringToNumber(fields []string) []int {
	var intFields []int

	for _, field := range fields {
		intField, _ := strconv.Atoi(field)
		intFields = append(intFields, intField)
	}

	return intFields
}

func parseLine(line string) Card {
	head, numbersString, _ := strings.Cut(line, ":")
	winningString, ownString, _ := strings.Cut(numbersString, "|")

	headFields := strings.Fields(head)

	cardId, _ := strconv.Atoi(headFields[1])
	winningNumbers := fieldsStringToNumber(strings.Fields(winningString))
	ownNumbers := fieldsStringToNumber(strings.Fields(ownString))

	count, points := countCardPoints(winningNumbers, ownNumbers)

	return Card{
		Id:             cardId,
		WinningNumbers: winningNumbers,
		OwnNumbers:     ownNumbers,
		WinningCount:   count,
		Points:         points,
	}
}

func parseLines(lines []string) []Card {
	var cards []Card

	for _, line := range lines {
		cards = append(cards, parseLine(line))
	}

	return cards
}

func main() {
	lines, err := readLinesFromFile("../../inputs/day04/input.txt")

	if err != nil {
		log.Fatal("Could not open the input file")
	}

	cards := parseLines(lines)

	sum := getTotalPoints(cards)
	fmt.Println("Sum:", sum)

	cardPool := createCardPool(cards)
	cardPool = processCardPool(cardPool)
	cardsCount := countCards(cardPool)

	fmt.Println("Cards count:", cardsCount)

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
