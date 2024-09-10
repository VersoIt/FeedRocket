package storage

import (
	"FeedRocket/internal/model"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer closeConnection(conn)
	if _, err := s.db.ExecContext(ctx, `INSERT INTO articles(source_id, title, link, summary, published_at) VALUES($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`, article.SourceID, article.Title, article.Link, article.Summary, time.Now()); err != nil {
		return err
	}

	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	var articles []dbArticle
	if err := conn.SelectContext(ctx, &articles, `SELECT * FROM articles 
         WHERE posted_at IS NULL
           AND published_at >= $1::timestamp
         ORDER BY published_at DESC
         LIMIT $2`, since.UTC().Format(time.RFC3339), limit); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			ID:          article.ID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary,
			PostedAt:    article.PostedAt.Time,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer closeConnection(conn)
	if _, err := conn.ExecContext(ctx, `UPDATE articles SET posted_at = $1::timestamp WHERE id = $2`, time.Now().UTC().Format(time.RFC3339), id); err != nil {
		return err
	}

	return nil
}

type dbArticle struct {
	ID          int64        `db:"id"`
	SourceId    int64        `db:"source_id"`
	Title       string       `db:"title"`
	Link        string       `db:"link"`
	Summary     string       `db:"summary"`
	PostedAt    sql.NullTime `db:"posted_at"`
	PublishedAt time.Time    `db:"published_at"`
	CreatedAt   time.Time    `db:"created_at"`
}
