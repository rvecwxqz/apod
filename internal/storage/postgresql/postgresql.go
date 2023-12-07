package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rvecwxqz/apod/internal/core"
	"time"
)

type PostgreSQL struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, DSN string) (*PostgreSQL, error) {
	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to create pool: %w", err)
	}

	_, err = pool.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS apod_info(
    			id SERIAL PRIMARY KEY,
    			title VARCHAR,
    			explanation VARCHAR, 
    			date DATE UNIQUE
			 )`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create table: %w", err)
	}

	return &PostgreSQL{pool}, nil
}

func (s PostgreSQL) SaveInfo(
	ctx context.Context,
	title string,
	explanation string,
	date core.Date,
) error {
	_, err := s.pool.Exec(
		ctx,
		`INSERT INTO apod_info(title, explanation, date) VALUES($1, $2, $3)`,
		title, explanation, date,
	)
	if err != nil {
		return fmt.Errorf("insert info error: %w", err)
	}
	return nil
}

func (s PostgreSQL) GetInfo(ctx context.Context, d core.Date) (core.APODInfo, error) {
	var info core.APODInfo
	var t time.Time
	err := s.pool.QueryRow(
		ctx,
		`SELECT title, explanation, date FROM apod_info WHERE DATE = $1`, d,
	).Scan(&info.Title, &info.Explanation, &t)
	info.Date = core.Date(t)

	if err != nil {
		return core.APODInfo{}, fmt.Errorf("error selecting row: %w", err)
	}
	return info, nil
}

func (s PostgreSQL) GetAllInfo(ctx context.Context) ([]core.APODInfo, error) {
	out := make([]core.APODInfo, 0, 10)
	rows, err := s.pool.Query(
		ctx,
		`SELECT title, explanation, date FROM apod_info`,
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows: %w", err)
	}

	for rows.Next() {
		var info core.APODInfo
		var t time.Time

		err = rows.Scan(&info.Title, &info.Explanation, &t)
		info.Date = core.Date(t)

		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}
		out = append(out, info)
	}

	return out, nil
}

func (s PostgreSQL) Stop() {
	s.pool.Close()
}
