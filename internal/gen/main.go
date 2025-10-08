// Package main generates tries of Unicode properties for string width calculation
package main

import (
	"fmt"
	"log"
	"path/filepath"
)

func main() {
	fmt.Println("Generating string width trie...")

	// Parse Unicode data
	data, err := ParseUnicodeData()
	if err != nil {
		log.Fatalf("Failed to parse Unicode data: %v", err)
	}

	// Generate trie
	trie, err := GenerateTrie(data)
	if err != nil {
		log.Fatalf("Failed to generate trie: %v", err)
	}

	// Write trie to output file
	outputPath := filepath.Join("..", "..", "trie.go")
	if err := WriteTrieGo(trie, outputPath); err != nil {
		log.Fatalf("Failed to write trie: %v", err)
	}

	fmt.Println("Trie generation completed successfully!")
}
