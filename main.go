package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	log.Println("hi")

	var err error
	// botToken := os.Getenv("botToken")
	// botSecretKey := os.Getenv("botSecretKey")
	// log.Printf(botToken)
	// log.Printf(botSecretKey)
	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", response)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func response(writer http.ResponseWriter, request *http.Request) {
	log.Printf("response")
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
					log.Println("Quota err:", err)
				}
				if _, err = bot.ReplyMessage(botEvent.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
