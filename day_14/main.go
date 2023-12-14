package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func replaceCharAt(line string, char rune, at int) string {
    runes := []rune(line)
    runes[at] = char
    return string(runes)
}

var fallthroughCache = map[string][]byte{}

func fallThrough(row []byte) []byte {
    key := string(row)
    for x, byte := range row {
        if byte == 'O' && x > 0 {
            xx := x
            for xx > 0 {
                xx--
                if row[xx] != '.' {
                    xx++
                    break
                }
            }
            row[x] = '.'
            row[xx] = 'O'
        }
    }
    fallthroughCache[key] = row
    return row
}

func tiltNorth(array [][]byte) [][]byte {
    numRows := len(array)
    rowLength := len(array[0])
    for x := 0; x < rowLength; x++ {
        row := make([]byte, numRows)
        for y := 0; y < numRows; y++ {
            row[y] = array[y][x]
        }
        row = fallThrough(row)
        for y, byte := range row {
            array[y][x] = byte
        }
    }
    return array
}

func tiltSouth(array [][]byte) [][]byte {
    numRows := len(array)
    rowLength := len(array[0])
    for x := 0; x < rowLength; x++ {
        row := make([]byte, numRows)
        for y := 0; y < numRows; y++ {
            row[numRows - 1 - y] = array[y][x]
        }
        row = fallThrough(row)
        for y, byte := range row {
            array[numRows - 1 - y][x] = byte
        }
    }
    return array
}

func tiltWest(array [][]byte) [][]byte {
    numRows := len(array)
    for y := 0; y < numRows; y++ {
        array[y] = fallThrough(array[y])
    }
    return array
}

func tiltEast(array [][]byte) [][]byte {
    numRows := len(array)
    for y := 0; y < numRows; y++ {
        row := array[y]
        rowLength := len(row)
        tmpRow := make([]byte, rowLength)
        for x, byte := range row {
            tmpRow[rowLength - 1 - x] = byte
        }
        tmpRow = fallThrough(tmpRow)
        for x, byte := range tmpRow {
            row[rowLength - 1 - x] = byte
        }
    }
    return array
}

var tiltCache = map[string][][]byte{}

func byteToString(array [][]byte) string {
    builder := bytes.Buffer{}
    for _, bs := range array {
        builder.Write(bs)
    }
    return builder.String()
}

func tiltCycle(array [][]byte) [][]byte {
    key := byteToString(array)
    if r, ok := tiltCache[key]; ok {
        return r
    }
    array = tiltNorth(array)
    array = tiltWest(array)
    array = tiltSouth(array)
    array = tiltEast(array)
    tiltCache[key] = array
    return array
}

func calc(array [][]byte) int {
    lineNo := len(array)
    sum := 0
    for _, row := range array {
        for _, byte := range row {
            if byte == 'O' {
                sum += lineNo
            }
        }
        lineNo--
        // fmt.Printf("%s\n", string(row))
    }
    return sum
}

func solvePt1(array [][]byte) {
    // for _, row := range array {
    //     fmt.Printf("%s\n", row)
    // }
    // fmt.Print("-------\n")
    array = tiltNorth(array)
    // fmt.Print("-------\n")
    sum := calc(array)
    fmt.Printf("%d\n", sum)
}

func solvePt2(array [][]byte) {
    cycles := 1000000000
    for i := 0; i < cycles; i++ {
        array = tiltCycle(array)
    }
    printRows(array)
    sum := calc(array)
    fmt.Printf("%d\n", sum)
    fmt.Printf("%d\n", len(tiltCache))
}

func printRows(array [][]byte) {
    fmt.Print("+")
    for range array {
        fmt.Print("-")
    }
    fmt.Print("+\n")
    for _, row := range array {
        fmt.Printf("|%s|\n", row)
    }
    fmt.Print("+")
    for range array {
        fmt.Print("-")
    }
    fmt.Print("+\n")
}

func printLines(lines []string) {
    for _, row := range lines {
        fmt.Printf("%s\n", row)
    }
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    builder := strings.Builder{}
    for scanner.Scan() {
        if builder.Len() > 0 {
            builder.WriteByte('\n')
        }
        builder.Write(scanner.Bytes())
    }
    s := builder.String()
    ls := strings.Split(s, "\n")
    lines := make([][]byte, len(ls))

    for i, line := range ls {
        lines[i] = []byte(line)
    }

    // printRows(array)
    // printLines(ls)
    printRows(lines)
    solvePt1(lines)
    // solvePt2(lines)
}
