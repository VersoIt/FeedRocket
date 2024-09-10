package bot

import (
	"FeedRocket/internal/bot/bottotal"
	"FeedRocket/internal/botkit"
	"FeedRocket/internal/model"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceStorage interface {
	Add(ctx context.Context, source *model.Source) (int64, error)
}

func ViewCmdAddSource(storage SourceStorage) botkit.ViewFunc {

	type addSourceArgs struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJson[addSourceArgs](update.Message.CommandArguments())
		if err != nil {
			_ = bottotal.SendMessage(update, bot, "Incorrect arguments")
			return err
		}
		source := model.Source{
			Name:    args.Name,
			FeedUrl: args.Url,
		}

		sourceId, err := storage.Add(ctx, &source)
		if err != nil {
			return err
		}

		return bottotal.SendMessage(update, bot, fmt.Sprintf("Source added with ID: `%v`\\.Use this ID to manage source\\.", sourceId))
	}
}
