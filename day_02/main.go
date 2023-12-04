package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type TokenType string;

const (
    GAME = "Game"
    RED = "RED"
    GREEN = "GREEN"
    BLUE = "BLUE"
    INVALID_IDENTIFIER = "INVALID_IDENTIFIER"
    COLON = ";"
    SEMICOLON = ";"
    COMMA = ","
    INT = "Int"
    EOL = "EOL"
)

type Token struct {
    Type   TokenType
    Literal string
}

var keywords = map[string]TokenType {
    "Game":   GAME,
    "red":    RED,
    "green":  GREEN,
    "blue":   BLUE,
}

func lookupIdent(ident []rune) TokenType {
    if tok, ok := keywords[string(ident)]; ok {
        return tok
    }
    return INVALID_IDENTIFIER
}

func isLetter(char rune) bool {
    return (char >= 65 && char <= 90) || (char >= 97 && char <= 122)
}

func isDigit(char rune) bool {
    return char >= 48 && char <= 57
}

type Lexer struct {
    input           []rune
    position        int
    readPosition    int
    char            rune
}

func newLexer(input []rune) *Lexer {
    lexer := Lexer{input: input}
    lexer.readChar()
    return &lexer;
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

func (l *Lexer) readIdentifier() []rune {
    position := l.position
    for isLetter(l.char) {
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
        case ':':
            l.readChar()
            return newToken(COLON, char)
        case ';':
            l.readChar()
            return newToken(SEMICOLON, char)
        case ',':
            l.readChar()
            return newToken(COMMA, char)
        case 0:
            return newToken(EOL, 0)
        default:
            if (isLetter(char)) {
                identifier := l.readIdentifier()
                tokenType := lookupIdent(identifier)
                if tokenType != INVALID_IDENTIFIER {
                    return Token{Type: tokenType, Literal: string(identifier)}
                } else {
                    panic(fmt.Errorf("Unknown identifier '%d'", identifier))
                }
            } else if  isDigit(char) {
                number := l.readNumber()
                return Token{Type: INT, Literal: string(number)}
            }
            panic(fmt.Errorf("Unknown char '%d'", char))
    }
}

func newToken(tokenType TokenType, char rune) Token {
    return Token{Type: tokenType, Literal: string(char)}
}

type Color string;

const (
    COLOR_RED = "red"
    COLOR_GREEN = "green"
    COLOR_BLUE = "blue"
)

var colors = map[string]Color {
    "red": COLOR_RED,
    "green": COLOR_GREEN,
    "blue": COLOR_BLUE,
}

type RevealStatement struct {
    Color   Color
    Amount  int
}

type SubsetStatement struct {
    Reveals []RevealStatement
}

type GameStatement struct {
    GameId  int
    Subsets []SubsetStatement
}

type Parser struct {
    lexer *Lexer
    curToken Token
    peekToken Token
}

func newParser(lexer *Lexer) *Parser {
    p := &Parser{lexer: lexer}

    p.nextToken()
    p.nextToken()

    return p
}

func parseInt(s string) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        panic(err)
    }
    return i
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.lexer.nextToken()
}

func (p *Parser) curTokenIs(tokenType TokenType) bool {
    return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType TokenType) bool {
    return p.peekToken.Type == tokenType
}

func (p *Parser) expectPeekToken(tokenType TokenType) {
    if !p.peekTokenIs(tokenType) {
        panic(fmt.Errorf("expected token '%s', got '%s'", tokenType, p.peekToken.Type))
    }
    p.nextToken()
}

func (p *Parser) parseGame() *[]GameStatement {
    statements := []GameStatement{}
    for !p.curTokenIs(EOL) {
        statement := p.parseGameStatement()
        statements = append(statements, statement)
    }
    return &statements
}

func (p *Parser) parseGameStatement() GameStatement {
    if !p.curTokenIs(GAME) {
        panic(fmt.Errorf("unexpected statement '%s'", p.curToken.Type))
    }
    p.expectPeekToken(INT)
    id := parseInt(p.curToken.Literal)
    p.expectPeekToken(COLON)
    subsets := p.parseSubsetStatements()
    return GameStatement{GameId: id, Subsets: subsets}
}

func (p *Parser) parseSubsetStatements() []SubsetStatement {
    subsets := []SubsetStatement{}
    for !p.curTokenIs(EOL) {
        subsets = append(subsets, p.parseSubsetStatement())
    }
    return subsets
}

func (p *Parser) parseSubsetStatement() SubsetStatement {
    reveals := []RevealStatement{}
    for !p.curTokenIs(EOL) {
        reveals = append(reveals, p.parseRevealStatement())
        if p.curTokenIs(SEMICOLON) {
            break;
        }
    }
    return SubsetStatement{Reveals: reveals}
}

func (p *Parser) parseRevealStatement() RevealStatement {
    p.expectPeekToken(INT)
    amount := parseInt(p.curToken.Literal)
    if p.peekTokenIs(GREEN) || p.peekTokenIs(RED) || p.peekTokenIs(BLUE) {
        p.nextToken()
    }
    color, ok := colors[p.curToken.Literal]
    if !ok {
        panic(fmt.Errorf("unexpected color '%s'", p.curToken.Literal))
    }
    p.nextToken()
    return RevealStatement{Color: color, Amount: amount}
}

type Subset struct {
    Red     int
    Green   int
    Blue    int
}

type Game struct {
    GameId  int
    Subsets []Subset
}

func parseGame(line string) Game {
    lexer := newLexer([]rune(line))
    parser := newParser(lexer)
    statements := parser.parseGame()
    if len(*statements) == 0 {
        panic(fmt.Errorf("got no statements"))
    }
    gameStatement := (*statements)[0]
    subsets := []Subset{}
    for _, statement := range gameStatement.Subsets {
        subset := Subset{}
        for _, revealStatement := range statement.Reveals {
            switch(revealStatement.Color) {
            case COLOR_RED:
                subset.Red += revealStatement.Amount
            case COLOR_GREEN:
                subset.Green += revealStatement.Amount
            case COLOR_BLUE:
                subset.Blue += revealStatement.Amount
            }
        }
        subsets = append(subsets, subset)
    }
    return Game{GameId: gameStatement.GameId, Subsets: subsets}
}

func isValidGame(game Game, predicate Subset) bool {
    for _, subset := range game.Subsets {
        if subset.Red > predicate.Red || subset.Green > predicate.Green || subset.Blue > predicate.Blue {
            return false
        }
    }
    return true
}

func maxInt(i1 int, i2 int) int {
    if i1 > i2 {
        return i1
    }
    return i2
}

func calculatePower(subsets []Subset) int {
    red := 0
    green := 0
    blue := 0
    for _, subset := range subsets {
        red = maxInt(red, subset.Red)
        green = maxInt(green, subset.Green)
        blue = maxInt(blue, subset.Blue)
    }
    return red * green * blue
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    sum := 0
    sumOfPowers := 0
    predicate := Subset{Red: 12, Green: 13, Blue: 14}
    for scanner.Scan() {
        line := scanner.Text()
        game := parseGame(line)
        if isValidGame(game, predicate) {
            sum += game.GameId
        }
        sumOfPowers += calculatePower(game.Subsets)
    }
    fmt.Printf("%d\n", sum)
    fmt.Printf("%d\n", sumOfPowers)
}

