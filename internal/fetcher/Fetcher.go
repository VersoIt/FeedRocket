package fetcher

import (
	"FeedRocket/internal/model"
	"FeedRocket/internal/source"
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

type ArticleStorage interface {
	Store(ctx context.Context, article model.Article) error
}

type SourceProvider interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articles ArticleStorage
	sources  SourceProvider

	fetchInterval  time.Duration
	filterKeywords []string
}

func New(
	articleStorage ArticleStorage,
	sourceProvider SourceProvider,
	fetchInterval time.Duration,
	filterKeywords []string,
) *Fetcher {
	return &Fetcher{
		articles:       articleStorage,
		sources:        sourceProvider,
		fetchInterval:  fetchInterval,
		filterKeywords: filterKeywords,
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	log.Println("starting fetching...")
	for _, src := range sources {
		wg.Add(1)
		rssSource := source.NewRSSSourceFromModel(src)
		go func(source Source) {
			defer wg.Done()
			items, err := source.Fetch(ctx)
			if err != nil {
				log.Printf("[ERROR] failed to fetch items: %v", err)
				return
			}
			if err := f.processItems(ctx, source, items); err != nil {
				log.Printf("[ERROR] failed to process items: %v", err)
				return
			}
		}(rssSource)
	}

	wg.Wait()
	return nil
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()
		if f.itemShouldBeSkipped(item) {
			continue
		}
		if err := f.articles.Store(ctx, model.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fetcher) itemShouldBeSkipped(item model.Item) bool {
	keywordsMap := sliceToMap(f.filterKeywords)

	for _, keyword := range f.filterKeywords {
		titleContainsKeyword := strings.Contains(strings.ToLower(item.Title), keyword)
		if _, containsKeyword := keywordsMap[keyword]; containsKeyword || titleContainsKeyword {
			return true
		}
	}

	return false
}

func sliceToMap[T comparable](slice []T) map[T]struct{} {
	result := make(map[T]struct{}, len(slice))
	for _, value := range slice {
		result[value] = struct{}{}
	}
	return result
}
