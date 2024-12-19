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

type Node struct {
	Raw       string
	Antenna   string
	Antinodes []string
}

type Point struct {
	X, Y int
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

	world := map[int]map[int]Node{}
	antennas := map[string][]Point{}
	for y, line := range strings.Split(input, "\n") {
		for x, v := range strings.Split(line, "") {
			if _, ok := world[x]; !ok {
				world[x] = map[int]Node{}
			}
			n := Node{Raw: v}
			if alphaNumeric.FindString(v) != "" {
				n.Antenna = v
				antennas[v] = append(antennas[v], Point{x, y})
			}
			world[x][y] = n
		}
	}

	maxX, maxY := len(world[0]), len(world)

	printMap(world)
	fmt.Printf("Antennas: %+v\n", antennas)

	for antenna, points := range antennas {
		for i := 0; i < len(points)-1; i++ {
			for j := i + 1; j < len(points); j++ {
				for _, ai := range antinodes(points[i], points[j], maxX, maxY) {
					if ai.X < maxX && ai.X >= 0 && ai.Y < maxY && ai.Y >= 0 {
						fmt.Printf("found antinode %+v\n", ai)
						if _, ok := world[ai.X]; !ok {
							world[ai.X] = map[int]Node{}
						}
						if _, ok := world[ai.X][ai.Y]; !ok {
							world[ai.X][ai.Y] = Node{}
						}
						n := world[ai.X][ai.Y]
						n.Antinodes = append(n.Antinodes, antenna)
						world[ai.X][ai.Y] = n
					}
				}
			}
		}
	}

	printMap(world)

	sum := 0
	for _, ys := range world {
		for _, node := range ys {
			if len(node.Antinodes) > 0 {
				sum++
			}
		}
	}
	fmt.Printf("Sum: %d\n", sum)
}

func antinodes(a, b Point, maxX, maxY int) []Point {
	diffX := b.X - a.X
	diffY := b.Y - a.Y
	arr := []Point{}
	x, y := b.X, b.Y
	for x >= 0 && y >= 0 && x < maxX && y < maxY {
		arr = append(arr, Point{X: x, Y: y})
		x += diffX
		y += diffY
	}
	x, y = a.X, a.Y
	for x >= 0 && y >= 0 && x < maxX && y < maxY {
		arr = append(arr, Point{X: x, Y: y})
		x -= diffX
		y -= diffY
	}
	return arr
}

func printMap(world map[int]map[int]Node) {
	for y := 0; y < len(world); y++ {
		for x := 0; x < len(world[0]); x++ {
			if world[x][y].Antenna != "" {
				fmt.Printf("%s", world[x][y].Antenna)
			} else if len(world[x][y].Antinodes) > 0 {
				fmt.Printf("#")
			} else {
				fmt.Printf("%s", world[x][y].Raw)
			}

		}
		fmt.Println()
	}
}
