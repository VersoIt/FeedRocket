package bot

import (
	"FeedRocket/internal/bot/bottotal"
	"FeedRocket/internal/botkit"
	"FeedRocket/internal/model"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceUpdater interface {
	UpdateSource(ctx context.Context, source *model.Source) error
}

func ViewCmdEditSource(sourceUpdater SourceUpdater) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		type EditSourceArgs struct {
			Id      int64  `json:"id"`
			Name    string `json:"name"`
			FeedUrl string `json:"url"`
		}

		commandArgs := update.Message.CommandArguments()
		args, err := botkit.ParseJson[EditSourceArgs](commandArgs)
		if err != nil {
			_ = bottotal.SendMessage(update, bot, "Invalid arguments")
			return err
		}

		err = sourceUpdater.UpdateSource(ctx, &model.Source{
			ID:      args.Id,
			Name:    args.Name,
			FeedUrl: args.FeedUrl,
		})

		if err != nil {
			return err
		}

		return bottotal.SendMessage(update, bot, "Source successfully edited")
	}
}
