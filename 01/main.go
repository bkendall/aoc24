package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

var number = regexp.MustCompile(`[0-9]+`)

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

	listA, listB := []int{}, []int{}
	counts := map[int]int{}
	for _, l := range strings.Split(input, "\n") {
		arr := number.FindAllStringSubmatch(l, -1)
		listA = append(listA, toInt(arr[0][0]))
		b := toInt(arr[1][0])
		listB = append(listB, b)
		counts[b] += 1
	}

	sort.Slice(listA, func(i, j int) bool { return listA[i] < listA[j] })
	sort.Slice(listB, func(i, j int) bool { return listB[i] < listB[j] })

	sum := 0
	for i := range listA {
		dist := abs(listA[i] - listB[i])
		sum += dist
	}
	fmt.Printf("distance: %d\n", sum)

	simSum := 0
	for _, v := range listA {
		val := v * counts[v]
		simSum += val
	}
	fmt.Printf("similarity: %d\n", simSum)
}
