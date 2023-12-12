package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Arrangement struct {
    Order []byte
    UnknownIndices []int
    Amounts []int
}

func readLine(line string) Arrangement {
    i := 0
    order := []byte{}
    unknownIndices := []int{}
    for line[i] != ' ' {
        order = append(order, line[i])
        if line[i] == '?' {
            unknownIndices = append(unknownIndices, i)
        }
        i++
    }
    i++
    k := i
    amounts := []int{}
    for i <= len(line) {
        if i == len(line) || line[i] == ',' {
            amount, err := strconv.Atoi(string(line[k:i]))
            if err == nil {
                amounts = append(amounts, amount)
            }
            k = i + 1
        }
        i++
    }
    return Arrangement{Order: order, UnknownIndices: unknownIndices, Amounts: amounts}
}

type BitArrangment struct {
    length int
    target int
    current int
    next int
}

func newBitArrangment(length int) BitArrangment {
    return BitArrangment{length: length, target: 1 << length, current: 0, next: 0}
}

func (b *BitArrangment) Next() bool {
    if b.current == b.target {
        return false
    }
    b.current = b.next
    b.next += 1
    return b.current < b.target
}

func (b *BitArrangment) Get() []int {
    arrangement := make([]int, b.length)
    bitString := strconv.FormatInt(int64(b.current), 2)
    numBits := len(bitString)
    i := numBits
    k := 0
    for i > 0 {
        i--
        if bitString[i] == '1' {
            arrangement[k] = 1
        }
        k++
    }
    return arrangement
}

// ???
// ...
// #..
// ##.
// #.#
// ###
// .#.
// .##
// ..#

func intsEqual(a []int, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func (a *Arrangement) isValidOrder(order []byte) bool {
    amount := 0
    amounts := []int{}
    for _, char := range order {
        if char == '#' {
            amount++
        } else if amount > 0 {
            amounts = append(amounts, amount)
            amount = 0
        }
        // fmt.Printf("%d, %d\n", char, amount)
    }
    if amount > 0 {
        amounts = append(amounts, amount)
    }
    // fmt.Printf("%s, %+v\n", string(order), amounts)
    // fmt.Printf("%+v, %+v, %+v, %t\n", order,  amounts, a.Amounts, intsEqual(amounts, a.Amounts))
    return intsEqual(amounts, a.Amounts)
}

func (a *Arrangement) countPossibleArrangements() int {
    bitArragement := newBitArrangment(len(a.UnknownIndices))
    sum := 0
    for bitArragement.Next() {
        bits := bitArragement.Get()
        order := a.Order
        for i, bit := range bits {
            sign := byte(0)
            if bit == 1 {
                sign = '#'
            } else {
                sign = '.'
            }
            index := a.UnknownIndices[i]
            order[index] = sign
        }
        if a.isValidOrder(order) {
            sum++
        }
        // fmt.Printf("%s\n", string(order))
    }
    // fmt.Printf("%d\n", sum)
    return sum
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    sum := 0
    for scanner.Scan() {
        line := scanner.Text()
        // fmt.Print(line)
        arrangement := readLine(line)
        sum += arrangement.countPossibleArrangements()
        //fmt.Printf("%+v\n", arrangement)
    }
    fmt.Printf("%d\n", sum)

}
