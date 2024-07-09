package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/*
Expected tokens spent to see each rarity level:

Frequent (1 leading zero):    Expected tokens = 16
Ordinary (2 leading zeros):   Expected tokens = 256
Common (3 leading zeros):     Expected tokens = 4096
Uncommon (4 leading zeros):   Expected tokens = 65536
Rare (5 leading zeros):       Expected tokens = 1048576
Epic (6 leading zeros):       Expected tokens = 16777216
Legendary (7 leading zeros):  Expected tokens = 268435456
Mythic (8 leading zeros):     Expected tokens = 4294967296
*/

var rarityLevels = map[string]string{
	"Frequent":  "0",
	"Ordinary":  "00",
	"Common":    "000",
	"Uncommon":  "0000",
	"Rare":      "00000",
	"Epic":      "000000",
	"Legendary": "0000000",
	"Mythic":    "00000000",
}

func generateRandomString(rng *rand.Rand, length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rng.Intn(len(letters))]
	}
	return string(result)
}

func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func checkRarity(hash string) string {
	for rarity, pattern := range rarityLevels {
		if strings.HasPrefix(hash, pattern) {
			return rarity
		}
	}
	return "None"
}

func findRareItemWithTokens(targetRarity string, tokenLimit int) (int, string) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	attempts := 0

	for attempts < tokenLimit {
		randomStr := generateRandomString(rng, 10)
		hashStr := hashString(randomStr)
		rarity := checkRarity(hashStr)
		fmt.Printf("Attempt %d: String: %s, Hash: %s, Rarity: %s\n", attempts+1, randomStr, hashStr, rarity) // Debugging output
		attempts++

		if rarity == targetRarity {
			return attempts, hashStr
		}
	}

	return attempts, ""
}

func main() {
	targetRarity := "Ordinary"
	tokenLimit := 10000
	runs := 30
	totalAttempts := 0

	for i := 0; i < runs; i++ {
		attempts, hashStr := findRareItemWithTokens(targetRarity, tokenLimit)
		if hashStr != "" {
			fmt.Printf("Run %d: Found %s item in %d attempts.\n", i+1, targetRarity, attempts)
		} else {
			fmt.Printf("Run %d: Did not find %s item within the token limit.\n", i+1, targetRarity)
		}
		totalAttempts += attempts
	}

	averageAttempts := totalAttempts / runs
	fmt.Printf("Average attempts to find %s item: %d\n", targetRarity, averageAttempts)
}
