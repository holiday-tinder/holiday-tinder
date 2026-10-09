package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// connect opens a pool and waits (bounded) until the database answers.
func connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(90 * time.Second)
	for {
		err = pool.Ping(ctx)
		if err == nil {
			return pool, nil
		}
		if time.Now().After(deadline) {
			pool.Close()
			return nil, err
		}
		log.Printf("waiting for database: %v", err)
		select {
		case <-ctx.Done():
			pool.Close()
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}

var schema = []string{
	`CREATE TABLE IF NOT EXISTS venues (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT NOT NULL CHECK (category IN ('restaurant', 'bar', 'club')),
		city TEXT NOT NULL,
		neighborhood TEXT NOT NULL,
		cuisine TEXT NOT NULL,
		rating NUMERIC(2,1) NOT NULL,
		review_count INT NOT NULL,
		price_level INT NOT NULL,
		emoji TEXT NOT NULL,
		description TEXT NOT NULL,
		highlights TEXT[] NOT NULL DEFAULT '{}',
		UNIQUE (name, city)
	)`,
	`CREATE TABLE IF NOT EXISTS sessions (
		code TEXT PRIMARY KEY,
		city TEXT NOT NULL,
		categories TEXT[] NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`,
	`CREATE TABLE IF NOT EXISTS participants (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		session_code TEXT NOT NULL REFERENCES sessions(code) ON DELETE CASCADE,
		name TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`,
	`CREATE INDEX IF NOT EXISTS participants_session_idx ON participants (session_code)`,
	`CREATE TABLE IF NOT EXISTS swipes (
		participant_id UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
		venue_id INT NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
		liked BOOLEAN NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		PRIMARY KEY (participant_id, venue_id)
	)`,
}

func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	for _, stmt := range schema {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("schema: %w", err)
		}
	}
	return seed(ctx, pool)
}
