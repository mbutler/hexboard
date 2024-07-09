package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

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

func findRareItemWithTokens(targetRarity string) (int, float64, string, string) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	attempts := 0
	startTime := time.Now()

	for {
		randomStr := generateRandomString(rng, 10)
		hashStr := hashString(randomStr)
		rarity := checkRarity(hashStr)
		attempts++

		if rarity == targetRarity {
			elapsedTime := time.Since(startTime).Seconds()
			return attempts, elapsedTime, hashStr, randomStr
		}
	}
}

func main() {
	rarityLevelsToTest := []string{"Frequent", "Ordinary", "Common", "Uncommon", "Rare", "Epic", "Legendary", "Mythic"}

	for _, rarity := range rarityLevelsToTest {
		fmt.Printf("Finding item of rarity level: %s\n", rarity)
		attempts, elapsedTime, hashStr, randomStr := findRareItemWithTokens(rarity)
		fmt.Printf("Found %s item!\n", rarity)
		fmt.Printf("Random String: %s\n", randomStr)
		fmt.Printf("Hash: %s\n", hashStr)
		fmt.Printf("Attempts: %d\n", attempts)
		fmt.Printf("Time: %.2f seconds\n", elapsedTime)
		fmt.Println()
	}
}
