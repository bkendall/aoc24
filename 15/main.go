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

type Entity int

const (
	Empty Entity = iota
	Box
	Wall
	Robot
)

func (e Entity) String() string {
	switch e {
	case Box:
		return "O"
	case Wall:
		return "#"
	case Robot:
		return "@"
	default:
		return "."
	}
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
	world := map[int]map[int]Entity{}
	instructions := false
	y := 0
	robot := Point{}
	inst := ""
	for _, line := range strings.Split(input, "\n") {
		if instructions {
			inst += line
			continue
		}
		if line == "" {
			instructions = true
			continue
		}
		for x, c := range strings.Split(line, "") {
			if _, ok := world[x]; !ok {
				world[x] = map[int]Entity{}
			}
			switch c {
			case "#":
				world[x][y] = Wall
			case ".":
				world[x][y] = Empty
			case "O":
				world[x][y] = Box
			case "@":
				world[x][y] = Robot
				robot.X, robot.Y = x, y
			}
		}
		y++
	}

	printWorld(world)

	for _, ins := range strings.Split(inst, "") {
		dx, dy := 0, 0
		fmt.Printf("ins: %s\n", ins)
		switch ins {
		case "<":
			dx = -1
		case "^":
			dy = -1
		case ">":
			dx = +1
		case "v":
			dy = +1
		default:
			fatalf("unknown instruction: %q\n", ins)
		}

		i := 1
		done := false
		wall := false
		for !done {
			px, py := robot.X+(dx*i), robot.Y+(dy*i)
			if world[px][py] == Box {
				i++
			} else if world[px][py] == Empty {
				// We move *something* here. Stop.
				break
			} else if world[px][py] == Wall {
				// If we get here, we cannot move anything.
				wall = true
				break
			}
		}

		if wall {
			continue
		}
		for ; i > 0; i-- {
			px, py := robot.X+(dx*i), robot.Y+(dy*i)
			sx, sy := robot.X+(dx*(i-1)), robot.Y+(dy*(i-1))
			fmt.Printf("Swapping %d, %d and %d, %d\n", px, py, sx, sy)
			world[px][py] = world[sx][sy]
		}
		// We moved, most definitely. Record where the robot is.
		px, py := robot.X+(dx), robot.Y+(dy)
		world[robot.X][robot.Y] = Empty
		robot.X, robot.Y = px, py
	}

	printWorld(world)

	fmt.Printf("Sum: %d\n", scoreWorld(world))
}

func scoreWorld(w map[int]map[int]Entity) int {
	sum := 0
	count := 0
	for y := 0; y < len(w[0]); y++ {
		for x := 0; x < len(w); x++ {
			if w[x][y] == Box {
				count++
				sum += 100*y + x
			}
		}
	}
	fmt.Printf("Count: %d\n", count)
	return sum
}

func printWorld(w map[int]map[int]Entity) {
	for y := 0; y < len(w[0]); y++ {
		for x := 0; x < len(w); x++ {
			fmt.Printf("%s", w[x][y])
		}
		fmt.Println()
	}
	fmt.Println()
}
