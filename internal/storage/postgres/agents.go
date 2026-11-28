package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaherro/MultiPing/internal/domain"
)

type AgentsRepo struct {
	pool *pgxpool.Pool
}

func NewAgentsRepo(pool *pgxpool.Pool) *AgentsRepo {
	return &AgentsRepo{pool: pool}
}

func (r *AgentsRepo) Create(ctx context.Context, a *domain.Agent) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO agents (
			name, token_hash, enabled, last_seen_at
		)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		a.Name, a.TokenHash, a.Enabled, a.LastSeenAt,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *AgentsRepo) List(ctx context.Context) ([]domain.Agent, error) {
	rows, err := r.pool.Query(ctx, 
		`SELECT id, name, token_hash, enabled, last_seen_at, created_at
		FROM agents ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := []domain.Agent{}
	for rows.Next() {
		var a domain.Agent
		if err := rows.Scan(&a.ID, &a.Name, &a.TokenHash, &a.Enabled, 
			&a.LastSeenAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}
	return agents, rows.Err()
}

func (r *AgentsRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM agents WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil 
}