package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

// DefaultEmail is the email to use for TOTP generation
const DefaultEmail = "aswai21@gmail.com"

// DefaultDigits is the number of digits for TOTP output (HENNGE requires 10)
const DefaultDigits = 10

// DefaultTimeStep is the time step in seconds (standard TOTP)
const DefaultTimeStep = int64(30)

func main() {
	// Secret key must be constructed as "HENNGECHALLENGE" + email_address
	secret := BuildSecret(DefaultEmail)
	timeStep := DefaultTimeStep
	digits := DefaultDigits

	fmt.Println("=== TOTP Demo ===")
	fmt.Printf("Email: %s\n", DefaultEmail)
	fmt.Printf("Secret: %s\n", secret)
	fmt.Printf("Time Step: %d seconds\n", timeStep)
	fmt.Printf("Digits: %d\n\n", digits)

	// Generate current TOTP
	totp, err := GenerateTOTP(secret, timeStep, digits)
	if err != nil {
		fmt.Printf("Error generating TOTP: %v\n", err)
		return
	}

	fmt.Printf("Current TOTP: %s\n", totp)
	fmt.Printf("Generated at: %s\n", time.Now().Format("15:04:05"))

	// Show time counter calculation
	now := time.Now().Unix()
	timeCounter := now / timeStep
	fmt.Printf("Current Unix time: %d\n", now)
	fmt.Printf("Time counter (T): %d\n", timeCounter)

	// Validate the generated code
	isValid := ValidateTOTP(secret, totp, timeStep, digits, 1)
	fmt.Printf("Validation result: %t\n", isValid)

	// Show how codes change over time (using recursion, no for loops)
	fmt.Println("\n=== Time-based Code Changes ===")
	printTimeBasedCodes(secret, now, timeStep, digits, -2, 2)
}

// printTimeBasedCodes prints TOTP codes for a range of time offsets using recursion
func printTimeBasedCodes(secret string, baseTime int64, timeStep int64, digits int, current int, end int) {
	if current > end {
		return
	}

	testTime := baseTime + int64(current)*timeStep
	testCounter := testTime / timeStep
	testTOTP, _ := generateHOTP(decodeSecret(secret), testCounter, digits)

	timeStr := time.Unix(testTime, 0).Format("15:04:05")
	if current == 0 {
		fmt.Printf("-> %s: %s (current)\n", timeStr, testTOTP)
	} else {
		fmt.Printf("   %s: %s\n", timeStr, testTOTP)
	}

	// Recursive call for next iteration
	printTimeBasedCodes(secret, baseTime, timeStep, digits, current+1, end)
}

// BuildSecret constructs the secret key as per HENNGE requirements:
// "HENNGECHALLENGE" + email_address
func BuildSecret(email string) string {
	return "HENNGECHALLENGE" + email
}

// GenerateTOTP generates a Time-based One-Time Password according to RFC 6238
// Uses HMAC-SHA512 as required by HENNGE challenge
func GenerateTOTP(secret string, timeStep int64, digits int) (string, error) {
	// Step 1: Calculate time counter T
	// T = (Current Unix Time - T0) / X
	// Where T0 = 0 (Unix epoch) and X = time step (usually 30 seconds)
	now := time.Now().Unix()
	timeCounter := now / timeStep

	// Step 2: Generate HOTP using the time counter
	return generateHOTP(decodeSecret(secret), timeCounter, digits)
}

// GenerateTOTPWithTime generates a TOTP for a specific Unix timestamp
// This is useful for testing and validation
func GenerateTOTPWithTime(secret string, unixTime int64, timeStep int64, digits int) (string, error) {
	timeCounter := unixTime / timeStep
	return generateHOTP(decodeSecret(secret), timeCounter, digits)
}

// generateHOTP implements HOTP algorithm from RFC 4226
// Uses HMAC-SHA512 (NOT SHA1) as required by HENNGE challenge
func generateHOTP(key []byte, counter int64, digits int) (string, error) {
	// Step 1: Convert counter to 8-byte big-endian format
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(counter))

	// Step 2: Compute HMAC-SHA512 (CRITICAL: Must be SHA512, not SHA1)
	mac := hmac.New(sha512.New, key)
	mac.Write(counterBytes)
	hash := mac.Sum(nil)

	// Step 3: Dynamic truncation
	// Extract 4 bytes from the hash using dynamic truncation
	offset := hash[len(hash)-1] & 0x0F // Last 4 bits of hash

	// Extract 4 bytes starting from offset
	truncatedHash := binary.BigEndian.Uint32(hash[offset : offset+4])

	// Clear the most significant bit (sign bit)
	truncatedHash &= 0x7FFFFFFF

	// Step 4: Compute the final OTP
	otp := truncatedHash % uint32(math.Pow10(digits))

	// Step 5: Format with leading zeros (CRITICAL: Must be padded)
	// e.g., 0012345678 not 12345678
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, otp), nil
}

// ValidateTOTP validates a TOTP code with time window tolerance
// Uses recursion instead of for loops
func ValidateTOTP(secret, code string, timeStep int64, digits, windowSize int) bool {
	currentTime := time.Now().Unix()
	return validateTOTPRecursive(secret, code, currentTime, timeStep, digits, -windowSize, windowSize)
}

// validateTOTPRecursive is the recursive helper for TOTP validation
func validateTOTPRecursive(secret, code string, currentTime int64, timeStep int64, digits int, offset int, maxOffset int) bool {
	if offset > maxOffset {
		return false
	}

	testTime := currentTime + int64(offset)*timeStep
	timeCounter := testTime / timeStep

	expectedCode, err := generateHOTP(decodeSecret(secret), timeCounter, digits)
	if err == nil && expectedCode == code {
		return true
	}

	// Recursive call for next offset
	return validateTOTPRecursive(secret, code, currentTime, timeStep, digits, offset+1, maxOffset)
}

// decodeSecret converts the secret string to bytes
// For HENNGE challenge, the secret is used directly as bytes (not Base32 encoded)
func decodeSecret(secret string) []byte {
	return []byte(secret)
}

// GenerateSecret creates a random base32-encoded secret using recursion
func GenerateSecret(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	secret := make([]byte, length)
	generateSecretRecursive(secret, charset, 0)
	return string(secret)
}

// generateSecretRecursive is the recursive helper for secret generation
func generateSecretRecursive(secret []byte, charset string, index int) {
	if index >= len(secret) {
		return
	}
	// Simplified for demo - in production use crypto/rand
	secret[index] = charset[index%len(charset)]
	generateSecretRecursive(secret, charset, index+1)
}
