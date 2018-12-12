package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

type CustomMessage struct {
}

func main() {
	handler, err := httphandler.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {

		log.Print(r)

		bot, err := handler.NewClient()
		if err != nil {
			log.Print(err)
		}
		for _, event := range events {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				b, _ := linebot.NewTextMessage(message.Text).MarshalJSON()
				log.Print(string(b))
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	})

	http.Handle("/callback", handler)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
