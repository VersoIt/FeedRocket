package middleware

import (
	"FeedRocket/internal/botkit"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OnlyAdmin(channelId int64, next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: channelId,
			},
		})
		if err != nil {
			return err
		}

		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {
				return next(ctx, bot, update)
			}
		}

		_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You don't have permission"))
		if err != nil {
			return err
		}

		return nil
	}
}
