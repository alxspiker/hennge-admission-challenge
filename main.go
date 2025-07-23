package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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
	value_of_x := string_to_int32(input.Text())
	input.Scan() // Second line Yn (-100 <= Yn <= 100)

	if input.Text() == "" {
		return
	}

	// The input will be a list of integers, each separated by a newline character.
	var yn []int32 = recursive_loop(strings.Split(input.Text(), " "), 0, string_to_int32)
	length_of_yn := int32(len(yn))

	// Note: There should be no output until all the input has been received.
	if value_of_x == 0 || length_of_yn == 0 {
		return
	}

	// Note 5: It is possible that X and the number of integers Yn may not be equal.
	// If that is the case, print -1 as the output.
	if value_of_x == length_of_yn {
		// For each test case, calculate the power of four of Yn, excluding when Yn is positive
		var processed []int32 = recursive_loop(yn, 0, power_of_four)
		aggregated := sum(processed)
		min := math.Pow(-2, 31)
		max := math.Pow(2, 31)
		// Note 4: The final output is guaranteed to be within the int32 range.
		if aggregated >= int32(min) || aggregated <= int32(max) {
			fmt.Println(aggregated)
		}
	} else {
		fmt.Println("-1")
	}
}

func sum(numbers []int32) int32 {
	if len(numbers) == 0 {
		return 0
	}
	return numbers[0] + sum(numbers[1:])
}

func power_of_four(number int32) int32 {
	if number < 0 {
		return int32(math.Pow(float64(number), 4))
	}
	return 0
}

func string_to_int32(value string) int32 {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return int32(num)
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
