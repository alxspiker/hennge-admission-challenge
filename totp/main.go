package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

func main() {
	// secret := "MFZXOYLJGIYUAZ3NMFUWYLTDN5WUQRKOJZDUKQ2IIFGEYRKOI5CTAMBU" // Base32 encoded secret
	secret := "aswai21@gmail.comHENNGECHALLENGE004" // Base32 encoded secret
	timeStep := int64(30)                           // 30 seconds (standard)
	digits := 10                                    // 10-digit codes (standard)

	fmt.Println("=== TOTP Demo ===")
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

	// Show how codes change over time
	fmt.Println("\n=== Time-based Code Changes ===")
	for i := -2; i <= 2; i++ {
		testTime := now + int64(i)*timeStep
		testCounter := testTime / timeStep
		testTOTP, _ := generateHOTP(decodeSecret(secret), testCounter, digits)

		timeStr := time.Unix(testTime, 0).Format("15:04:05")
		if i == 0 {
			fmt.Printf("-> %s: %s (current)\n", timeStr, testTOTP)
		} else {
			fmt.Printf("   %s: %s\n", timeStr, testTOTP)
		}
	}
}

// TOTP generates a Time-based One-Time Password according to RFC 6238
func GenerateTOTP(secret string, timeStep int64, digits int) (string, error) {
	// Step 1: Calculate time counter T
	// T = (Current Unix Time - T0) / X
	// Where T0 = 0 (Unix epoch) and X = time step (usually 30 seconds)
	now := time.Now().Unix()
	timeCounter := now / timeStep

	// Step 2: Generate HOTP using the time counter
	return generateHOTP(decodeSecret(secret), timeCounter, digits)
}

// generateHOTP implements HOTP algorithm from RFC 4226
func generateHOTP(key []byte, counter int64, digits int) (string, error) {
	// Step 1: Convert counter to 8-byte big-endian format
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(counter))

	// Step 2: Compute HMAC-SHA1
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

	// Step 5: Format with leading zeros
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, otp), nil
}

// ValidateTOTP validates a TOTP code with time window tolerance
func ValidateTOTP(secret, code string, timeStep int64, digits, windowSize int) bool {
	currentTime := time.Now().Unix()

	// Check current time window and adjacent windows for clock skew tolerance
	for i := -windowSize; i <= windowSize; i++ {
		testTime := currentTime + int64(i)*timeStep
		timeCounter := testTime / timeStep

		expectedCode, err := generateHOTP(decodeSecret(secret), timeCounter, digits)
		if err != nil {
			continue
		}

		if expectedCode == code {
			return true
		}
	}
	return false
}

// Helper function to decode base32 secret
// func decodeSecret(secret string) []byte {
// 	key, _ := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
// 	return key
// }

func decodeSecret(secret string) []byte {
	return []byte(secret)
}

// GenerateSecret creates a random base32-encoded secret
func GenerateSecret(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	secret := make([]byte, length)

	// In production, use crypto/rand for secure random generation
	for i := range secret {
		secret[i] = charset[i%len(charset)] // Simplified for demo
	}
	return string(secret)
}
