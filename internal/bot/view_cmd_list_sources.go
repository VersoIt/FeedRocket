package bot

import (
	"FeedRocket/internal/bot/bottotal"
	"FeedRocket/internal/botkit"
	"FeedRocket/internal/botkit/markup"
	"FeedRocket/internal/model"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"strings"
)

type SourceLister interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

func ViewCmdListSources(lister SourceLister) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sources, err := lister.Sources(ctx)
		if err != nil {
			return err
		}

		var (
			sourceInfos = lo.Map(sources, func(source model.Source, _ int) string {
				return formatSource(&source)
			})
			msgText = fmt.Sprintf("Source list \\(total %d\\): \n\n%s", len(sources), strings.Join(sourceInfos, "\n\n"))
		)

		return bottotal.SendMessage(update, bot, msgText)
	}
}

func formatSource(source *model.Source) string {
	return fmt.Sprintf("*%s*\nID: `%d`\nURL: %s", markup.EscapeForMarkdown(source.Name), source.ID, markup.EscapeForMarkdown(source.FeedUrl))
}
