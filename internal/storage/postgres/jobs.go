package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaherro/MultiPing/internal/domain"
)

type JobsRepo struct {
	pool *pgxpool.Pool
}

func NewJobsRepo(pool *pgxpool.Pool) *JobsRepo {
	return &JobsRepo{pool: pool}
}

func (r *JobsRepo) Create(ctx context.Context, j *domain.Job) error {
    return r.pool.QueryRow(ctx, `
        INSERT INTO jobs (site_id, agent_id, status)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`,
        j.SiteID, j.AgentID, domain.JobStatusPending,
    ).Scan(&j.ID, &j.CreatedAt)
}

func (r *JobsRepo) List(ctx context.Context) ([]domain.Job, error) {
	rows, err := r.pool.Query(ctx, 
		`SELECT id, site_id, agent_id, status, created_at, assigned_at, deadline_at, finished_at,
		FROM jobs ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []domain.Job{}
	for rows.Next() {
		var j domain.Job
		if err := rows.Scan(&j.ID, &j.SiteID, &j.AgentID, &j.Status, 
			&j.CreatedAt, &j.AssignedAt, &j.DeadlineAt, &j.FinishedAt); 
			err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

func (r *JobsRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM jobs WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil 
}