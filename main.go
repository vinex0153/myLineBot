package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	fmt.Println("hi")

	botToken := os.Getenv("botToken")
	botSecretKey := os.Getenv("botSecretKey")

	bot, err := linebot.New(botSecretKey, botToken)
	http.HandleFunc("/callback", response)
	fmt.Println(bot)
	fmt.Println(err)
}

func response(writer http.ResponseWriter, request *http.Request) {
	botEvents, err := bot.ParseRequest(request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			writer.WriteHeader(400)
		} else {
			writer.WriteHeader(500)
		}
		return
	}

	for _, botEvent := range botEvents {
		if botEvent.Type == linebot.EventTypeMessage {
			switch message := botEvent.Message.(type) {
			case *linebot.TextMessage:
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {

				}
				if _, err = bot.ReplyMessage(botEvent.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {

				}
			}
		}
	}
}
