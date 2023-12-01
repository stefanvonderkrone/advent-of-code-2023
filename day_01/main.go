package main

import (
	"bufio"
	"fmt"
	"os"
)

var numbers = map[rune]int{
    '1': 1,
    '2': 2,
    '3': 3,
    '4': 4,
    '5': 5,
    '6': 6,
    '7': 7,
    '8': 8,
    '9': 9,
}

var numberNames = map[string]int{
    "one": 1,
    "two": 2,
    "three": 3,
    "four": 4,
    "five": 5,
    "six": 6,
    "seven": 7,
    "eight": 8,
    "nine": 9,
}

func readNameAt(line string, index int) (int, bool) {
    chars := []rune(line)
    numChars := len(chars)
    for name, num := range numberNames {
        nameLength := len([]rune(name))
        endIndex := index + nameLength
        if endIndex > numChars {
            endIndex = numChars
        }
        n := string(chars[index:endIndex])
        if n != name {
            continue
        }
        if _, ok := numberNames[n]; ok {
            return num, true
        }
    }
    return -1, false
}

func readFirstNum(line string) int {
    for index, char := range line {
        if num, ok := numbers[char]; ok {
            return num
        }
        if num, ok := readNameAt(line, index); ok {
            return num
        }
    }
    return 0
}

func readLastNum(line string) int {
    chars := []rune(line)
    length := len(chars) - 1
    for i := length; i >= 0; i-- {
        char := chars[i]
        if num, ok := numbers[char]; ok {
            return num
        }
        if num, ok := readNameAt(line, i); ok {
            return num
        }
    }
    return 0;
}

func readCalibration(line string) int {
    firstNum := readFirstNum(line)
    lastNum := readLastNum(line)
    return firstNum * 10 + lastNum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

    sum := 0
    for scanner.Scan() {
        line := scanner.Text()
        calibration := readCalibration(line)
        sum = sum + calibration
    }
    fmt.Printf("%d\n", sum)
}

