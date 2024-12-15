package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

var number = regexp.MustCompile(`[0-9]+`)
var mulExpression = regexp.MustCompile(`^mul\([0-9]{1,3},[0-9]{1,3}\)`)
var doExpression = regexp.MustCompile(`^do\(\)`)
var doNotExpression = regexp.MustCompile(`^don't\(\)`)

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fatalf("%v", err)
	}
	return i
}

func abs(i int) int {
	return int(math.Abs(float64(i)))
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		fatalf("Error: %v", err)
	}
	input := string(buff)
	input = strings.TrimSpace(input)

	orders := false
	forward := map[int][]int{}
	reverse := map[int][]int{}
	pages := [][]int{}
	for _, line := range strings.Split(input, "\n") {
		if !orders && line == "" {
			orders = true
			continue
		}
		if !orders {
			arr := strings.Split(line, "|")
			one, two := toInt(arr[0]), toInt(arr[1])
			forward[one] = append(forward[one], two)
			reverse[two] = append(reverse[two], one)
			continue
		}
		arr := strings.Split(line, ",")
		nums := make([]int, 0, len(arr))
		for _, n := range arr {
			nums = append(nums, toInt(n))
		}
		pages = append(pages, nums)
	}

	fmt.Printf("forward: %+v\n", forward)
	fmt.Printf("reverse: %+v\n", reverse)
	fmt.Printf("pages: %+v\n", pages)

	sum, fixedSum := 0, 0
	for _, page := range pages {
		valid := true
		for i := 1; i < len(page); i++ {
			for j := 0; j < i; j++ {
				checking, current := page[j], page[i]
				if _, exists := forward[current]; !exists {
					if _, exists := reverse[current]; !exists {
						continue
					}
				}
				if slices.Contains(forward[current], checking) {
					fmt.Printf("%+v is invalid, since %d came before %d\n", page, checking, current)
					valid = false
				}
			}
			for j := i + 1; j < len(page); j++ {
				checking, current := page[j], page[i]
				if _, exists := forward[current]; !exists {
					if _, exists := reverse[current]; !exists {
						continue
					}
				}
				if slices.Contains(reverse[current], checking) {
					fmt.Printf("%+v is invalid, since %d came before %d\n", page, checking, current)
					valid = false
				}
			}
		}
		if valid {
			fmt.Printf("%+v is valid!\n", page)
			middle := int(math.Ceil(float64(len(page) / 2)))
			sum += page[middle]
			continue
		}
		fixed := tryAndFix(forward, reverse, page)
		middle := int(math.Ceil(float64(len(fixed) / 2)))
		fixedSum += fixed[middle]
	}

	fmt.Printf("Sum: %d\n", sum)
	fmt.Printf("Fixed Sum: %d\n", fixedSum)
}

func tryAndFix(forward, reverse map[int][]int, page []int) []int {
	changed := true
	for changed {
		changed = false
	whole:
		for i := 1; i < len(page); i++ {
			for j := 0; j < i; j++ {
				checking, current := page[j], page[i]
				if _, exists := forward[current]; !exists {
					if _, exists := reverse[current]; !exists {
						continue
					}
				}
				if slices.Contains(forward[current], checking) {
					fmt.Printf("%+v is invalid, since %d came before %d\n", page, checking, current)
					page[j], page[i] = current, checking
					changed = true
					break whole
				}
			}
			for j := i + 1; j < len(page); j++ {
				checking, current := page[j], page[i]
				if _, exists := forward[current]; !exists {
					if _, exists := reverse[current]; !exists {
						continue
					}
				}
				if slices.Contains(reverse[current], checking) {
					fmt.Printf("%+v is invalid, since %d came before %d\n", page, checking, current)
					page[j], page[i] = current, checking
					changed = true
					break whole
				}
			}
		}
	}
	fmt.Printf("%+v is valid!\n", page)
	return page
}
