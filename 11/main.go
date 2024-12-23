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

var alphaNumeric = regexp.MustCompile(`[0-9a-zA-Z]`)
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

type Point struct {
	X, Y int
}

var findPaths func(p Point) int
var findPeaks func(p Point) []Point

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		fatalf("Error: %v", err)
	}
	input := string(buff)
	input = strings.TrimSpace(input)

	world := map[int]int{}
	for _, line := range strings.Split(input, "\n") {
		for _, c := range strings.Split(line, " ") {
			n := toInt(c)
			world[n] += 1
		}
	}

	fmt.Printf("Initial: %+v\n", world)
	steps := 75
	for step := range steps {
		if step%5 == 0 {
			fmt.Printf("Step %d...\n", step)
		}
		newWorld := map[int]int{}
		for v, c := range world {
			if v == 0 {
				newWorld[1] += c
				continue
			}
			if s := fmt.Sprintf("%d", v); len(s)%2 == 0 {
				h := len(s) / 2
				newWorld[toInt(s[:h])] += c
				newWorld[toInt(s[h:])] += c
				continue
			}
			newWorld[v*2024] += c
		}
		world = newWorld
		// fmt.Printf("After step %d: %+v\n", step, newWorld)
	}

	sum := 0
	for _, c := range world {
		sum += c
	}
	fmt.Printf("Stones: %d\n", sum)
}

func printWorld(w map[int]int, max int) {
	for l := 0; l < max; l++ {
		v, ok := w[l]
		if ok {
			fmt.Printf("%d", v)
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()
}
