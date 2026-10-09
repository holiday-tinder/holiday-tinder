package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Result struct {
	Venue   Venue    `json:"venue"`
	Likes   int      `json:"likes"`
	Likers  []string `json:"likers"`
	Matched bool     `json:"matched"`
}

type Summary struct {
	City         string   `json:"city"`
	Categories   []string `json:"categories"`
	Participants []string `json:"participants"`
	Results      []Result `json:"results"`
}

// Summary returns the session, who is in it, and every liked venue.
func (s *Store) Summary(ctx context.Context, code string) (Summary, error) {
	var sum Summary
	err := s.pool.QueryRow(ctx, `SELECT city, categories FROM sessions WHERE code = $1`, code).
		Scan(&sum.City, &sum.Categories)
	if errors.Is(err, pgx.ErrNoRows) {
		return sum, ErrNotFound
	}
	if err != nil {
		return sum, err
	}

	rows, err := s.pool.Query(ctx, `SELECT name FROM participants WHERE session_code = $1 ORDER BY created_at`, code)
	if err != nil {
		return sum, err
	}
	sum.Participants, err = pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return sum, err
	}
	if sum.Participants == nil {
		sum.Participants = []string{}
	}

	rows, err = s.pool.Query(ctx, `SELECT `+venueCols+`, count(*)::int AS likes,
			array_agg(p.name ORDER BY p.created_at) AS likers
		FROM swipes w
		JOIN participants p ON p.id = w.participant_id
		JOIN venues v ON v.id = w.venue_id
		WHERE p.session_code = $1 AND w.liked
		GROUP BY v.id
		ORDER BY likes DESC, v.rating DESC, v.review_count DESC`, code)
	if err != nil {
		return sum, err
	}
	defer rows.Close()

	sum.Results = []Result{}
	total := len(sum.Participants)
	for rows.Next() {
		var r Result
		r.Venue, err = scanVenue(rows, &r.Likes, &r.Likers)
		if err != nil {
			return sum, err
		}
		r.Matched = total >= 2 && r.Likes == total
		sum.Results = append(sum.Results, r)
	}
	return sum, rows.Err()
}
