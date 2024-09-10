package main

import (
	"FeedRocket/internal/bot"
	"FeedRocket/internal/bot/middleware"
	"FeedRocket/internal/botkit"
	"FeedRocket/internal/config"
	"FeedRocket/internal/fetcher"
	"FeedRocket/internal/notifier"
	"FeedRocket/internal/storage"
	"FeedRocket/internal/summary"
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	botApi, err := tgbotapi.NewBotAPI(config.Get().BotConfig.Token)
	if err != nil {
		log.Printf("Failed to connect to telegram bot: %v", err)
		return
	}

	db, err := sqlx.Connect(config.Get().DbConfig.Driver, getDSN(config.Get().DbConfig.Driver, "feed_rocket_bot", "postgres", "postgres", config.Get().DbConfig.Address, "disable", config.Get().DbConfig.Port))
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}(db)

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetch          = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().ServerConfig.FetchInterval,
			config.Get().ServerConfig.FilterKeywords,
		)

		notify = notifier.New(
			articleStorage,
			summary.NewOpenAISummarizer(config.Get().OpenAiConfig.Key, config.Get().OpenAiConfig.Prompt),
			botApi,
			config.Get().ServerConfig.NotificationInterval,
			2*config.Get().ServerConfig.FetchInterval,
			config.Get().BotConfig.ChannelId)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.New(botApi)
	newsBot.RegisterCmdView("editsource", bot.ViewCmdEditSource(sourceStorage))
	newsBot.RegisterCmdView("help", bot.ViewCmdHelp())
	newsBot.RegisterCmdView("deletesource", middleware.OnlyAdmin(config.Get().BotConfig.ChannelId, bot.ViewCmdDelete(sourceStorage)))
	newsBot.RegisterCmdView("start", middleware.OnlyAdmin(config.Get().BotConfig.ChannelId, bot.ViewCmdStart()))
	newsBot.RegisterCmdView("addsource", middleware.OnlyAdmin(config.Get().BotConfig.ChannelId, bot.ViewCmdAddSource(sourceStorage)))
	newsBot.RegisterCmdView("sourcelist", middleware.OnlyAdmin(config.Get().BotConfig.ChannelId, bot.ViewCmdListSources(sourceStorage)))

	go func(ctx context.Context) {
		if err := fetch.Fetch(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("failed to start fetcher: %v", err)
				return
			}
		}
	}(ctx)

	go func() {
		if err := notify.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start notifier: %v", err)
				return
			}

			log.Println("notifier stopped")
		}
	}()

	if err := newsBot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to start news bot: %v", err)
		}
		log.Println("bot stopped")
	}
}

func getDSN(driver, dbName, login, password, address, sslMode string, port int) string {
	return fmt.Sprintf("%s://%s:%s@%s:%v/%s?sslmode=%s", driver, login, password, address, port, dbName, sslMode)
}
