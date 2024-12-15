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
	gX, gY := -1, -1
	g := "^"
	for y, line := range strings.Split(input, "\n") {
		for x, c := range strings.Split(line, "") {
			if _, ok := world[x]; !ok {
				world[x] = make(map[int]string)
			}
			world[x][y] = c
			if c == "^" {
				gX, gY = x, y
			}
		}
	}

	yM := len(world)
	visited := make([]bool, yM*len(world[0]))
	visited[gY*yM+gX] = true

loop:
	for {
		// printMap(world)
		switch g {
		case "^":
			if gY == 0 {
				break loop
			}
			if world[gX][gY-1] == "#" {
				g = ">"
				world[gX][gY] = g
			} else {
				world[gX][gY], world[gX][gY-1] = ".", g
				gY--
				visited[gY*yM+gX] = true
			}
		case ">":
			if gX == len(world)-1 {
				break loop
			}
			if world[gX+1][gY] == "#" {
				g = "v"
				world[gX][gY] = g
			} else {
				world[gX][gY], world[gX+1][gY] = ".", g
				gX++
				visited[gY*yM+gX] = true
			}
		case "v":
			if gY == len(world[0])-1 {
				break loop
			}
			if world[gX][gY+1] == "#" {
				g = "<"
				world[gX][gY] = g
			} else {
				world[gX][gY], world[gX][gY+1] = ".", g
				gY++
				visited[gY*yM+gX] = true
			}
		case "<":
			if gX == 0 {
				break loop
			}
			if world[gX-1][gY] == "#" {
				g = "^"
				world[gX][gY] = g
			} else {
				world[gX][gY], world[gX-1][gY] = ".", g
				gX--
				visited[gY*yM+gX] = true
			}
		default:
			fatalf("Unknown direction: %q\n", g)
		}
	}

	positions := 0
	for _, b := range visited {
		if b {
			positions++
		}
	}
	fmt.Printf("Positions: %d\n", positions)
}

func printMap(world map[int]map[int]string) {
	for y := 0; y < len(world); y++ {
		for x := 0; x < len(world[0]); x++ {
			fmt.Printf("%s", world[x][y])
		}
		fmt.Println()
	}
}
