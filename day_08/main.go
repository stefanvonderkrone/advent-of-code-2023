package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Pair struct {
    Left string
    Right string
}

func isLetter(char byte) bool {
    return (char >= 48 && char <= 57) || (char >= 65 && char <= 90)
}

func readWordAt(line string, at int) (string, int) {
    start := at
    for !isLetter(line[start]) {
        start++
    }
    end := start
    for isLetter(line[end]) {
        end++;
    }
    return line[start:end], end
}

func parseLine(line string) (string, Pair) {
    key, end := readWordAt(line, 0)
    pair := Pair{}
    left, endLeft := readWordAt(line, end)
    right, _ := readWordAt(line, endLeft)
    pair.Left = left
    pair.Right = right
    return key, pair
}

func solvePt1(instructions []rune, coordinates map[string]Pair) {
    key := "AAA"
    current, _ := coordinates[key]
    steps := 0
    instructionIndex := 0
    numInstructions := len(instructions)
    for key != "ZZZ" {
        // fmt.Printf("current %+v, %s\n", current, key)
        steps++
        instruction := instructions[instructionIndex]
        if instruction == 'L' {
            key = current.Left
        } else {
            key = current.Right
        }
        current, _ = coordinates[key]
        instructionIndex++
        if instructionIndex >= numInstructions {
            instructionIndex = 0
        }
    }
    fmt.Printf("%d\n", steps)
}

func endsWith(word string, char byte) bool {
    return word[len(word) - 1] == char
}

func allEndWith(words []string, char byte) bool {
    for _, word := range words {
        if !endsWith(word, char) {
            return false
        }
    }
    return true
}

func stepsFrom(instructions []rune, coordinates map[string]Pair, from string) int {
    steps := 0
    instructionIndex := 0
    numInstructions := len(instructions)
    key := from
    for !endsWith(key, 'Z') {
        steps++
        instruction := instructions[instructionIndex]
        pair, _ := coordinates[key]
        if instruction == 'L' {
            key = pair.Left
        } else {
            key = pair.Right
        }
        instructionIndex++
        if instructionIndex >= numInstructions {
            instructionIndex = 0
        }
    }
    return steps
}

func gcd(a int, b int) int {
    if a > b {
        b, a = a, b
    }
    if a == 0 {
        return b
    }
    return gcd(a, b % a)
}

func lcm(a int, b int) int {
    return int(math.Abs(float64(a * b))) / gcd(a, b)
}

func reduceLcm(xs []int) int {
    a := xs[0]
    b := xs[1]
    a = lcm(a, b)
    xs = xs[2:]
    for len(xs) > 0 {
        b = xs[0]
        xs = xs[1:]
        a = lcm(a, b)
    }
    return a
}

func solvePt2(instructions []rune, coordinates map[string]Pair) {
    keys := []string{}
    steps := []int{}
    for key := range coordinates {
        if endsWith(key, 'A') {
            keys = append(keys, key)
            steps = append(steps, stepsFrom(instructions, coordinates, key))
        }
    }
    fmt.Printf("%+v\n", steps)
    fmt.Printf("%d\n", reduceLcm(steps))
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    index := 0
    instructions := []rune{}
    coordinates := map[string]Pair{}
    for scanner.Scan() {
        line := scanner.Text()
        if index == 0 {
            instructions = []rune(line)
        }
        if index > 1 {
            key, pair := parseLine(line)
            coordinates[key] = pair
        }
        index++
    }
    // solvePt1(instructions, coordinates)
    solvePt2(instructions, coordinates)
    // fmt.Printf("%s\n", string(instructions))
    // fmt.Printf("%+v\n", coordinates)
}
