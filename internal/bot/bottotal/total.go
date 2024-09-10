package bottotal

import (
	"FeedRocket/internal/model"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceById interface {
	SourceById(ctx context.Context, id int64) (*model.Source, error)
}

func SendMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, message string) error {
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	reply.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := bot.Send(reply)
	return err
}
