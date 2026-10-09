package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// Deck returns the venues this participant has not swiped yet.
func (s *Store) Deck(ctx context.Context, code, pid string) ([]Venue, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+venueCols+`
		FROM venues v JOIN sessions s ON s.code = $1
		WHERE v.city = s.city AND v.category = ANY (s.categories)
		  AND NOT EXISTS (SELECT 1 FROM swipes w WHERE w.participant_id = $2::uuid AND w.venue_id = v.id)
		ORDER BY v.rating DESC, v.review_count DESC`, code, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	venues := []Venue{}
	for rows.Next() {
		v, err := scanVenue(rows)
		if err != nil {
			return nil, err
		}
		venues = append(venues, v)
	}
	return venues, rows.Err()
}

// Swipe records a decision and returns the venue when it just became a
// match: liked by every participant of a session of at least two.
func (s *Store) Swipe(ctx context.Context, code, pid string, venueID int, liked bool) (*Venue, error) {
	_, err := s.pool.Exec(ctx, `INSERT INTO swipes (participant_id, venue_id, liked) VALUES ($1::uuid, $2, $3)
		ON CONFLICT (participant_id, venue_id) DO UPDATE SET liked = EXCLUDED.liked, created_at = now()`,
		pid, venueID, liked)
	if isForeignKeyViolation(err) {
		return nil, ErrNotFound
	}
	if err != nil || !liked {
		return nil, err
	}

	row := s.pool.QueryRow(ctx, `SELECT `+venueCols+` FROM venues v
		WHERE v.id = $2
		  AND (SELECT count(*) FROM participants p WHERE p.session_code = $1) >= 2
		  AND NOT EXISTS (
		    SELECT 1 FROM participants p WHERE p.session_code = $1
		      AND NOT EXISTS (SELECT 1 FROM swipes w
		        WHERE w.participant_id = p.id AND w.venue_id = v.id AND w.liked))`, code, venueID)
	v, err := scanVenue(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}
