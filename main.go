package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/greymd/mamadm/generator"
)

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_TOKEN"),
		Intents:  disgord.IntentGuildMessages,
	})

	defer client.Gateway().StayConnectedUntilInterrupted()

	u, err := client.BotAuthorizeURL(disgord.PermissionUseSlashCommands, []string{
		"bot",
		"applications.commands",
	})
	if err != nil {
		panic(err)
	}
	log.Println(u)

	client.Gateway().BotReady(func() {
		appID, err := strconv.Atoi(os.Getenv("DISCORD_APPLICATION_ID"))
		if err != nil {
			log.Fatalf("failed to convert string to int: %+v", err)
		}

		if err := client.ApplicationCommand(snowflake.NewSnowflake(uint64(appID))).
			Global().
			Create(&disgord.CreateApplicationCommand{
				Name:        "mama",
				Description: "Send Mama DM",
			}); err != nil {
			log.Fatalf("failed to create application command: %+v", err)
		}
	})

	client.Gateway().
		InteractionCreate(func(s disgord.Session, event *disgord.InteractionCreate) {
			if event.Data.Name != "mama" {
				return
			}

			msg, err := generator.Generate(0)
			if err != nil {
				log.Printf("failed to generate: %+v", err)
				return
			}

			if err = s.SendInteractionResponse(context.Background(), event, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: msg,
				},
			}); err != nil {
				log.Printf("failed to send response: %+v", err)
				return
			}
		})
}
