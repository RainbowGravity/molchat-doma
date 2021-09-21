package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type JsonData struct {
	Name string `json:"name"`
	Html string `json:"html_url"`
}

var reading []JsonData

func main() {
	var msgHmrwk string

	bot, err := tgbotapi.NewBotAPI("____________________TOKEN_____________________")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	gitRepos := "https://api.github.com/repos/"
	gitUser := "RainbowGravity/"
	gitRepo := "course"
	gitUrl := gitRepos + gitUser + gitRepo

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "markdown"
			switch update.Message.Command() {
			case "git":
				msg.Text = "*Here is the link to my repository:* \n\n[RainbowGravity/course](https://github.com/RainbowGravity/course)"
			case "tasks":
				msgHmrwk, err = completedHomework(gitUrl)
				msg.Text = errorHandling(msgHmrwk, err)
			case "task":
				hmwrkStr := (update.Message.CommandArguments())
				msgHmrwk, err = specifiedHomework(gitUrl, hmwrkStr)
				msg.Text = errorHandling(msgHmrwk, err)
			default:
				err = errors.New("there is no *" + update.Message.Text + "* command. \n\n*Try one of these:*\n*/git* — Link to my Github repository;\n*/tasks* — List of my completed homework;\n*/task* — Specified homework (*e.g. /task 2*)")
				msg.Text = errorHandling(msgHmrwk, err)
				err = nil
			}
			bot.Send(msg)
		}

	}
}

func getContents(gitUrl string) {

	resp, err := http.Get(gitUrl + "/contents/")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	jsonString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(jsonString), &reading)
	if err != nil {
		fmt.Println(err)
	}
}

func completedHomework(gitUrl string) (hmwrkAll string, err error) {

	getContents(gitUrl)

	for _, val := range reading {
		if !strings.Contains(val.Name+" — "+val.Html, "README.md") && !strings.Contains(val.Name+" — "+val.Html, "WIP") {
			//[RainbowGravity/course](https://github.com/RainbowGravity/course)
			tmp := string("\n" + "[" + val.Name + "]" + "(" + val.Html + ")")
			hmwrkAll = hmwrkAll + tmp
			if err != nil {
				err = errors.New("there is no homework")
			}
		}
	}
	return hmwrkAll, err
}

func specifiedHomework(gitUrl string, hmwrkStr string) (hmwrkSpc string, err error) {

	getContents(gitUrl)
	hmwrkNum, err := strconv.Atoi(hmwrkStr)
	hmwrkNum = hmwrkNum - 1

	if hmwrkNum < len(reading)-1 && hmwrkNum > -1 {
		hmwrkSpc = ("\n" + "[" + reading[hmwrkNum].Name + "]" + "(" + reading[hmwrkNum].Html + ")")
	} else {
		if hmwrkNum <= -1 {
			err = errors.New("there is no homework with number less than 1")
		} else {
			err = errors.New("there is no homework with number more than " + fmt.Sprint(len(reading)-1))
		}
	}
	return hmwrkSpc, err
}

//error handling function
func errorHandling(msgHmrwk string, err error) (botMessage string) {
	if err != nil {
		botMessage = ("*An error occurred: *" + fmt.Sprint(err) + ".")
	} else {
		botMessage = ("*List of completed homework:* \n " + msgHmrwk)
	}
	return
}