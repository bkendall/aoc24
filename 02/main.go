package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
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

	levels := [][]int{}
	for _, l := range strings.Split(input, "\n") {
		arr := number.FindAllStringSubmatch(l, -1)
		level := []int{}
		for _, v := range arr {
			n := toInt(v[0])
			level = append(level, n)
		}
		levels = append(levels, level)
	}

	safeLevels := 0
	safeLevelsWithDamper := 0
	for _, level := range levels {
		if isSafe(level) {
			safeLevels++
			safeLevelsWithDamper++
		} else {
			if isSafeWithDamper(level) {
				safeLevelsWithDamper++
			}
		}
	}
	fmt.Printf("safe levels: %d\n", safeLevels)
	fmt.Printf("safe levels with damper: %d\n", safeLevelsWithDamper)
}

type Direction int

const (
	Unknown Direction = iota
	Up
	Down
)

func isSafeWithDamper(l []int) bool {
	for skip := 0; skip < len(l); skip++ {
		arr := []int{}
		arr = append(arr, l[:skip]...)
		arr = append(arr, l[skip+1:]...)
		if isSafe(arr) {
			return true
		}
	}
	return false
}

func isSafe(l []int) bool {
	var d Direction
	if l[0] < l[1] {
		d = Up
	} else if l[0] > l[1] {
		d = Down
	} else {
		return false
	}

	for i := 0; i < len(l)-1; i++ {
		if d == Up && !(l[i] < l[i+1]) {
			return false
		}
		if d == Down && !(l[i] > l[i+1]) {
			return false
		}
		diff := abs(l[i] - l[i+1])
		if diff < 1 || diff > 3 {
			return false
		}
	}
	return true
}
