package storage

import (
	"FeedRocket/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"log"
	"time"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {

	}
	defer closeConnection(conn)

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, "SELECT * FROM sources"); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(source dbSource, _ int) model.Source {
		return model.Source(source)
	}), nil
}

func (s *SourcePostgresStorage) SourceById(ctx context.Context, id int64) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	var source dbSource
	if err := conn.GetContext(ctx, &source, "SELECT * FROM sources WHERE id = $1", id); err != nil {
		return nil, err
	}

	return (*model.Source)(&source), nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source *model.Source) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer closeConnection(conn)

	var id int64
	row := conn.QueryRowxContext(
		ctx,
		"INSERT INTO sources(name, feed_url, created_at) VALUES($1, $2, $3) RETURNING id",
		source.Name,
		source.FeedUrl,
		source.CreatedAt)

	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer closeConnection(conn)
	if _, err := conn.ExecContext(ctx, "DELETE FROM sources WHERE id = $1", id); err != nil {
		return err
	}
}

func closeConnection(conn *sqlx.Conn) {
	err := conn.Close()
	if err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedUrl   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}
