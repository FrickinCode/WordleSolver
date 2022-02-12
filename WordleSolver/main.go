package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var initialGuess string = "petal"

func main() {
	var word string
	fmt.Print("Input a word to guess: ")
	fmt.Scanln(&word)

	if len(word) != 5 {
		fmt.Println("Input word must be 5 characters long")
		return
	}

	word = strings.ToLower(word)

	validWords := readDictionary()
	guessWords := validWords
	var guesses [6]string

	var guess string
	var isCorrect bool
	var results [6][5]CharResult

	for i := 0; i < 6; i++ {
		var singleResult [5]CharResult
		var correctPositions int
		var correctLetters int

		if len(validWords) == 0 {
			fmt.Println("No valid words found")
			return
		}

		var corrects int = correctLetters + correctPositions
		if i == 0 {
			guess = initialGuess
		} else if corrects >= 3 && correctPositions >= 1 {
			guess = validWords[0]
		} else if len(guessWords) > 0 {
			guess = guessWords[0]
		} else {
			guess = validWords[0]
		}
		guesses[i] = guess

		isCorrect, singleResult = doGuess(word, guess)

		results[i] = singleResult
		if isCorrect {
			break
		}

		correctLetters = 0
		correctPositions = 0

		for j := 0; j < 5; j++ {
			if singleResult[j] == CorrectPosition {
				validWords = filterDictionary(validWords, guess[j], j, comparerIsCorrectLetterPosition)
				guessWords = filterDictionary(guessWords, guess[j], j, comparerStringContainsLetter)
				correctPositions++
			} else if singleResult[j] == CorrectLetter {
				validWords = filterDictionary(validWords, guess[j], j, comparerStringContainsLetter)
				guessWords = filterDictionary(guessWords, guess[j], j, comparerRemoveLetterPosition)
				correctLetters++
			} else if singleResult[j] == NoMatch {
				validWords = filterDictionary(validWords, guess[j], j, comparerStringDoesNotContainLetter)
				guessWords = filterDictionary(guessWords, guess[j], j, comparerStringDoesNotContainLetter)
			}
		}

	}

	if isCorrect {
		fmt.Println("Successfully solved the wordle!")
	} else {
		fmt.Println("Was unable to find the correct word.")
	}

	fancyPrint(guesses, results)
}

func readDictionary() []string {
	file, err := os.Open("words.txt")
	checkError(err)

	scanner := bufio.NewScanner(file)
	wordList := make([]string, 0)

	for scanner.Scan() {
		wordList = append(wordList, scanner.Text())
	}

	return wordList

}

func filterDictionary(dict []string, c byte, index int, comparer func(string, byte, int) bool) []string {
	var i int = 0
	for _, s := range dict {
		if comparer(s, c, index) {
			dict[i] = s
			i++
		} else {
			//fmt.Println("Removing", s)
		}
	}

	dict = dict[:i]
	return dict
}

func comparerIsCorrectLetterPosition(s string, c byte, index int) bool {
	return s[index] == c
}

func comparerRemoveLetterPosition(s string, c byte, index int) bool {
	return s[index] != c
}

func comparerStringContainsLetter(s string, c byte, index int) bool {
	return strings.Contains(s, string(c))
}

func comparerStringDoesNotContainLetter(s string, c byte, index int) bool {
	return !strings.Contains(s, string(c))
}

func doGuess(ans string, guess string) (bool, [5]CharResult) {
	var isCorrect bool = false
	var results [5]CharResult
	if ans == guess {
		isCorrect = true
	}

	for index := range guess {
		if ans[index] == guess[index] {
			results[index] = CorrectPosition
		} else if strings.Contains(ans, string(guess[index])) {
			results[index] = CorrectLetter
		} else {
			results[index] = NoMatch
		}
	}

	return isCorrect, results
}

func fancyPrint(guesses [6]string, results [6][5]CharResult) {
	for index, s := range guesses {
		if s != "" {
			fmt.Println("Guess", index+1, ":", s)
		}
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

type CharResult int32

const (
	NoMatch         CharResult = 0
	CorrectLetter   CharResult = 1
	CorrectPosition CharResult = 2
)
