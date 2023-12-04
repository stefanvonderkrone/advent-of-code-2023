package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Stack3L struct {
    upper []rune
    middle []rune
    lower []rune
}

func newStack3L() Stack3L {
    return Stack3L{}
}

func (s *Stack3L) move() {
    s.upper = s.middle
    s.middle = s.lower
    s.lower = nil
}

func (s *Stack3L) push(line string) {
    s.move()
    s.lower = []rune(line)
}

const DOT = rune('.')

func isDigit(char rune) bool {
    return char >= 48 && char <= 57
}

func isSymbol(char rune) bool {
    if isDigit(char) || char == DOT {
        return false
    }
    return true
}

func readNumberAt(line []rune, at int) int {
    lastIndex := len(line) - 1
    left := at
    for left > 0 && isDigit(line[left - 1]) {
        left--
    }
    right := at
    for right < lastIndex && isDigit(line[right + 1]) {
        right++
    }
    numberSlice := line[left:right + 1]
    numberString := string(numberSlice)
    // fmt.Printf("found number: '%s'\n", numberString)
    n, err := strconv.Atoi(string(numberString));
    if err != nil {
        return 0
    }
    return n
}

func extractLineAt(line []rune, at int) []int {
    numbers := []int{}
    // look at the center
    if isDigit(line[at]) {
        numbers = append(numbers, readNumberAt(line, at))
        return numbers
    }
    if at > 0 && isDigit(line[at - 1]) {
        numbers = append(numbers, readNumberAt(line, at - 1))
    }
    if at < len(line) - 2 && isDigit(line[at + 1]) {
        numbers = append(numbers, readNumberAt(line, at + 1))
    }
    return numbers
}

func (s *Stack3L) extractAt(at int) []int {
    numbers := []int{}
    // check upper
    if s.upper != nil {
        numbers = append(numbers, extractLineAt(s.upper, at)...)
    }
    if s.middle != nil {
        numbers = append(numbers, extractLineAt(s.middle, at)...)
    }
    if s.lower != nil {
        numbers = append(numbers, extractLineAt(s.lower, at)...)
    }
    return numbers
}

func (s *Stack3L) extractCurent() []int {
    numbers := []int{}
    if s.middle == nil {
        return numbers
    }
    for index, char := range s.middle {
        if isSymbol(char) {
            numbers = append(numbers, s.extractAt(index)...)
            // fmt.Printf("found numbers arround symbol: %v\n", numbers)
            // current = append(current, s.extractAt(index)...)
            // fmt.Printf("char '%s' at '%d' is a symbol\n", string(char), index)
        }
    }
    return numbers
}

func (s *Stack3L) extractGearRationAt(index int) int {
    numbers := s.extractAt(index)
    if len(numbers) == 2 {
        a := numbers[0]
        b := numbers[1]
        return a * b
    }
    return 0
}

const GEAR = rune('*')

func (s *Stack3L) extractCurrentGearRatio() []int {
    numbers := []int{}
    if s.middle == nil {
        return numbers
    }
    for index, char := range s.middle {
        if char == GEAR {
            ratio := s.extractGearRationAt(index);
            if ratio > 0 {
                numbers = append(numbers, ratio)
            }
        }
    }
    return numbers
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    stack := newStack3L()
    numbers := []int{}
    gearRatios := []int{}
    for scanner.Scan() {
        line := scanner.Text()
        stack.push(line)
        numbers = append(numbers, stack.extractCurent()...)
        gearRatios = append(gearRatios, stack.extractCurrentGearRatio()...)
    }
    stack.move()
    numbers = append(numbers, stack.extractCurent()...)
    gearRatios = append(gearRatios, stack.extractCurrentGearRatio()...)
    // fmt.Printf("numbers: %v\n", numbers)
    sum := 0
    for _, number := range numbers {
        sum += number
    }
    fmt.Printf("%d\n", sum)
    sumGR := 0
    for _, gearRatio := range gearRatios {
        sumGR += gearRatio
    }
    fmt.Printf("%d\n", sumGR)
}
