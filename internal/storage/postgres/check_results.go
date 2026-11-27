package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaherro/MultiPing/internal/domain"
)

type CheckResultRepo struct {
	pool *pgxpool.Pool
}

func NewCheckResultRepo(pool *pgxpool.Pool) *CheckResultRepo {
	return &CheckResultRepo{pool: pool}
}

func (r *CheckResultRepo) Create(ctx context.Context, c *domain.CheckResult) error {
	var errText any
	if c.Error != nil {
		errText = c.Error.Error()
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO check_results (
			site_id, ok, http_code, latency_ms, error, checked_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		c.SiteID, c.OK, c.HTTPCode, c.LatencyMS, errText, c.CheckedAT,
	)

	return err
}

func (r *CheckResultRepo) List(ctx context.Context) ([]domain.CheckResult, error) {
	rows, err := r.pool.Query(ctx, 
		`SELECT site_id, ok, http_code, latency_ms, error, checked_at
		FROM check_results ORDER BY checked_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]domain.CheckResult, 0)

	for rows.Next() {
		var c domain.CheckResult
		var errText pgtype.Text

		err := rows.Scan(
			&c.SiteID,
			&c.OK,
			&c.HTTPCode,
			&c.LatencyMS,
			&errText,
			&c.CheckedAT,
		)
		if err != nil {
			return nil, err
		}

		if errText.Valid {
			c.Error = errors.New(errText.String)
		}

		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *CheckResultRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM check_results WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil 
}