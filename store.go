package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Store struct{ pool *pgxpool.Pool }

type Venue struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	City         string   `json:"city"`
	Neighborhood string   `json:"neighborhood"`
	Cuisine      string   `json:"cuisine"`
	Rating       float64  `json:"rating"`
	ReviewCount  int      `json:"reviewCount"`
	PriceLevel   int      `json:"priceLevel"`
	Emoji        string   `json:"emoji"`
	Description  string   `json:"description"`
	Highlights   []string `json:"highlights"`
}

const venueCols = `v.id, v.name, v.category, v.city, v.neighborhood, v.cuisine, v.rating::float8,
	v.review_count, v.price_level, v.emoji, v.description, v.highlights`

func scanVenue(row pgx.Row, extra ...any) (Venue, error) {
	var v Venue
	dest := append([]any{&v.ID, &v.Name, &v.Category, &v.City, &v.Neighborhood, &v.Cuisine, &v.Rating,
		&v.ReviewCount, &v.PriceLevel, &v.Emoji, &v.Description, &v.Highlights}, extra...)
	err := row.Scan(dest...)
	if v.Highlights == nil {
		v.Highlights = []string{}
	}
	return v, err
}

func (s *Store) Cities(ctx context.Context) ([]string, error) {
	rows, err := s.pool.Query(ctx, `SELECT DISTINCT city FROM venues ORDER BY city`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (s *Store) CityExists(ctx context.Context, city string) (bool, error) {
	var ok bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM venues WHERE city = $1)`, city).Scan(&ok)
	return ok, err
}

// CreateSession creates a session with a fresh code and its first participant.
func (s *Store) CreateSession(ctx context.Context, city string, cats []string, name string) (string, string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback(ctx)

	code := ""
	for attempt := 0; attempt < 5 && code == ""; attempt++ {
		candidate, err := newCode()
		if err != nil {
			return "", "", err
		}
		tag, err := tx.Exec(ctx, `INSERT INTO sessions (code, city, categories) VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING`, candidate, city, cats)
		if err != nil {
			return "", "", err
		}
		if tag.RowsAffected() == 1 {
			code = candidate
		}
	}
	if code == "" {
		return "", "", errors.New("could not allocate a session code")
	}

	var pid string
	err = tx.QueryRow(ctx, `INSERT INTO participants (session_code, name) VALUES ($1, $2) RETURNING id::text`,
		code, name).Scan(&pid)
	if err != nil {
		return "", "", err
	}
	return code, pid, tx.Commit(ctx)
}

func (s *Store) Join(ctx context.Context, code, name string) (string, error) {
	var pid string
	err := s.pool.QueryRow(ctx, `INSERT INTO participants (session_code, name)
		SELECT code, $2 FROM sessions WHERE code = $1 RETURNING id::text`, code, name).Scan(&pid)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	return pid, err
}

func (s *Store) IsMember(ctx context.Context, code, pid string) (bool, error) {
	var ok bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM participants
		WHERE id = $1::uuid AND session_code = $2)`, pid, code).Scan(&ok)
	return ok, err
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
