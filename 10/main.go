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

	world := map[int]map[int]int{}
	heads := []Point{}
	for y, line := range strings.Split(input, "\n") {
		for x, c := range strings.Split(line, "") {
			if _, ok := world[x]; !ok {
				world[x] = map[int]int{}
			}
			n := toInt(c)
			world[x][y] = n
			if n == 0 {
				heads = append(heads, Point{x, y})
			}
		}
	}

	nextSteps := func(p Point) []Point {
		arr := []Point{}
		x, y := p.X, p.Y
		v := world[x][y]
		if x > 0 {
			arr = append(arr, Point{x - 1, y})
		}
		if x < len(world[0])-1 {
			arr = append(arr, Point{x + 1, y})
		}
		if y > 0 {
			arr = append(arr, Point{x, y - 1})
		}
		if y < len(world)-1 {
			arr = append(arr, Point{x, y + 1})
		}
		up := []Point{}
		for _, p := range arr {
			if world[p.X][p.Y] == v+1 {
				up = append(up, p)
			}
		}
		return up
	}

	findPeaks = func(p Point) []Point {
		v := world[p.X][p.Y]
		if v == 9 {
			return []Point{p}
		}
		arr := []Point{}
		// fmt.Printf("at %d, %d: %d\n", p.X, p.Y, v)
		for _, pt := range nextSteps(p) {
			// fmt.Printf("next: %d, %d (%d)\n", pt.X, pt.Y, world[pt.X][pt.Y])
			arr = append(arr, findPeaks(pt)...)
		}
		return arr
	}

	sum := 0
	unique := 0
	for _, h := range heads {
		// fmt.Printf("starting at %d, %d: ", h.X, h.Y)
		peaks := findPeaks(h)
		unique += len(peaks)
		found := map[int]map[int]bool{}
		count := 0
		for _, p := range peaks {
			if _, ok := found[p.X]; !ok {
				found[p.X] = map[int]bool{}
			}
			if found[p.X][p.Y] {
				continue
			}
			found[p.X][p.Y] = true
			count++
		}
		// fmt.Printf("%d\n", count)
		sum += count
	}
	fmt.Printf("Sum: %d\n", sum)
	fmt.Printf("Unique paths: %d\n", unique)
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
