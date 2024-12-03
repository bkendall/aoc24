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

	stream := ""
	for _, l := range strings.Split(input, "\n") {
		stream += l
	}

	do := true
	sum := 0
	i := 0
	for i < len(stream) {
		// fmt.Printf("i: %d, s: %q\n", i, stream[i:])
		if e := doExpression.FindStringIndex(stream[i:]); e != nil {
			// fmt.Printf("found: %q\n", stream[i:i+e[1]])
			do = true
			i += e[1]
			continue
		}
		if e := doNotExpression.FindStringIndex(stream[i:]); e != nil {
			// fmt.Printf("found: %q\n", stream[i:i+e[1]])
			do = false
			i += e[1]
			continue
		}
		if e := mulExpression.FindStringSubmatch(stream[i:]); e != nil {
			if !do {
				i += len(e[0])
				continue
			}
			nums := number.FindAllStringSubmatch(e[0], -1)
			one, two := toInt(nums[0][0]), toInt(nums[1][0])
			sum += one * two
			i += len(e[0])
			continue
		}
		i++
	}
	fmt.Printf("sum: %d\n", sum)
}
