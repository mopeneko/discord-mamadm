package main

import (
	"context"
	"log"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/greymd/mamadm/generator"
)

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_TOKEN"),
		Intents:  disgord.IntentGuildMessages,
	})

	defer client.Gateway().StayConnectedUntilInterrupted()

	client.Gateway().
		MessageCreate(func(s disgord.Session, event *disgord.MessageCreate) {
			if event.Message.Content != "!mama" {
				return
			}

			msg, err := generator.Generate(0)
			if err != nil {
				log.Printf("failed to generate: %+v", err)
				return
			}

			_, err = event.Message.Reply(context.Background(), s, msg)
			if err != nil {
				log.Printf("failed to send: %+v", err)
				return
			}
		})
}
