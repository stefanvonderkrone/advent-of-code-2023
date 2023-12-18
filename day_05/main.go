package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
    Dest int
    Source int
    Length int
}

type Category struct {
    Ranges []Range
    From string
    To string
}

type Garden struct {
    Seeds []int
    Relations map[string]Category
}

func (g *Garden) find(key string, value int) int {
    cat, ok := g.Relations[key]
    if !ok {
        return value
    }
    // fmt.Printf("key: %s, value: %d, cat: %+v\n", key, value, cat)
    for _, r := range cat.Ranges {
        if value >= r.Source && value < r.Source + r.Length {
            a := value - r.Source
            return g.find(cat.To, r.Dest + a)
        }
    }
    return g.find(cat.To, value)
}

func parseSeeds(line string) []int {
    parts := strings.Split(line, " ")[1:]
    seeds := make([]int, len(parts))
    for i, s := range parts {
        seed, err := strconv.Atoi(s)
        if err != nil {
            panic(err)
        }
        seeds[i] = seed
    }
    return seeds
}

func parseFromTo(line string) (string, string) {
    relation := strings.Split(line, " ")[0]
    // fmt.Printf("%s\n", relation)
    parts := strings.Split(relation, "-")
    from := parts[0]
    to := parts[2]
    return from, to
}

func parseRange(line string) Range {
    parts := strings.Split(line, " ")
    dest, err := strconv.Atoi(parts[0])
    if err != nil {
        panic(err)
    }
    source, err := strconv.Atoi(parts[1])
    if err != nil {
        panic(err)
    }
    r, err := strconv.Atoi(parts[2])
    if err != nil {
        panic(err)
    }
    return Range{dest, source, r}
}

func readCategory(scanner *bufio.Scanner) Category {
    cat := Category{Ranges: []Range{}}
    j := 0
    line := scanner.Text()
    for {
        if line == "" {
            break
        }
        // fmt.Printf("%s\n", line)
        if j == 0 {
            j++
            from, to := parseFromTo(line)
            cat.From = from
            cat.To = to
            // fmt.Printf("got relation: %s, %s\n", from, to)
        } else {
            cat.Ranges = append(cat.Ranges, parseRange(line))
        }
        if !scanner.Scan() {
            break
        }
        line = scanner.Text()
    }
    return cat
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    i := 0
    garden := Garden{Relations: map[string]Category{}}
    currentCategory := Category{}
    for scanner.Scan() {
        line := scanner.Text()
        if i == 0 {
            i++
            garden.Seeds = parseSeeds(line)
            // fmt.Printf("got seeds: %+v\n", garden.Seeds)
            continue
        }
        if line == "" {
            continue
        }
        currentCategory = readCategory(scanner)
        garden.Relations[currentCategory.From] = currentCategory
    }
    // for key, cat := range garden.Relations {
    //     fmt.Printf("key: %s, category: %+v\n", key, cat)
    // }
    // fmt.Printf("%+v\n", garden)
    // for _, seed := range garden.Seeds {
    //     fmt.Printf("%d, %d\n", seed, garden.find("seed", seed))
    // }
    minLoc := -1
    for _, seed := range garden.Seeds {
        loc := garden.find("seed", seed)
        if minLoc < 0 || loc < minLoc {
            minLoc = loc
        }
    }
    fmt.Printf("location: %d\n", minLoc)
}

