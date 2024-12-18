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

	sum := 0
	for _, line := range strings.Split(input, "\n") {
		arr := strings.Split(line, ": ")
		total, s := toInt(arr[0]), arr[1]
		arr = strings.Split(s, " ")
		nums := make([]int, len(arr))
		for i := range arr {
			nums[i] = toInt(arr[i])
		}

		fmt.Printf("***** %d: %v\n", total, nums)

		if checkForTotal(total, nums) {
			sum += total
			continue
		}
	}

	fmt.Printf("Sum: %d\n", sum)
}

const base int = 2

func checkForTotal(total int, nums []int) bool {
	// fmt.Printf("CheckForTotal: %d, %+v\n", total, nums)
	bits := len(nums) - 1
	for x := 0; x < int(math.Pow(float64(base), float64(bits))); x++ {
		// fmt.Printf("x: %d\n", x)
		x_3 := strconv.FormatInt(int64(x), base)
		for len(x_3) < bits {
			x_3 = "0" + x_3
		}
		test := nums[0]
		// fmt.Printf("initial: %d\n", test)
		for i, v := range nums[1:] {
			bit := bits - i - 1
			// fmt.Printf("bit: %d\n", bit)
			// fmt.Printf("x: %d, x_3: %s\n", x, x_3)
			// 0 == +
			// 1 == *
			// 2 == ||
			switch toInt(string(x_3[bit])) {
			case 0:
				// fmt.Printf("add %d\n", v)
				test += v
			case 1:
				// fmt.Printf("mult\n")
				test *= v
			case 2:
				test = toInt(fmt.Sprintf("%d%d", test, v))
			}
		}
		// fmt.Printf("test: %d\n", test)
		if total == test {
			// fmt.Printf("found match!\n")
			return true
		}
	}
	return false
}
