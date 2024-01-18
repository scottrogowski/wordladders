package main

import (
	"fmt"
	"io/ioutil"
	"bytes"
	// json
	"encoding/json"
	// rand
	// "math/rand"
	// cli args
	"os"
	// sort
	"sort"
)

// str is a type which is []byte
type Ladder []string


func isLadder(word1, word2 string) bool {
	diff := 0
	for i := 0; i < len(word1); i++ {
		if word1[i] != word2[i] {
			diff++
		}
	}
	return diff == 1
}

func filter(words []string, fn func(string) bool) []string {
	filtered := make([]string, 0)
	for _, word := range words {
		if fn(word) {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func main() {
	// filename is the first argument
	filename := os.Args[1]

	// read words.txt
	words, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	// split by newline into an array of bytes
	allWords := bytes.Split(words, []byte("\r\n"))
	// filter for words that are 5 characters long
	fiveWords := []string{}
	for _, word := range allWords {
		if len(word) == 5 {
			fiveWords = append(fiveWords, string(word))
		}
	}

	// print length of array
	fmt.Println("Found ", len(fiveWords), " words that are 5 characters long.")

	// generate a map of word pairs which we can use to find ladders
	// e.g. "hello" -> ["jello", "cello", ...]
	wordPairs := make(map[string][]string)
	for _, word1 := range fiveWords {
		for _, word2 := range fiveWords {
			if isLadder(word1, word2) {
				wordPairs[word1] = append(wordPairs[word1], word2)
			}
		}
	}
	fmt.Println("Word pairs: ", wordPairs)

	fmt.Println("Found ", len(wordPairs), " words with >=1 pair.")

	wordPairKeys := make([]string, 0, len(wordPairs))
	for k := range wordPairs {
		wordPairKeys = append(wordPairKeys, k)
	}

	// generate all possible ladders
	fmt.Println("Generating ladders..")
	allLadders := make([]Ladder, 0)
	for i, startWord := range wordPairKeys {
		if i % 100 == 0 {
			fmt.Printf("Generating ladders for word %d/%d\n", i, len(wordPairKeys))
		}
		startLadder := Ladder{startWord}
		// fmt.Println("Start: ", startWord)
		ladders := findLadders(startLadder, wordPairs)
		// fmt.Println("Found ", len(ladders), " ladders.")
		allLadders = append(allLadders, ladders...)
	}
	fmt.Println("Found ", len(allLadders), " ladders.")

	// filter to shortest ladders of each beginning/end pair
	ladderByPair := make(map[string]Ladder)
	for _, ladder := range allLadders {
		key := ladder[0] + ":" + ladder[len(ladder)-1]
		if _, ok := ladderByPair[key]; !ok {
			ladderByPair[key] = ladder
		} else if len(ladder) < len(ladderByPair[key]) {
			ladderByPair[key] = ladder
		}
	}
	allLadders = make([]Ladder, 0, len(ladderByPair))
	for _, ladder := range ladderByPair {
		allLadders = append(allLadders, ladder)
	}

	// filter out all ladders with length == 2. These are trivial
	laddersWithoutLength2 := make([]Ladder, 0)
	for _, ladder := range allLadders {
		if len(ladder) > 2 {
			laddersWithoutLength2 = append(laddersWithoutLength2, ladder)
		}
	}
	allLadders = laddersWithoutLength2

	// filter out all ladders where the first letter of the first word is = to the first letter of the last word
	// these are less interesting
	laddersWithoutSameFirstLetter := make([]Ladder, 0)
	for _, ladder := range allLadders {
		if ladder[0][0] != ladder[len(ladder)-1][0] {
			laddersWithoutSameFirstLetter = append(laddersWithoutSameFirstLetter, ladder)
		}
	}
	allLadders = laddersWithoutSameFirstLetter

	// sort the latters by length and then alphabetically
	sort.Slice(allLadders, func(i, j int) bool {
		if len(allLadders[i]) == len(allLadders[j]) {
			return allLadders[i][0] < allLadders[j][0]
		}
		return len(allLadders[i]) < len(allLadders[j])
	})

	// save it as ladders.json
	fmt.Println("After filtering, found ", len(allLadders), " ladders.")
	jsonLadders, err := json.Marshal(allLadders)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("ladders.json", jsonLadders, 0644)
	if err != nil {
		panic(err)
	}
}

func findLadders(ladder Ladder, wordPairs map[string][]string) []Ladder {
	possibleWords := wordPairs[ladder[len(ladder)-1]]
	// exclude words that are already in the ladder
	possibleWords = filter(possibleWords, func(word string) bool {
		for _, w := range ladder {
			if w == word {
				return false
			}
		}
		return true
	})
	// if there are no possible words, return the ladder
	if len(possibleWords) == 0 {
		return []Ladder{ladder}
	}
	// if the length of the ladder is 5, just return the ladder with the last word
	ladders := []Ladder{}
	if len(ladder) == 7 {
		for _, word := range possibleWords {
			ladders = append(ladders, append(ladder, word))
		}
		return ladders
	}

	// otherwise, recurse
	for _, word := range possibleWords {
		ladders = append(ladders, findLadders(append(ladder, word), wordPairs)...)
	}
	// add the initial ladder
	ladders = append(ladders, ladder)
	return ladders
}
