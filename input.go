package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// ScanStringWithPrompt : Print the input, and using the provided
// scanner, scan a string and return it
func ScanStringWithPrompt(message string, scanner *bufio.Scanner) string {
	fmt.Print(message)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// ScanYesNoQuestion : Print the question, scan an answer and parse it into a bool
func ScanYesNoQuestion(question string, defaultAns rune, scanner *bufio.Scanner) bool {
	answerOptions := "[y/n] "
	if defaultAns == 'y' {
		answerOptions = "[Y/n] "
	} else if defaultAns == 'n' {
		answerOptions = "[y/N] "
	}
	prompt := fmt.Sprintf("%s %s", question, answerOptions)
	for {
		switch ScanStringWithPrompt(prompt, scanner) {
		case "Y", "y":
			return true
		case "N", "n":
			return false
		case "":
			if defaultAns == 'y' {
				return true
			} else if defaultAns == 'n' {
				return false
			}
		default:
			fmt.Println("Sorry, an answer couldn't be determined")
		}
	}
}

// ScanNumberWithPrompt : Print the input, and using the provided
// scanner, scan a string, parse it as a float64 and return it
func ScanNumberWithPrompt(message string, scanner *bufio.Scanner) float64 {
	for {
		maybeNumber := ScanStringWithPrompt(message, scanner)
		newValue, err := strconv.ParseFloat(maybeNumber, 64)
		if err != nil {
			fmt.Printf("Value '%s' couldn't be parsed as a number\n", maybeNumber)
		} else {
			return newValue
		}
	}
}

// ScanStringArray : Scans strings while the user hasn't left the input empty
func ScanStringArray(scanner *bufio.Scanner) []interface{} {
	var stringArr []interface{}
	message := "Write a new element or leave it empty to finish adding elements: "
	newElem := ScanStringWithPrompt(message, scanner)
	for newElem != "" {
		stringArr = append(stringArr, newElem)
		newElem = ScanStringWithPrompt(message, scanner)
	}
	return stringArr
}
