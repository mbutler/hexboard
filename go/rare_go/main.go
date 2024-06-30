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
	"Frequent": "0",
	"Ordinary": "00",
	"Common":   "000",
	// Uncomment the below lines for higher rarity levels if you wish to simulate them
	"Uncommon": "0000",
	"Rare":     "00000",
	"Epic":     "000000",
	// "Legendary": "0000000",
	// "Mythic":    "00000000",
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
		attempts++

		if rarity == targetRarity {
			return attempts, hashStr
		}
	}

	return attempts, ""
}

func main() {
	rarityLevelsToTest := []string{"Frequent", "Ordinary", "Common", "Uncommon", "Rare", "Epic"}
	tokenLimit := 10000000
	runs := 30

	for _, rarity := range rarityLevelsToTest {
		totalAttempts := 0

		fmt.Printf("Testing rarity level: %s\n", rarity)
		for i := 0; i < runs; i++ {
			attempts, _ := findRareItemWithTokens(rarity, tokenLimit)
			totalAttempts += attempts
		}

		averageAttempts := totalAttempts / runs
		fmt.Printf("Average attempts to find %s item: %d\n", rarity, averageAttempts)
	}
}
