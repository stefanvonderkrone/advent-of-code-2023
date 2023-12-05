package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/exp/slices"
)

const (
    PIPE = "PIPE"
    EOL = "EOL"
    INT = "INT"
)

type TokenType string

type Token struct {
    Type TokenType
    Literal []rune
}

func newToken(tokenType TokenType, literal rune) Token {
    return Token{Type: tokenType, Literal: []rune{literal}}
}

func isDigit(char rune) bool {
    return char >= 48 && char <= 57
}

type Lexer struct {
    input []rune
    position int
    readPosition int
    char rune
}

func newLexer(input string) *Lexer {
    lexer := Lexer{input: []rune(input)}
    lexer.readChar()
    for lexer.char != ':' {
        lexer.readChar()
    }
    lexer.readChar()
    return &lexer
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.char = 0;
    } else {
        l.char = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition += 1
}

func (l *Lexer) readNumber() []rune {
    position := l.position
    for isDigit(l.char) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
    for l.char == ' ' {
        l.readChar()
    }
}

func (l *Lexer) nextToken() Token {
    l.skipWhitespace()

    char := l.char

    switch(char) {
        case '|':
            l.readChar()
            return newToken(PIPE, char)
        case 0:
            return newToken(EOL, char)
        default:
            if isDigit(char) {
                number := l.readNumber()
                return Token{Type: INT, Literal: number}
            }
    }
    panic(fmt.Errorf("Unknown char '%s'", string(char)))
}

func (l *Lexer) parseNumbers() []int {
    token := l.nextToken()
    numbers := []int{}
    for token.Type == INT {
        if n, err := strconv.Atoi(string(token.Literal)); err == nil {
            numbers = append(numbers, n)
        }
        token = l.nextToken()
    }
    return numbers
}

type Card struct {
    winning []int
    owning []int
}

func (c *Card) calculatePower() int {
    power := 0
    for _, n := range c.owning {
        if slices.Contains(c.winning, n) {
            if power == 0 {
                power = 1
            } else {
                power *= 2
            }
        }
    }
    return power
}

func (c *Card) calculateScore() int {
    score := 0
    for _, n := range c.owning {
        if slices.Contains(c.winning, n) {
            score += 1
        }
    }
    return score
}

type Scoreboard struct {
    counts []int
}

func newScoreboard() Scoreboard {
    return Scoreboard{counts: []int{}}
}

func (s *Scoreboard) addScoreAt(score int, at int) {
    requiredLength := at + score + 1
    for len(s.counts) < requiredLength {
        s.counts = append(s.counts, 1)
    }
    count := s.counts[at]
    fmt.Printf("score %d at %d with count %d\n", score, at, count)
    for n := 0; n < count; n++ {
        for i := at + 1; i < requiredLength; i++ {
            s.counts[i]++;
        }
    }
}

func (s *Scoreboard) sum() int {
    sum := 0
    for _, n := range s.counts {
        sum += n
    }
    return sum
}

func parseLine(line string) Card {
    lexer := newLexer(line)
    winning := lexer.parseNumbers()
    owning := lexer.parseNumbers()
    return Card{winning: winning, owning: owning}
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    sum := 0
    scoreboard := newScoreboard()
    index := 0
    for scanner.Scan() {
        line := scanner.Text()
        card := parseLine(line)
        power := card.calculatePower()
        sum += power
        score := card.calculateScore()
        scoreboard.addScoreAt(score, index)
        index++
    }
    fmt.Printf("%d\n", sum)
    fmt.Printf("%d\n", scoreboard.sum())
}
