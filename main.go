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

func (m *CustomMessage) Message() {
}

var testMsg = `{
	"type": "template",
	"altText": "this is a confirm template",
	"template": {
	  "type": "confirm",
	  "actions": [
		{
		  "type": "uri",
		  "label": "Yes",
		  "uri": "http://google.com"
		},
		{
		  "type": "message",
		  "label": "No",
		  "text": "No"
		}
	  ],
	  "text": "Anything else ? "
	},
	"quickReply": {
		"items": [
		  {
			"type": "action", 
			"imageUrl": "https://example.com/sushi.png",
			"action": {
			  "type": "message",
			  "label": "Sushi",
			  "text": "Sushi"
			}
		  },
		  {
			"type": "action",
			"imageUrl": "https://example.com/tempura.png",
			"action": {
			  "type": "message",
			  "label": "Tempura",
			  "text": "Tempura"
			}
		  },
		  {
			"type": "action",
			"action": {
			  "type": "location",
			  "label": "Send location"
			}
		  }
		]
	  }
  }`

func (m *CustomMessage) MarshalJSON() ([]byte, error) {
	return []byte(testMsg), nil
}

// WithQuickReplies method of CustomMessage
func (m *CustomMessage) WithQuickReplies(items *linebot.QuickReplyItems) linebot.SendingMessage {
	return m
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
				message.Message()
				msg := new(CustomMessage)
				b, _ := msg.MarshalJSON()
				log.Print(string(b))
				if _, err = bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
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
