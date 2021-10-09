package main

import (
	"bufio"
	"flag"
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

func init() {
	flag.StringVar(&Token, "t", "", "Bot token")
	flag.Parse()
}

func main() {
	quotes = ParseFile(path)
	dg, err := discordgo.New("Bot " + Token)
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
		output := fmt.Sprintf("--- Guess the quote! ---  \n%s", currentQuoute)
		sendChannelMessage(s, m, output)

	} else if strings.HasPrefix(m.Content, "!quote guess") {
		if m.Content == "!quote guess" {
			sendChannelMessage(s, m, "You have to guess, dummy!")
		} else {

			guess := m.Content[13:len(m.Content)]

			if strings.EqualFold(guess, currentAuthor) {
				output := fmt.Sprintf("Congratulations, %s is the correct person! :)", currentAuthor)

				sendChannelMessage(s, m, output)
			} else {
				output := fmt.Sprintf("I'm sorry, %s is not the correct person :(", guess)

				sendChannelMessage(s, m, output)
			}
		}
	} else if m.Content == "!quote all" {
		sendChannelMessage(s, m, GetAllQuotes())
	} else if m.Content == "!quote help" {
		output := fmt.Sprintf("Hello %s and welcome to Gøgler bot!\n\n"+
			"The following commands are supported at the moment:\n"+
			"!quote : Returns a random quote\n"+
			"!quote guess <name> : A guess to the last random quote\n"+
			"!quote all : Returns all quotes (but without the authors, otherwise that would be cheating!)\n"+
			"!tue : Returns \"Yo, fuck Tue!\"\n"+
			"!tue send : Sends \"Yo, fuck Tue!\" to Tue in a private dm", m.Author.Username)

		sendChannelMessage(s, m, output)
	} else if strings.HasPrefix(m.Content, "!quote ") {
		authorToSearchQuotesFor := m.Content[7:len(m.Content)]

		sendChannelMessage(s, m, GetAllQuotesFromAuthor(authorToSearchQuotesFor))
	} else if m.Content == "!tue" {
		sendChannelMessage(s, m, "Yo, fuck Tue!")
	} else if m.Content == "!tue send" {
		a, err := s.UserChannelCreate("245253768021540864")
		if err != nil {
			fmt.Printf("Error in sending to user: %v", err)
		}
		s.ChannelMessageSend(a.ID, "Yo, fuck Tue!")
		sendChannelMessage(s, m, "Successfully sent message to Tue")
	} else if m.Content == "!lasse" {
		sendChannelMessage(s, m, "Yo fuck Tue!")
	} else if m.Content == "!nød" {
		sendChannelMessage(s, m, ":weary:")
	}
}

func sendChannelMessage(s *discordgo.Session, m *discordgo.MessageCreate, output string) {
	_, err := s.ChannelMessageSend(m.ChannelID, output)
	if err != nil {
		log.Fatal(err)
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

func GetAllQuotes() string {
	var allQuotes strings.Builder
	for key := range quotes {
		allQuotes.WriteString(fmt.Sprintf("\"%s \"\n", key))
	}
	return allQuotes.String()
}

func GetAllQuotesAndAuthors() string {
	var allQuotes strings.Builder
	for key, value := range quotes {
		allQuotes.WriteString(fmt.Sprintf("\"%s \" -%s \n", key, value))
	}
	return allQuotes.String()
}

func GetAllQuotesFromAuthor(author string) string {
	var allQuotes strings.Builder
	allQuotes.WriteString(fmt.Sprintf("%s has said all of these quotes:\n", strings.Title(author)))

	for key, value := range quotes {
		if strings.EqualFold(value, author) {
			allQuotes.WriteString(fmt.Sprintf("\"%s \"\n", key))
		}
	}

	if allQuotes.Len() == 1 {
		return fmt.Sprintf("Hmm, doesn't look like %s has any quotes :(", author)
	} else {
		return allQuotes.String()
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
