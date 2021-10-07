package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var quotes map[string]string

func main() {
	path := "data\\quote.txt"
	quotes = ParseFile(path)

	GuessRandomQuote()
}

func ParseFile(path string) map[string]string {
	f, err := os.Open("data\\quote.txt")
	if err != nil {
		log.Fatalf("Error in opening file: %v", err)
	}

	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)

	quouteToAuthor := make(map[string]string)

	for sc.Scan() {

		line := sc.Text()
		split := strings.Split(line, ":")

		quouteToAuthor[split[0]] = split[1]
	}
	return quouteToAuthor
}

func PrintAllQuotes() {
	for key, value := range quotes {
		fmt.Printf("\"%s \" ~%s \n", key, value)
	}
}

func FindRandomQuote() (string, string) {
	rand.Seed(time.Now().UnixNano())

	var list []string
	for key := range quotes {
		list = append(list, key)
	}

	randomNumber := rand.Intn((len(list) - 1) + 1)

	return fmt.Sprintf("\"%s\"", list[randomNumber]), quotes[list[randomNumber]]
}

func GuessRandomQuote() {
	quote, author := FindRandomQuote()

	fmt.Println("----- Gæt et quote! -----")
	fmt.Println(quote)

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	sentence, err := buf.ReadBytes('\n')
	result := strings.TrimSpace(string(sentence))
	if err != nil {
		fmt.Println(err)
	} else {
		if strings.EqualFold(result, author) {
			fmt.Println("Tillykke, du gættede korrekt!")
		} else {
			fmt.Print("Forkert gæt :(")
		}
	}
}
