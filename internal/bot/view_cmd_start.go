package bot

import (
	"FeedRocket/internal/botkit"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Started!")); err != nil {
			return err
		}

		return nil
	}
}
