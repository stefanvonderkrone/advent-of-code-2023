package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Galaxy struct {
    X int
    Y int
}

func parseRow(line string, y int) ([]Galaxy, bool) {
    isRowEmpty := true
    row := []Galaxy{}
    for i := 0; i < len(line); i++ {
        if line[i] == '.' {
            continue
        }
        galaxy := Galaxy{X: i, Y: y}
        row = append(row, galaxy)
        isRowEmpty = false
    }
    return row, isRowEmpty
}

func length(g1 Galaxy, g2 Galaxy) int {
    return int(math.Abs(float64(g1.X - g2.X))) + int(math.Abs(float64(g1.Y - g2.Y)))
}

func solve(universe []Galaxy, rowDeltas []int) int {
    universeSize := len(universe)
    sum := 0
    for i := 0; i < universeSize; i++ {
        g1 := universe[i]
        for k := i+1; k < universeSize; k++ {
            g2 := universe[k]
            sum += length(Galaxy{X: g1.X + rowDeltas[g1.X], Y: g1.Y}, Galaxy{X: g2.X + rowDeltas[g2.X], Y: g2.Y})
        }
    }
    return sum
}

const EXPAND_DELTA = 1000000

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    universe := []Galaxy{}
    rowCoords := []int{}
    universeY := 0
    for scanner.Scan() {
        line := scanner.Text()
        if len(rowCoords) < len(line) {
            rowCoords = append(rowCoords, make([]int, len(line))...)
        }
        galaxies, isRowEmpty := parseRow(line, universeY)
        if isRowEmpty {
            universeY += EXPAND_DELTA
        } else {
            universeY += 1
        }
        for _, galaxy := range galaxies {
            rowCoords[galaxy.X] += 1
        }
        universe = append(universe, galaxies...)
    }
    rowDeltas := make([]int, len(rowCoords))
    delta := 0
    for i, count := range rowCoords {
        if count == 0 {
            delta += EXPAND_DELTA - 1
        }
        rowDeltas[i] = delta
    }

    // fmt.Printf("universe: %+v\n", universe)
    // fmt.Printf("rowDeltas: %+v\n", rowDeltas)
    
    sumOfLengths := solve(universe, rowDeltas)
    fmt.Printf("sumOfLengths: %d\n", sumOfLengths)
}
