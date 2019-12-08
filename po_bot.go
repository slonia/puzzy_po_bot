package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"gopkg.in/telegram-bot-api.v4"
)

type Configuration struct {
	Token string
}

var err error
var bot *tgbotapi.BotAPI
var puzzies = [...]string{"Тинки-Винки", "Дипси", "Ляля"}

func main() {
	reg, _ := regexp.Compile("^по")
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	fatal(err)
	bot, err = tgbotapi.NewBotAPI(config.Token)
	fatal(err)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	fatal(err)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		rand.Seed(time.Now().UTC().UnixNano())
		words := regexp.MustCompile("\\p{L}+").FindAllString(update.Message.Text, -1)
		po_words := []string{}
		for _, word := range words {
			matched := reg.MatchString(strings.ToLower(word))
			fatal(err)
			if utf8.RuneCountInString(word) > 3 && matched {
				po_words = append(po_words, word)
			}
		}
		probability := rand.Float32() + float32(len(po_words))/10 - 0.5
		if probability > 0.5 {
			word := po_words[rand.Intn(len(po_words))]
			word = stripText(word)
			word = string([]rune(word)[2 : len(word)-1])
			puzzy := puzzies[rand.Intn(len(puzzies))]
			verb := "говорил"
			if puzzy == "Ляля" {
				verb = "говорила"
			}
			text := fmt.Sprintf("Как %s %s: По, %s", verb, puzzy, word)
			log.Print(text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
		}
	}
}

func stripText(str string) string {
	for strings.HasSuffix(str, ",") || strings.HasSuffix(str, ".") || strings.HasSuffix(str, "!") || strings.HasSuffix(str, "?") {
		str = str[:len(str)-1]
	}
	return str
}

func fatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
