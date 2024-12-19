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
	index := 0
	pointer := 0
	for _, line := range strings.Split(input, "\n") {
		file := true
		for _, c := range strings.Split(line, "") {
			n := toInt(c)
			if file {
				end := pointer + n
				for pointer < end {
					world[pointer] = index
					pointer++
				}
				index++
			} else {
				pointer += n
			}
			file = !file
		}
	}

	// fmt.Printf("%+v\n", world)

	for {
		if _, ok := world[pointer]; !ok {
			pointer--
		} else {
			break
		}
	}

	head := 0
	for {
		if _, ok := world[head]; ok {
			head++
		} else {
			break
		}
	}

	for pointer > head {
		world[head] = world[pointer]
		delete(world, pointer)
		for {
			if _, ok := world[pointer]; !ok {
				pointer--
			} else {
				break
			}
		}
		for {
			if _, ok := world[head]; ok {
				head++
			} else {
				break
			}
		}
	}

	// fmt.Printf("%+v\n", world)

	sum := 0
	for i, v := range world {
		sum += i * v
	}
	fmt.Printf("Sum: %d\n", sum)
}
