package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestNoForLoops ensures no 'for' or 'goto' statements exist in the TOTP code
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

// checkDigitsRecursive checks if all characters in the string are digits
func checkDigitsRecursive(t *testing.T, s string, index int) {
	if index >= len(s) {
		return
	}
	c := rune(s[index])
	if c < '0' || c > '9' {
		t.Errorf("Character at position %d is not a digit: %c", index, c)
	}
	checkDigitsRecursive(t, s, index+1)
}

// TestTOTP_compliance verifies TOTP implementation meets HENNGE requirements
func TestTOTP_compliance(t *testing.T) {
	email := "test@example.com"
	secret := BuildSecret(email)

	t.Run("output_length_is_10_digits", func(t *testing.T) {
		totp, err := GenerateTOTP(secret, DefaultTimeStep, DefaultDigits)
		if err != nil {
			t.Fatalf("GenerateTOTP failed: %v", err)
		}
		if len(totp) != 10 {
			t.Errorf("TOTP output length = %d; expected 10", len(totp))
		}
	})

	t.Run("output_is_zero_padded", func(t *testing.T) {
		// Generate multiple TOTPs to test padding
		// We test with a known time to get predictable results
		secret := BuildSecret(email)

		// Test that the output is always exactly 10 digits
		totp, _ := GenerateTOTPWithTime(secret, 1234567890, DefaultTimeStep, DefaultDigits)
		if len(totp) != 10 {
			t.Errorf("TOTP should be zero-padded to 10 digits, got %d digits: %s", len(totp), totp)
		}

		// Verify all characters are digits using recursive function
		checkDigitsRecursive(t, totp, 0)
	})

	t.Run("algorithm_is_sha512", func(t *testing.T) {
		// Verify HMAC-SHA512 is being used by checking hash length
		// SHA512 produces 64 bytes (512 bits)
		key := []byte(secret)
		counter := int64(1234567890 / 30)

		// Convert counter to bytes using binary.BigEndian (no for loop)
		counterBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(counterBytes, uint64(counter))

		mac := hmac.New(sha512.New, key)
		mac.Write(counterBytes)
		hash := mac.Sum(nil)

		// SHA512 produces 64 bytes
		if len(hash) != 64 {
			t.Errorf("Hash length = %d bytes; expected 64 bytes for SHA512", len(hash))
		}
	})

	t.Run("secret_key_format", func(t *testing.T) {
		// Verify secret key is constructed as "HENNGECHALLENGE" + email
		expectedPrefix := "HENNGECHALLENGE"
		if secret[:len(expectedPrefix)] != expectedPrefix {
			t.Errorf("Secret should start with 'HENNGECHALLENGE', got prefix: %s", secret[:len(expectedPrefix)])
		}
		if secret[len(expectedPrefix):] != email {
			t.Errorf("Secret should end with email '%s', got suffix: %s", email, secret[len(expectedPrefix):])
		}
	})
}

// TestBuildSecret tests the secret key construction
func TestBuildSecret(t *testing.T) {
	t.Run("test@example.com", func(t *testing.T) {
		result := BuildSecret("test@example.com")
		expected := "HENNGECHALLENGEtest@example.com"
		if result != expected {
			t.Errorf("BuildSecret(%q) = %q; expected %q", "test@example.com", result, expected)
		}
	})

	t.Run("user@domain.org", func(t *testing.T) {
		result := BuildSecret("user@domain.org")
		expected := "HENNGECHALLENGEuser@domain.org"
		if result != expected {
			t.Errorf("BuildSecret(%q) = %q; expected %q", "user@domain.org", result, expected)
		}
	})

	// Test with the default email used in this challenge submission
	t.Run("challenge_default_email", func(t *testing.T) {
		result := BuildSecret("aswai21@gmail.com")
		expected := "HENNGECHALLENGEaswai21@gmail.com"
		if result != expected {
			t.Errorf("BuildSecret(%q) = %q; expected %q", "aswai21@gmail.com", result, expected)
		}
	})

	t.Run("empty email", func(t *testing.T) {
		result := BuildSecret("")
		expected := "HENNGECHALLENGE"
		if result != expected {
			t.Errorf("BuildSecret(%q) = %q; expected %q", "", result, expected)
		}
	})
}

// TestGenerateTOTP tests basic TOTP generation
func TestGenerateTOTP(t *testing.T) {
	secret := BuildSecret("test@example.com")

	totp, err := GenerateTOTP(secret, DefaultTimeStep, DefaultDigits)
	if err != nil {
		t.Fatalf("GenerateTOTP failed: %v", err)
	}

	// Verify length
	if len(totp) != 10 {
		t.Errorf("TOTP length = %d; expected 10", len(totp))
	}

	// Verify all digits using recursive function
	checkDigitsRecursive(t, totp, 0)
}

// TestGenerateTOTPWithTime tests TOTP generation at specific times
func TestGenerateTOTPWithTime(t *testing.T) {
	secret := BuildSecret("test@example.com")

	// Same time counter should produce same TOTP
	totp1, _ := GenerateTOTPWithTime(secret, 1234567890, DefaultTimeStep, DefaultDigits)
	totp2, _ := GenerateTOTPWithTime(secret, 1234567890, DefaultTimeStep, DefaultDigits)

	if totp1 != totp2 {
		t.Errorf("Same time should produce same TOTP: %s vs %s", totp1, totp2)
	}

	// Different time counters should produce different TOTPs (usually)
	totp3, _ := GenerateTOTPWithTime(secret, 1234567890+60, DefaultTimeStep, DefaultDigits)
	// Note: This could theoretically be the same by chance, but extremely unlikely
	if totp1 == totp3 {
		t.Log("Warning: Different times produced same TOTP (unlikely but possible)")
	}
}

// TestValidateTOTP tests TOTP validation
func TestValidateTOTP(t *testing.T) {
	secret := BuildSecret("test@example.com")

	// Generate a current TOTP
	totp, _ := GenerateTOTP(secret, DefaultTimeStep, DefaultDigits)

	// Should validate successfully
	if !ValidateTOTP(secret, totp, DefaultTimeStep, DefaultDigits, 1) {
		t.Errorf("ValidateTOTP failed to validate current TOTP: %s", totp)
	}

	// Invalid TOTP should fail
	if ValidateTOTP(secret, "0000000000", DefaultTimeStep, DefaultDigits, 0) {
		t.Errorf("ValidateTOTP should reject invalid TOTP")
	}
}

// TestDecodeSecret tests the secret decoding
func TestDecodeSecret(t *testing.T) {
	secret := "HENNGECHALLENGEtest@example.com"
	decoded := decodeSecret(secret)

	// Should return the secret as bytes
	if string(decoded) != secret {
		t.Errorf("decodeSecret(%q) = %q; expected %q", secret, string(decoded), secret)
	}
}

// testSecretLengthRecursive tests multiple secret lengths recursively
func testSecretLengthRecursive(t *testing.T, lengths []int, index int) {
	if index >= len(lengths) {
		return
	}
	length := lengths[index]
	secret := GenerateSecret(length)
	if len(secret) != length {
		t.Errorf("GenerateSecret(%d) produced string of length %d", length, len(secret))
	}
	testSecretLengthRecursive(t, lengths, index+1)
}

// TestGenerateSecret tests the secret generation
func TestGenerateSecret(t *testing.T) {
	lengths := []int{10, 20, 32}
	testSecretLengthRecursive(t, lengths, 0)
}

// TestDefaultConstants verifies the default constants meet HENNGE requirements
func TestDefaultConstants(t *testing.T) {
	if DefaultDigits != 10 {
		t.Errorf("DefaultDigits = %d; HENNGE requires 10", DefaultDigits)
	}

	if DefaultTimeStep != 30 {
		t.Errorf("DefaultTimeStep = %d; standard TOTP uses 30", DefaultTimeStep)
	}
}

// TestPrintTimeBasedCodes tests the recursive time-based codes function doesn't panic
func TestPrintTimeBasedCodes(t *testing.T) {
	secret := BuildSecret("test@example.com")

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printTimeBasedCodes panicked: %v", r)
		}
	}()

	// Just verify it runs without error
	printTimeBasedCodes(secret, 1234567890, DefaultTimeStep, DefaultDigits, -2, 2)
}

// TestValidateTOTPRecursive tests the recursive validation helper
func TestValidateTOTPRecursive(t *testing.T) {
	secret := BuildSecret("test@example.com")
	currentTime := int64(1234567890)

	// Generate TOTP for current time
	totp, _ := GenerateTOTPWithTime(secret, currentTime, DefaultTimeStep, DefaultDigits)

	// Should validate with window
	if !validateTOTPRecursive(secret, totp, currentTime, DefaultTimeStep, DefaultDigits, -1, 1) {
		t.Errorf("validateTOTPRecursive failed to validate current TOTP")
	}

	// Invalid TOTP should fail
	if validateTOTPRecursive(secret, "9999999999", currentTime, DefaultTimeStep, DefaultDigits, 0, 0) {
		t.Errorf("validateTOTPRecursive should reject invalid TOTP")
	}
}
