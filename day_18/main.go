package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Point struct {
    X int
    Y int
}

var directions = map[string]Point{
    "U": Point{0, -1},
    "D": Point{0, 1},
    "L": Point{-1, 0},
    "R": Point{1, 0},
}

const directionIndices = "RDLU"

func scanPoints() ([]Point, int) {
    scanner := bufio.NewScanner(os.Stdin)

    points := []Point{Point{}}
    boundary := 0
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Split(line, " ")
        direction := parts[0]
        stepsString := parts[1]
        // pt2
        x := parts[2][2:]
        stepsHex := x[0:5]
        dir := x[5:6]
        s, err := strconv.ParseInt(stepsHex, 16, 64)
        if err != nil {
            panic(err)
        }
        dirIndex, err := strconv.Atoi(dir)
        if err != nil {
            panic(err)
        }
        steps, err := strconv.Atoi(stepsString)
        if err != nil {
            panic(err)
        }
        steps = int(s)
        boundary += steps
        directionPoint := directions[direction]
        directionPoint = directions[directionIndices[dirIndex:dirIndex + 1]]
        lastPoint := points[len(points) - 1]
        point := Point{lastPoint.X + directionPoint.X * steps, lastPoint.Y + directionPoint.Y * steps}
        points = append(points, point)
    }
    return points, boundary
}

func main() {

    points, boundary := scanPoints()

    area := 0
    for i, point := range points {
        lastPoint := point
        if i == 0 {
            lastPoint = points[len(points) - 1]
        } else {
            lastPoint = points[i - 1]
        }
        nextPoint := point
        if i == len(points) - 1 {
            nextPoint = points[0]
        } else {
            nextPoint = points[i + 1]
        }
        area += point.Y * (lastPoint.X - nextPoint.X)
    }
    area = int(math.Abs(float64(area / 2)))
    innerArea := area - boundary / 2 + 1
    // fmt.Printf("%+v\n", points)
    fmt.Printf("area: %d, boundary: %d, inner area + boundary: %d\n", area, boundary, innerArea + boundary)
}

