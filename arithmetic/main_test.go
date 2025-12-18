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
	tests := []struct {
		name     string
		input    []int64
		expected int64
	}{
		{"empty slice", []int64{}, 0},
		{"single element", []int64{5}, 5},
		{"multiple elements", []int64{1, 2, 3, 4, 5}, 15},
		{"with negatives", []int64{-1, 2, -3, 4}, 2},
		{"all zeros", []int64{0, 0, 0}, 0},
		{"large numbers", []int64{1000000, 2000000, 3000000}, 6000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sum(tt.input)
			if result != tt.expected {
				t.Errorf("sum(%v) = %d; expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSquarePositive tests the squarePositive function
func TestSquarePositive(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected int64
	}{
		// Positive integers - should be squared
		{"positive 1", 1, 1},
		{"positive 2", 2, 4},
		{"positive 3", 3, 9},
		{"positive 5", 5, 25},
		{"positive 10", 10, 100},
		{"positive 100", 100, 10000},

		// CRITICAL: Negative integers - should return 0, NOT be squared
		{"negative -1", -1, 0},
		{"negative -2", -2, 0},
		{"negative -5", -5, 0},  // NOT 25!
		{"negative -10", -10, 0}, // NOT 100!
		{"negative -100", -100, 0},

		// Zero - should return 0
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := squarePositive(tt.input)
			if result != tt.expected {
				t.Errorf("squarePositive(%d) = %d; expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSumOfSquaresPositiveOnly tests the full algorithm
func TestSumOfSquaresPositiveOnly(t *testing.T) {
	tests := []struct {
		name     string
		input    []int64
		expected int64
	}{
		// Only positive numbers
		{"only positives", []int64{1, 2, 3}, 14}, // 1+4+9 = 14

		// Mixed with negatives - negatives should be excluded
		{"mixed with negatives", []int64{-1, 2, -3, 4}, 20}, // 0+4+0+16 = 20

		// All negatives - should be 0
		{"all negatives", []int64{-1, -2, -3}, 0},

		// With zero - zero is not positive
		{"with zeros", []int64{0, 1, 0, 2, 0}, 5}, // 0+1+0+4+0 = 5

		// Large numbers to test int64
		{"large numbers", []int64{1000, 2000}, 5000000}, // 1000000 + 4000000 = 5000000

		// Edge case: single positive
		{"single positive", []int64{5}, 25},

		// Edge case: single negative
		{"single negative", []int64{-5}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processed := recursive_loop(tt.input, 0, squarePositive)
			result := sum(processed)
			if result != tt.expected {
				t.Errorf("sum of squares for %v = %d; expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestStringToInt64 tests the string to int64 conversion
func TestStringToInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"positive number", "42", 42},
		{"negative number", "-42", -42},
		{"zero", "0", 0},
		{"large positive", "9223372036854775807", 9223372036854775807},
		{"large negative", "-9223372036854775808", -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string_to_int64(tt.input)
			if result != tt.expected {
				t.Errorf("string_to_int64(%s) = %d; expected %d", tt.input, result, tt.expected)
			}
		})
	}
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

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("recursive_loop[%d] = %d; expected %d", i, result[i], expected[i])
		}
	}
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
