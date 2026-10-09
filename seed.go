package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// seedVenue is a TripAdvisor-style sample listing.
type seedVenue struct {
	City, Neighborhood, Name, Category, Cuisine, Emoji string
	Rating                                             float64
	Reviews, Price                                     int
	Description                                        string
	Highlights                                         []string
}

func venue(city, hood, name, cat, cuisine, emoji string, rating float64, reviews, price int, desc string, highlights ...string) seedVenue {
	return seedVenue{city, hood, name, cat, cuisine, emoji, rating, reviews, price, desc, highlights}
}

// seed inserts the sample venues; existing rows are left untouched.
func seed(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	all := append(append(lisbon(), barcelona()...), amsterdam()...)
	for _, v := range all {
		_, err := tx.Exec(ctx, `INSERT INTO venues
			(name, category, city, neighborhood, cuisine, rating, review_count, price_level, emoji, description, highlights)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (name, city) DO NOTHING`,
			v.Name, v.Category, v.City, v.Neighborhood, v.Cuisine, v.Rating, v.Reviews, v.Price, v.Emoji, v.Description, v.Highlights)
		if err != nil {
			return fmt.Errorf("seed %s: %w", v.Name, err)
		}
	}
	return tx.Commit(ctx)
}
