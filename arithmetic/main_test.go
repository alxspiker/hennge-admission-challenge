package arithmetic

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestNoForLoops ensures no 'for' or 'goto' statements exist in the codebase
func TestNoForLoops(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "main.go", nil, 0)
	if err != nil {
		t.Fatalf("Error parsing file: %v", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.ForStmt:
			t.Errorf("FAIL: 'for' statement found at %s", fset.Position(stmt.Pos()))
		case *ast.RangeStmt:
			t.Errorf("FAIL: 'for range' statement found at %s", fset.Position(stmt.Pos()))
		case *ast.BranchStmt:
			if stmt.Tok == token.GOTO {
				t.Errorf("FAIL: 'goto' statement found at %s", fset.Position(stmt.Pos()))
			}
		}
		return true
	})
}

// TestSum tests the recursive sum function
func TestSum(t *testing.T) {
	// Test empty slice
	t.Run("empty slice", func(t *testing.T) {
		if result := sum([]int64{}); result != 0 {
			t.Errorf("sum([]) = %d; expected 0", result)
		}
	})

	// Test single element
	t.Run("single element", func(t *testing.T) {
		if result := sum([]int64{5}); result != 5 {
			t.Errorf("sum([5]) = %d; expected 5", result)
		}
	})

	// Test multiple elements
	t.Run("multiple elements", func(t *testing.T) {
		if result := sum([]int64{1, 2, 3, 4, 5}); result != 15 {
			t.Errorf("sum([1,2,3,4,5]) = %d; expected 15", result)
		}
	})

	// Test with negatives
	t.Run("with negatives", func(t *testing.T) {
		if result := sum([]int64{-1, 2, -3, 4}); result != 2 {
			t.Errorf("sum([-1,2,-3,4]) = %d; expected 2", result)
		}
	})

	// Test all zeros
	t.Run("all zeros", func(t *testing.T) {
		if result := sum([]int64{0, 0, 0}); result != 0 {
			t.Errorf("sum([0,0,0]) = %d; expected 0", result)
		}
	})

	// Test large numbers
	t.Run("large numbers", func(t *testing.T) {
		if result := sum([]int64{1000000, 2000000, 3000000}); result != 6000000 {
			t.Errorf("sum([1000000,2000000,3000000]) = %d; expected 6000000", result)
		}
	})
}

// TestSquarePositive tests the squarePositive function
func TestSquarePositive(t *testing.T) {
	// Positive integers - should be squared
	t.Run("positive 1", func(t *testing.T) {
		if result := squarePositive(1); result != 1 {
			t.Errorf("squarePositive(1) = %d; expected 1", result)
		}
	})

	t.Run("positive 2", func(t *testing.T) {
		if result := squarePositive(2); result != 4 {
			t.Errorf("squarePositive(2) = %d; expected 4", result)
		}
	})

	t.Run("positive 3", func(t *testing.T) {
		if result := squarePositive(3); result != 9 {
			t.Errorf("squarePositive(3) = %d; expected 9", result)
		}
	})

	t.Run("positive 5", func(t *testing.T) {
		if result := squarePositive(5); result != 25 {
			t.Errorf("squarePositive(5) = %d; expected 25", result)
		}
	})

	t.Run("positive 10", func(t *testing.T) {
		if result := squarePositive(10); result != 100 {
			t.Errorf("squarePositive(10) = %d; expected 100", result)
		}
	})

	t.Run("positive 100", func(t *testing.T) {
		if result := squarePositive(100); result != 10000 {
			t.Errorf("squarePositive(100) = %d; expected 10000", result)
		}
	})

	// CRITICAL: Negative integers - should return 0, NOT be squared
	t.Run("negative -1", func(t *testing.T) {
		if result := squarePositive(-1); result != 0 {
			t.Errorf("squarePositive(-1) = %d; expected 0", result)
		}
	})

	t.Run("negative -2", func(t *testing.T) {
		if result := squarePositive(-2); result != 0 {
			t.Errorf("squarePositive(-2) = %d; expected 0", result)
		}
	})

	t.Run("negative -5", func(t *testing.T) {
		// NOT 25! Negatives should be excluded
		if result := squarePositive(-5); result != 0 {
			t.Errorf("squarePositive(-5) = %d; expected 0", result)
		}
	})

	t.Run("negative -10", func(t *testing.T) {
		// NOT 100! Negatives should be excluded
		if result := squarePositive(-10); result != 0 {
			t.Errorf("squarePositive(-10) = %d; expected 0", result)
		}
	})

	t.Run("negative -100", func(t *testing.T) {
		if result := squarePositive(-100); result != 0 {
			t.Errorf("squarePositive(-100) = %d; expected 0", result)
		}
	})

	// Zero - should return 0
	t.Run("zero", func(t *testing.T) {
		if result := squarePositive(0); result != 0 {
			t.Errorf("squarePositive(0) = %d; expected 0", result)
		}
	})
}

// TestSumOfSquaresPositiveOnly tests the full algorithm
func TestSumOfSquaresPositiveOnly(t *testing.T) {
	// Only positive numbers: 1+4+9 = 14
	t.Run("only positives", func(t *testing.T) {
		processed := recursive_loop([]int64{1, 2, 3}, 0, squarePositive)
		result := sum(processed)
		if result != 14 {
			t.Errorf("sum of squares for [1,2,3] = %d; expected 14", result)
		}
	})

	// Mixed with negatives: 0+4+0+16 = 20
	t.Run("mixed with negatives", func(t *testing.T) {
		processed := recursive_loop([]int64{-1, 2, -3, 4}, 0, squarePositive)
		result := sum(processed)
		if result != 20 {
			t.Errorf("sum of squares for [-1,2,-3,4] = %d; expected 20", result)
		}
	})

	// All negatives - should be 0
	t.Run("all negatives", func(t *testing.T) {
		processed := recursive_loop([]int64{-1, -2, -3}, 0, squarePositive)
		result := sum(processed)
		if result != 0 {
			t.Errorf("sum of squares for [-1,-2,-3] = %d; expected 0", result)
		}
	})

	// With zero: 0+1+0+4+0 = 5
	t.Run("with zeros", func(t *testing.T) {
		processed := recursive_loop([]int64{0, 1, 0, 2, 0}, 0, squarePositive)
		result := sum(processed)
		if result != 5 {
			t.Errorf("sum of squares for [0,1,0,2,0] = %d; expected 5", result)
		}
	})

	// Large numbers: 1000000 + 4000000 = 5000000
	t.Run("large numbers", func(t *testing.T) {
		processed := recursive_loop([]int64{1000, 2000}, 0, squarePositive)
		result := sum(processed)
		if result != 5000000 {
			t.Errorf("sum of squares for [1000,2000] = %d; expected 5000000", result)
		}
	})

	// Single positive
	t.Run("single positive", func(t *testing.T) {
		processed := recursive_loop([]int64{5}, 0, squarePositive)
		result := sum(processed)
		if result != 25 {
			t.Errorf("sum of squares for [5] = %d; expected 25", result)
		}
	})

	// Single negative
	t.Run("single negative", func(t *testing.T) {
		processed := recursive_loop([]int64{-5}, 0, squarePositive)
		result := sum(processed)
		if result != 0 {
			t.Errorf("sum of squares for [-5] = %d; expected 0", result)
		}
	})
}

// TestStringToInt64 tests the string to int64 conversion
func TestStringToInt64(t *testing.T) {
	t.Run("positive number", func(t *testing.T) {
		if result := string_to_int64("42"); result != 42 {
			t.Errorf("string_to_int64(\"42\") = %d; expected 42", result)
		}
	})

	t.Run("negative number", func(t *testing.T) {
		if result := string_to_int64("-42"); result != -42 {
			t.Errorf("string_to_int64(\"-42\") = %d; expected -42", result)
		}
	})

	t.Run("zero", func(t *testing.T) {
		if result := string_to_int64("0"); result != 0 {
			t.Errorf("string_to_int64(\"0\") = %d; expected 0", result)
		}
	})

	t.Run("large positive", func(t *testing.T) {
		if result := string_to_int64("9223372036854775807"); result != 9223372036854775807 {
			t.Errorf("string_to_int64(\"9223372036854775807\") = %d; expected 9223372036854775807", result)
		}
	})

	t.Run("large negative", func(t *testing.T) {
		if result := string_to_int64("-9223372036854775808"); result != -9223372036854775808 {
			t.Errorf("string_to_int64(\"-9223372036854775808\") = %d; expected -9223372036854775808", result)
		}
	})
}

// compareSlicesRecursive compares two int64 slices recursively
func compareSlicesRecursive(t *testing.T, result, expected []int64, index int) {
	if index >= len(expected) {
		return
	}
	if result[index] != expected[index] {
		t.Errorf("recursive_loop[%d] = %d; expected %d", index, result[index], expected[index])
	}
	compareSlicesRecursive(t, result, expected, index+1)
}

// TestRecursiveLoop tests the generic recursive loop function
func TestRecursiveLoop(t *testing.T) {
	// Test with int64 -> int64 transformation
	input := []int64{1, 2, 3, 4, 5}
	expected := []int64{2, 4, 6, 8, 10}
	double := func(n int64) int64 { return n * 2 }

	result := recursive_loop(input, 0, double)

	if len(result) != len(expected) {
		t.Fatalf("recursive_loop returned %d elements; expected %d", len(result), len(expected))
	}

	// Use recursive comparison instead of for loop
	compareSlicesRecursive(t, result, expected, 0)
}

// TestInt64Overflow ensures no overflow with large squares
func TestInt64Overflow(t *testing.T) {
	// Maximum value that can be squared without overflow in int64
	// sqrt(9223372036854775807) ≈ 3037000499
	maxSafe := int64(3037000499)

	result := squarePositive(maxSafe)
	if result < 0 {
		t.Errorf("Overflow detected: squarePositive(%d) = %d (negative)", maxSafe, result)
	}

	// Test with value from challenge constraint (-100 to 100)
	// Maximum square would be 100^2 = 10000, well within int64
	result = squarePositive(100)
	if result != 10000 {
		t.Errorf("squarePositive(100) = %d; expected 10000", result)
	}
}
