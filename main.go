package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

var convetedData []map[string]interface{}

func main() {
	getEnvData()

	var err error
	// botToken := os.Getenv("botToken")
	// botSecretKey := os.Getenv("botSecretKey")
	// log.Printf(botToken)
	// log.Printf(botSecretKey)
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
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
				//quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				if _, err = bot.ReplyMessage(botEvent.ReplyToken, linebot.NewTextMessage(findData(message.Text))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func getEnvData() {
	resp, err := http.Get("https://data.epa.gov.tw/api/v1/aqx_p_432?limit=1000&api_key=9be7b239-557b-4c10-9775-78cadfc555e9&format=json")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	data := string(body)

	var jsonObj map[string]interface{}
	json.Unmarshal([]byte(data), &jsonObj)

	records := jsonObj["records"].([]interface{})
	//fmt.Println(records)

	for _, r := range records {
		record := r.(map[string]interface{})
		//fmt.Println(record["SiteName"])
		convetedData = append(convetedData, record)
	}
}

func findData(site string) string {
	for _, c := range convetedData {
		if c["SiteName"].(string) == site {
			return c["AQI"].(string) + "時間:" + c["ImportDate"].(string)
		}
	}

	return "NoData"
}
