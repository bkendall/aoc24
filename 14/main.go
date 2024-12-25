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
var coordinate = regexp.MustCompile(`(\-?[0-9]+),(\-?[0-9]+)`)

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

type Robot struct {
	Velocity Point
	Position Point
}

func (r *Robot) Move() {
	r.Position.X += r.Velocity.X
	for r.Position.X < 0 {
		r.Position.X += maxX
	}
	for r.Position.X >= maxX {
		r.Position.X -= maxX
	}

	r.Position.Y += r.Velocity.Y
	for r.Position.Y < 0 {
		r.Position.Y += maxY
	}
	for r.Position.Y >= maxY {
		r.Position.Y -= maxY
	}
}

// const maxX, maxY = 11, 7

const maxX, maxY = 101, 103

func robotsToWorld(robots []*Robot) map[int]map[int][]*Robot {
	m := map[int]map[int][]*Robot{}
	for _, r := range robots {
		if _, ok := m[r.Position.X]; !ok {
			m[r.Position.X] = map[int][]*Robot{}
		}
		m[r.Position.X][r.Position.Y] = append(m[r.Position.X][r.Position.Y], r)
	}
	return m
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

	// p=0,4 v=3,-3
	robots := []*Robot{}
	for _, line := range strings.Split(input, "\n") {
		arr := strings.Split(line, " ")
		p := strings.Split(arr[0], "=")[1]
		v := strings.Split(arr[1], "=")[1]
		matchP := coordinate.FindStringSubmatch(p)
		matchV := coordinate.FindStringSubmatch(v)
		r := &Robot{
			Position: Point{X: toInt(matchP[1]), Y: toInt(matchP[2])},
			Velocity: Point{X: toInt(matchV[1]), Y: toInt(matchV[2])},
		}
		robots = append(robots, r)
	}

	qWidth := maxX / 2
	qHeight := maxY / 2

	for i := 0; i < 100_000; i++ {
		for _, r := range robots {
			r.Move()
		}
		if allUnique(robotsToWorld(robots)) {
			fmt.Printf("%d:\n", i+1)
			printWorld(robotsToWorld(robots))
			break
		}
	}

	counts := quadCounts(robots, qWidth, qHeight)
	fmt.Printf("Counts: %+v\n", counts)
	factor := 1
	for _, c := range counts {
		factor *= c
	}

	fmt.Printf("Factor: %d\n", factor)
}

func quadCounts(r []*Robot, qWidth, qHeight int) []int {
	w := robotsToWorld(r)
	arr := []int{}
	sum := 0
	for x := 0; x < qWidth; x++ {
		if _, ok := w[x]; !ok {
			continue
		}
		for y := 0; y < qHeight; y++ {
			sum += len(w[x][y])
		}
	}
	arr = append(arr, sum)
	sum = 0
	for x := qWidth + 1; x <= maxX; x++ {
		if _, ok := w[x]; !ok {
			continue
		}
		for y := 0; y < qHeight; y++ {
			sum += len(w[x][y])
		}
	}
	arr = append(arr, sum)
	sum = 0
	for x := 0; x < qWidth; x++ {
		if _, ok := w[x]; !ok {
			continue
		}
		for y := qHeight + 1; y <= maxY; y++ {
			sum += len(w[x][y])
		}
	}
	arr = append(arr, sum)
	sum = 0
	for x := qWidth + 1; x <= maxX; x++ {
		if _, ok := w[x]; !ok {
			continue
		}
		for y := qHeight + 1; y <= maxY; y++ {
			sum += len(w[x][y])
		}
	}
	arr = append(arr, sum)
	return arr
}

func allUnique(w map[int]map[int][]*Robot) bool {
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			row, ok := w[x]
			if !ok {
				continue
			}
			if len(row[y]) > 1 {
				return false
			}
		}
	}
	return true
}

func printWorld(w map[int]map[int][]*Robot) {
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			row, ok := w[x]
			if !ok {
				fmt.Printf(".")
				continue
			}
			if len(row[y]) == 0 {
				fmt.Printf(".")
				continue
			}
			fmt.Printf("%d", len(row[y]))
		}
		fmt.Println()
	}
	fmt.Println()
}
