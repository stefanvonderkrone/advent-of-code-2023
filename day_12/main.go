package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readLine(line string) (string, []int) {
    i := 0
    for line[i] != ' ' {
        i++
    }
    order := line[0:i]
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
    return order, amounts
}

var cache = map[string]int{}

func count(cfg string, nums []int) int {
    if cfg == "" {
        if len(nums) == 0 {
            return 1
        }
        return 0
    }

    if len(nums) == 0 {
        if strings.Contains(cfg, "#") {
            return 0
        }
        return 1
    }

    key := strings.Join([]string{cfg, fmt.Sprintf("%+v", nums)}, " ")

    if r, ok := cache[key]; ok {
        return r
    }

    result := 0

    if cfg[0] == '.' || cfg[0] == '?' {
        result += count(cfg[1:], nums)
    }

    if cfg[0] == '#' || cfg[0] == '?' {
        if nums[0] > len(cfg) {
            return result
        }
        if strings.Contains(cfg[:nums[0]], ".") {
            return result
        }

        if nums[0] == len(cfg) || cfg[nums[0]] != '#' {
            startIndex := nums[0] + 1
            if startIndex > len(cfg) {
                startIndex = len(cfg)
            }
            result += count(cfg[startIndex:], nums[1:])
        }
    }

    cache[key] = result

    return result
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    sum := 0
    start := time.Now()
    for scanner.Scan() {
        line := scanner.Text()
        // fmt.Print(line)
        cfg, nums := readLine(line)
        // fmt.Printf("%+v\n", nums)
        cfg = strings.Join([]string{cfg, cfg, cfg, cfg, cfg}, "?")
        tmpNums := nums
        nums = append(nums, tmpNums...)
        nums = append(nums, tmpNums...)
        nums = append(nums, tmpNums...)
        nums = append(nums, tmpNums...)
        // fmt.Printf("%+v\n", nums)
        // fmt.Print("-----------\n")
        sum += count(cfg, nums)
        // sum += arrangement.countPossibleArrangements()
        //fmt.Printf("%+v\n", arrangement)
    }
    end := time.Since(start)
    fmt.Printf("%d\n", sum)
    fmt.Printf("took %s\n", end)

}
