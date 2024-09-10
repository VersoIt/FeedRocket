package bot

import (
	"FeedRocket/internal/bot/bottotal"
	"FeedRocket/internal/botkit"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdHelp() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		return bottotal.SendMessage(update, bot, "Command list:\n/start\n/deletesource\n/addsource\n/sourcelist\n/editsource")
	}
}
