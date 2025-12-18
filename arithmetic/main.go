package arithmetic

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Note 3: Take input from standard input, and output to standard output.
	input := bufio.NewScanner(os.Stdin)
	input.Scan() // First line N (1 <= N <= 100)

	if input.Text() == "" {
		return
	}

	// The first line of the input will indicate the number of test cases to follow.
	value_of_x := string_to_int64(input.Text())
	input.Scan() // Second line Yn (-100 <= Yn <= 100)

	if input.Text() == "" {
		return
	}

	// The input will be a list of integers, each separated by a space character.
	var yn []int64 = recursive_loop(strings.Split(input.Text(), " "), 0, string_to_int64)
	length_of_yn := int64(len(yn))

	// Note: There should be no output until all the input has been received.
	if value_of_x == 0 || length_of_yn == 0 {
		return
	}

	// Note 5: It is possible that X and the number of integers Yn may not be equal.
	// If that is the case, print -1 as the output.
	if value_of_x == length_of_yn {
		// Calculate the sum of squares for positive integers only (>0)
		// Negatives and zero are excluded before calculation
		var processed []int64 = recursive_loop(yn, 0, squarePositive)
		aggregated := sum(processed)
		fmt.Println(aggregated)
	} else {
		fmt.Println("-1")
	}
}

// sum calculates the sum of numbers using recursion (no for loops)
// Uses int64 to prevent overflow on large squares
func sum(numbers []int64) int64 {
	if len(numbers) == 0 {
		return 0
	}
	return numbers[0] + sum(numbers[1:])
}

// squarePositive returns n^2 for positive integers (n > 0), 0 otherwise
// This ensures negative numbers are excluded BEFORE any calculation
func squarePositive(number int64) int64 {
	if number > 0 {
		return number * number
	}
	return 0
}

func string_to_int64(value string) int64 {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return num
}

func recursive_loop[T, U any](list []T, index int, callback func(T) U) []U {
	result := make([]U, len(list))
	transformer(list, result, index, callback)
	return result
}

func transformer[T, U any](input []T, output []U, index int, callback func(T) U) {
	if index+1 > len(input) {
		return
	}
	output[index] = callback(input[index])
	transformer(input, output, index+1, callback)
}
