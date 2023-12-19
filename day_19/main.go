package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Part struct {
    A int
    M int
    S int
    X int
}

type Range struct  {
    Start int
    End int
}

type RangedPart struct {
    A Range
    M Range
    S Range
    X Range
}

func (r *RangedPart) getRange(name string) Range {
    switch name {
    case "a":
        return r.A
    case "m":
        return r.M
    case "s":
        return r.S
    case "x":
        return r.X
    }
    panic(fmt.Errorf("unknown name '%s'!", name))
}

func (r *RangedPart) setRange(name string, rng Range) {
    switch name {
    case "a":
        r.A = rng
    case "m":
        r.M = rng
    case "s":
        r.S = rng
    case "x":
        r.X = rng
    }
}

type Condition struct {
    Part string
    Operator string
    Value int
    Result string
    IsResultOnly bool
}

func parseConditions(line string) (string, []Condition) {
    conditions := []Condition{}
    parts := strings.Split(line, "{")
    key := parts[0]
    conditionsString := parts[1][:len(parts[1])-1]
    conditionsStrings := strings.Split(conditionsString, ",")
    // fmt.Printf("%s, %+v\n", key, conditionsStrings)
    for _, s := range conditionsStrings {
        parts := strings.Split(s, ":")
        if len(parts) == 1 {
            conditions = append(conditions, Condition{Result: parts[0], IsResultOnly: true})
        } else {
            result := parts[1]
            operator := "<"
            if strings.Contains(parts[0], ">") {
                operator = ">"
            }
            parts = strings.Split(parts[0], operator)
            part := parts[0]
            valueString := parts[1]
            value, err := strconv.Atoi(valueString);
            if err != nil {
                panic(err)
            }
            conditions = append(conditions, Condition{part, operator, value, result, false})
        }
    }
    return key, conditions
}

func parsePart(line string) Part {
    line = line[1:len(line) - 1]
    components := strings.Split(line, ",")
    part := Part{}
    for _, component := range components {
        name := component[:1]
        valueString := component[2:]
        value, err := strconv.Atoi(valueString)
        if err != nil {
            panic(err)
        }
        switch name {
        case "m":
            part.M = value
        case "a":
            part.A = value
        case "s":
            part.S = value
        case "x":
            part.X = value
        }
    }
    return part
}

func gt(partValue int,value int) bool {
    return partValue > value
}

func lt(partValue int, value int) bool {
    return partValue < value
}

var predicates = map[string]func(a int, b int) bool{
    ">": gt,
    "<": lt,
}

func (p *Part) getValue(name string) int {
    switch name {
    case "a":
        return p.A
    case "m":
        return p.M
    case "s":
        return p.S
    case "x":
        return p.X
    }
    return 0
}

func (p *Part) isAccpeted(conditions map[string][]Condition, key string) bool {
    conditionList, ok := conditions[key]
    if !ok {
        return false
    }
    for _, condition := range conditionList {
        // fmt.Printf("%+v\n", condition)
        if condition.IsResultOnly {
            if condition.Result == "A" {
                return true
            }
            if condition.Result == "R" {
                return false
            }
            return p.isAccpeted(conditions, condition.Result)
        }
        predicate, ok := predicates[condition.Operator]
        if !ok {
            break
        }
        if predicate(p.getValue(condition.Part), condition.Value) {
            if condition.Result == "A" {
                return true
            }
            if condition.Result == "R" {
                return false
            }
            return p.isAccpeted(conditions, condition.Result)
        }
    }
    return false
}

func acceptedRanges(conditions map[string][]Condition, rangedPart RangedPart, key string) []RangedPart {
    ranges := []RangedPart{}
    cs, ok := conditions[key]
    if !ok {
        return []RangedPart{}
    }
    for _, c := range cs {
        if c.IsResultOnly {
            if c.Result == "A" {
                ranges = append(ranges, rangedPart)
                break
            }
            if c.Result == "R" {
                break
            }
            ranges = append(ranges, acceptedRanges(conditions, rangedPart, c.Result)...)
            break
        }
        rp := rangedPart
        r := rp.getRange(c.Part)
        operator := c.Operator
        if c.Result == "R" {
            if operator == "<" {
                operator = ">"
            } else {
                operator = "<"
            }
        }
        if operator == "<" {
            if r.Start < c.Value && r.End >= c.Value {
                r.End = c.Value - 1
            } else {
                break
            }
        }
        if operator == ">" {
            if r.End > c.Value && r.Start <= c.Value {
                r.Start = c.Value + 1
            } else {
                break
            }
        }
        rp.setRange(c.Part, r)
        if c.Result == "A" || c.Result == "R" {
            return []RangedPart{rp}
        } else {
            return []RangedPart{}
        }
    }
    return []RangedPart{}
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    conditions := map[string][]Condition{}

    for scanner.Scan() {
        line := scanner.Text()

        if line == "" {
            break
        }
        key, conditionList := parseConditions(line)

        conditions[key] = conditionList
    }
    // fmt.Printf("%+v\n", conditions)
    parts := []Part{}

    for scanner.Scan() {
        line := scanner.Text()

        part := parsePart(line)
        parts = append(parts, part)
    }

    // fmt.Printf("%+v\n", parts)

    accepted := []Part{}

    for _, part := range parts {
        if part.isAccpeted(conditions, "in") {
            accepted = append(accepted, part)
        }
    }
    // fmt.Printf("%+v\n", accepted)
    sum := 0
    for _, part := range accepted {
        sum += part.A + part.M + part.S + part.X
    }
    fmt.Printf("%d\n", sum)

    /* rangedPart := RangedPart{
        Range{0, 4000},
        Range{0, 4000},
        Range{0, 4000},
        Range{0, 4000},
    } */

    // acceptedRanges(conditions, rangedPart, "in")
}
