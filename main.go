package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v5"
	mamagen "github.com/greymd/mamadm/generator"
	ojigen "github.com/greymd/ojichat/generator"
	"golang.org/x/xerrors"
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
		applicationCommands := []*disgord.CreateApplicationCommand{
			{
				Name:        "mama",
				Description: "Send Mama DM",
			},
			{
				Name:        "oji",
				Description: "Send ojichat",
				Options: []*disgord.ApplicationCommandOption{
					{
						Type:        disgord.OptionTypeString,
						Name:        "name",
						Description: "Target name",
					},
				},
			},
		}

		appID, err := strconv.Atoi(os.Getenv("DISCORD_APPLICATION_ID"))
		if err != nil {
			log.Fatalf("failed to convert string to int: %+v", err)
		}

		functions := client.ApplicationCommand(snowflake.NewSnowflake(uint64(appID))).Global()

		for _, applicationCommand := range applicationCommands {
			if err := functions.Create(applicationCommand); err != nil {
				log.Fatalf("failed to create application command: %+v", err)
			}
		}
	})

	client.Gateway().
		InteractionCreate(func(s disgord.Session, event *disgord.InteractionCreate) {
			ctx := context.Background()

			go interactionCreate(ctx, s, event)
		})
}

func interactionCreate(ctx context.Context, s disgord.Session, event *disgord.InteractionCreate) {
	var (
		msg string
		err error
	)

	switch event.Data.Name {
	case "mama":
		msg, err = mama(ctx)

	case "oji":
		msg, err = oji(ctx, event.Data.Options...)
	}

	if err != nil {
		log.Printf("failed to execute mama: %+v", err)
		msg = "failed to execute"
	}

	if err = s.SendInteractionResponse(ctx, event, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: msg,
		},
	}); err != nil {
		log.Printf("failed to send response: %+v", err)
	}
}

func mama(ctx context.Context) (string, error) {
	msg, err := mamagen.Generate(0)
	if err != nil {
		return "", xerrors.Errorf("failed to generate: %w", err)
	}

	return msg, nil
}

func oji(ctx context.Context, opts ...*disgord.ApplicationCommandDataOption) (string, error) {
	cfg := ojigen.Config{
		EmojiNum: 4,
	}

	if len(opts) > 0 {
		v, ok := opts[0].Value.(string)
		if !ok {
			return "", xerrors.New("failed to cast to string")
		}

		cfg.TargetName = v
	}

	msg, err := ojigen.Start(cfg)
	if err != nil {
		return "", xerrors.Errorf("failed to generate: %w", err)
	}

	return msg, nil
}
