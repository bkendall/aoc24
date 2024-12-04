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

	world := map[int]map[int]string{}
	y := 0
	x := 0
	for _, l := range strings.Split(input, "\n") {
		x = 0
		for _, c := range l {
			if _, ok := world[x]; !ok {
				world[x] = map[int]string{}
			}
			world[x][y] = string(c)
			x++
		}
		y++
	}

	maxX, maxY := x, y

	hasXMAS := func(x, y int) int {
		count := 0
		if world[x+1][y] == "M" && world[x+2][y] == "A" && world[x+3][y] == "S" {
			count++
		}
		if world[x-1][y] == "M" && world[x-2][y] == "A" && world[x-3][y] == "S" {
			count++
		}
		if world[x][y+1] == "M" && world[x][y+2] == "A" && world[x][y+3] == "S" {
			count++
		}
		if world[x][y-1] == "M" && world[x][y-2] == "A" && world[x][y-3] == "S" {
			count++
		}

		if world[x+1][y+1] == "M" && world[x+2][y+2] == "A" && world[x+3][y+3] == "S" {
			count++
		}
		if world[x-1][y+1] == "M" && world[x-2][y+2] == "A" && world[x-3][y+3] == "S" {
			count++
		}
		if world[x+1][y-1] == "M" && world[x+2][y-2] == "A" && world[x+3][y-3] == "S" {
			count++
		}
		if world[x-1][y-1] == "M" && world[x-2][y-2] == "A" && world[x-3][y-3] == "S" {
			count++
		}
		return count
	}

	count := 0
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			if world[x][y] == "X" {
				if c := hasXMAS(x, y); c > 0 {
					count += c
				}
			}
		}
	}
	fmt.Printf("count: %d\n", count)

	hasMAS := func(x, y int) int {
		count := 0
		if world[x+1][y+1] == "M" && world[x+1][y-1] == "M" && world[x-1][y+1] == "S" && world[x-1][y-1] == "S" {
			count++
		}
		if world[x+1][y+1] == "S" && world[x+1][y-1] == "M" && world[x-1][y+1] == "S" && world[x-1][y-1] == "M" {
			count++
		}
		if world[x+1][y+1] == "S" && world[x+1][y-1] == "S" && world[x-1][y+1] == "M" && world[x-1][y-1] == "M" {
			count++
		}
		if world[x+1][y+1] == "M" && world[x+1][y-1] == "S" && world[x-1][y+1] == "M" && world[x-1][y-1] == "S" {
			count++
		}
		return count
	}

	count = 0
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			if world[x][y] == "A" {
				if c := hasMAS(x, y); c > 0 {
					count += c
				}
			}
		}
	}
	fmt.Printf("count x-mas: %d\n", count)
}
