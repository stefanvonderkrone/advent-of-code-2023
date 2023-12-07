package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type HandType int

const (
    FIVE_OF_A_KIND = 10000
    FOUR_OF_A_KIND = 1001
    FULL_HOUSE = 110
    THREE_OF_A_KIND = 102
    TWO_PAIR = 31
    ONE_PAIR = 13
    HIGH_CARD = 5
)

type CardType int

const (
    ASS = 14
    KING = 13
    QUEEN = 12
    JACK = 11
    TEN = 10
    NINE = 9
    EIGHT = 8
    SEVEN = 7
    SIX = 6
    FIVE = 5
    FOUR = 4
    THREE = 3
    TWO = 2
)

var cardsMap = map[rune]CardType {
    'A': ASS,
    'K': KING,
    'Q': QUEEN,
    'J': JACK,
    'T': TEN,
    '9': NINE,
    '8': EIGHT,
    '7': SEVEN,
    '6': SIX,
    '5': FIVE,
    '4': FOUR,
    '3': THREE,
    '2': TWO,
}

func handType(hand []CardType) HandType {
    cardCounts := make([]int, 15)
    for _, cardType := range hand {
        cardCounts[cardType]++
    }
    counts := make([]int, len(hand))
    for _, count := range cardCounts {
        if count > 0 {
            counts[count - 1]++
        }
    }
    hash := 0
    for i, n := range counts {
        hash += n * int(math.Pow10(i))
    }
    return HandType(hash)
}

func parseHand(hand []rune) []CardType {
    cards := make([]CardType, len(hand))
    for index, card := range hand {
        if cardType, ok := cardsMap[hand[index]]; ok {
            cards[index] = cardType
        } else {
            panic(fmt.Errorf("Unknown CardType for '%s'", string([]rune{card})))
        }
    }
    return cards
}

func parseCard(line string) Card {
    card := Card{}
    lineR := []rune(line)
    index := 0
    for lineR[index] != ' ' {
        index++;
    }
    card.Hand = parseHand(lineR[0:index])
    index++
    bid, err := strconv.Atoi(string(lineR[index:]))
    if err != nil {
        panic(err)
    }
    card.Bid = bid
    card.Type = handType(card.Hand)
    return card
}

type Card struct {
    Hand []CardType
    Bid int
    Type HandType
}

type Cards []Card

func (c Cards) Len() int {
    return len(c)
}

func (c Cards) Swap(i, j int) {
    c[j], c[i] = c[i], c[j]
}

func (c Cards) Less(i, j int) bool {
    card1 := c[i]
    card2 := c[j]
    if card1.Type == card2.Type {
        for i := 0; i < len(card1.Hand); i++ {
            if card1.Hand[i] == card2.Hand[i] {
                continue
            }
            return card1.Hand[i] < card2.Hand[i]
        }
    }
    return card1.Type < card2.Type
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    cards := []Card{}
    for scanner.Scan() {
        line := scanner.Text()
        card := parseCard(line);
        cards = append(cards, card)
    }
    sort.Sort(Cards(cards))
    sum := 0
    for i, card := range cards {
        sum += (i+1) * card.Bid
    }
    fmt.Printf("%d\n", sum)
}
