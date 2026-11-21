package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaherro/MultiPing/internal/domain"
)

type SiteRepo struct {
	pool *pgxpool.Pool
}

func NewSiteRepo(pool *pgxpool.Pool) *SiteRepo {
	return &SiteRepo{pool: pool}
}

func (r *SiteRepo) Create(ctx context.Context, s *domain.Site) error {
	err := r.pool.QueryRow(ctx, 
		`INSERT INTO sites (url, interval_sec, timeout_ms) 
		VALUES ($1, $2, $3) 
		RETURNING id, enabled, created_at`, 
		s.URL, s.IntervalSec, s.TimeoutMS,
		).Scan(&s.ID, &s.Enabled, &s.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrAlreadyExists
	}
	return err
}

func (r *SiteRepo) List(ctx context.Context) ([]domain.Site, error) {
	rows, err := r.pool.Query(ctx, 
	`SELECT id, url, interval_sec, timeout_ms, enabled, created_at
	FROM sites ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close() 

	sites := []domain.Site{}
	for rows.Next() {
		var s domain.Site
		if err := rows.Scan(&s.ID, &s.URL, &s.IntervalSec, &s.TimeoutMS, 
			&s.Enabled, &s.CreatedAt); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

func (r *SiteRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM sites WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil 
}