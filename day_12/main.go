package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Order []byte

type Cache struct {
    cache map[string]bool
}

func newCache() Cache {
    return Cache{cache: map[string]bool{}}
}

func (c *Cache) add(key string, value bool) {
    c.cache[key] = value
}

func (c *Cache) get(key string) (bool, bool) {
    v, ok := c.cache[key]
    return v, ok
}

type Arrangement struct {
    Order []byte
    UnknownIndices []int
    Amounts *[]int
    SumAmounts int
    Cache *Cache
}

func getUnknownIndices(order []byte) []int {
    u := []int{}
    for i, char := range order {
        if char == '?' {
            u = append(u, i)
        }
    }
    return u
}

func readLine(line string) Arrangement {
    i := 0
    for line[i] != ' ' {
        i++
    }
    order := []byte(line[0:i])
    unknownIndices := getUnknownIndices(order)
    fmt.Printf("%+v\n", unknownIndices)
    for _, index := range unknownIndices {
        order[index] = '.'
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
    cache := newCache()
    sumAmounts := 0
    for _, a := range amounts {
        sumAmounts += a
    }
    return Arrangement{Order: order, UnknownIndices: unknownIndices, Amounts: &amounts, SumAmounts: sumAmounts, Cache: &cache}
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
    // fmt.Printf("%d/%d\n", b.next, b.target)
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
    return intsEqual(amounts, *a.Amounts)
}

func (a *Arrangement) isValid() bool {
    amount := 0
    sumAmounts := 0
    numAmounts := 0
    i := 0
    amounts := *a.Amounts
    for _, char := range a.Order {
        if char == '#' {
            amount++
        } else if amount > 0 {
            if amounts[i] != amount {
                return false
            }
            sumAmounts += amount
            if sumAmounts > a.SumAmounts {
                return false
            }
            numAmounts++
            if numAmounts > len(amounts) {
                return false
            }
            amount = 0
        }
    }
    return true
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

func buildTree(a Arrangement) int {
    if len(a.UnknownIndices) == 0 {
        if _, ok := a.Cache.get(string(a.Order)); ok {
            return 0
        }
        // fmt.Printf("%s, %+v\n", string(a.Order), *a.Amounts)
        if a.isValidOrder(a.Order) {
            return 1
        }
        return 0
    }
    u := a.UnknownIndices[0]
    unknowns := a.UnknownIndices[1:]
    order := a.Order
    // fmt.Printf("%s\n", string(order))
    order[u] = '#'
    left := buildTree(Arrangement{order, unknowns, a.Amounts, a.SumAmounts, a.Cache})
    order[u] = '.'
    right := buildTree(Arrangement{order, unknowns, a.Amounts, a.SumAmounts, a.Cache})
    return left + right
}

func (a *Arrangement) countTree() int {
    return buildTree(*a)
}

func (a *Arrangement) unfold(count int) {
    // fmt.Printf("%s\n", string(a.Order))
    order := a.Order
    amounts := *a.Amounts
    unknownIndices := a.UnknownIndices
    for count > 0 {
        orderLength := len(a.Order)
        newUnknowns := append(a.UnknownIndices, orderLength)
        for _, u := range unknownIndices {
            newUnknowns = append(newUnknowns, u + orderLength + 1)
        }
        a.UnknownIndices = newUnknowns
        a.Order = append(a.Order, '?')
        a.Order = append(a.Order, order...)
        newAmounts := append(*a.Amounts, amounts...)
        a.Amounts = &newAmounts
        count--
    }
    // a.UnknownIndices = getUnknownIndices(a.Order)
    // fmt.Printf("%s\n", string(a.Order))
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    sum := 0
    start := time.Now()
    for scanner.Scan() {
        line := scanner.Text()
        // fmt.Print(line)
        arrangement := readLine(line)
        arrangement.unfold(3)
        fmt.Printf("%+v\n",arrangement)
        sum += arrangement.countTree()
        // sum += arrangement.countPossibleArrangements()
        //fmt.Printf("%+v\n", arrangement)
    }
    end := time.Since(start)
    fmt.Printf("%d\n", sum)
    fmt.Printf("took %s\n", end)

}
