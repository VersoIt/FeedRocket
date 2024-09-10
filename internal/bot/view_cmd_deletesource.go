package bot

import (
	"FeedRocket/internal/bot/bottotal"
	"FeedRocket/internal/botkit"
	"FeedRocket/internal/botkit/markup"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeleteStorage interface {
	Delete(ctx context.Context, id int64) error
	bottotal.SourceById
}

func ViewCmdDelete(storage DeleteStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

		type DeleteSourceArgs struct {
			Id int64 `json:"id"`
		}

		jsonArgs := update.Message.CommandArguments()
		args, err := botkit.ParseJson[DeleteSourceArgs](jsonArgs)
		if err != nil {
			_ = bottotal.SendMessage(update, bot, "Incorrect arguments")
			return err
		}

		sourceToRemove, err := storage.SourceById(ctx, args.Id)
		if err != nil {
			_ = bottotal.SendMessage(update, bot, fmt.Sprintf("Can't find source with ID: %d", args.Id))
			return err
		}

		err = storage.Delete(ctx, args.Id)
		if err != nil {
			_ = bottotal.SendMessage(update, bot, "Can't delete source")
			return err
		}

		return bottotal.SendMessage(update, bot, markup.EscapeForMarkdown(fmt.Sprintf("Source %s succeffully deleted", sourceToRemove.Name)))
	}
}
