package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

var quotes map[string]string
var path string = "data\\quotes.txt"
var currentAuthor string
var currentQuoute string

// func init() {
// 	flag.StringVar(&Token, "t", "", "Bot token")
// 	flag.Parse()
// }

func main() {
	quotes = ParseFile(path)
	dg, err := discordgo.New("Bot " + "ODk1ODAyNDgzMTE1NDU0NDY0.YV93Ew.M-yOKKwLAAh357ohQPps-U2sSUQ")
	if err != nil {
		log.Fatalf("Error creating Discord session %v", err)
	}
	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection %v", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!quote" {
		currentQuoute, currentAuthor = FindRandomQuote()
		_, err := s.ChannelMessageSend(m.ChannelID, currentQuoute)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.Contains(m.Content, "!quote guess") {
		words := strings.Fields(m.Content)
		guess := words[len(words)-1]

		if strings.EqualFold(guess, currentAuthor) {
			output := fmt.Sprintf("Congratulations, %s is the correct person! :)", currentAuthor)

			_, err := s.ChannelMessageSend(m.ChannelID, output)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			output := fmt.Sprintf("I'm sorry, %s is not the correct person :(", guess)

			_, err := s.ChannelMessageSend(m.ChannelID, output)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func ParseFile(path string) map[string]string {
	f, err := os.Open("data\\quotes.txt")
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
