package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
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
	fileStarts := map[int]int{}
	fileSizes := map[int]int{}
	index := 0
	pointer := 0
	for _, line := range strings.Split(input, "\n") {
		file := true
		for _, c := range strings.Split(line, "") {
			n := toInt(c)
			if file {
				fileStarts[index] = pointer
				fileSizes[index] = n
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

	max := pointer
	lastFile := index - 1

	findFree := func() map[int]int {
		m := map[int]int{}
		i := 0
		free := false
		start := 0
		space := 0
		for ; i < max; i++ {
			if _, ok := world[i]; ok {
				if free {
					m[start] = space
					start = -1
					space = 0
				}
				free = false
			} else {
				if !free {
					start = i
				}
				free = true
				space++
			}
		}
		return m
	}

	// fmt.Printf("%+v\n", world)
	// fmt.Printf("Free: %+v\n", findFree())

	for index = lastFile; index >= 0; index-- {
		fileLoc := fileStarts[index]
		// Count the file
		size := fileSizes[index]
		// fmt.Printf("looking at file %d, size %d\n", index, size)

		freeSpace := findFree()
		locs := make([]int, 0, len(freeSpace))
		for k := range freeSpace {
			locs = append(locs, k)
		}
		slices.Sort(locs)
		// fmt.Printf("available space indexes: %+v\n", locs)

		for _, l := range locs {
			available := freeSpace[l]
			if l >= fileLoc {
				break
			}
			if size <= available {
				for i := 0; i < size; i++ {
					world[l+i] = index
					delete(world, fileLoc+i)
				}
				break
			}
		}

		// printWorld(world, max)
	}

	// fmt.Printf("%+v\n", world)

	sum := 0
	for i, v := range world {
		sum += i * v
	}
	fmt.Printf("Sum: %d\n", sum)
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
