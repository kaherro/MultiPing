package checker

import (
	"context"

	"github.com/kaherro/MultiPing/internal/domain"
)

type Checker interface {
	Check(ctx context.Context, site domain.Site) domain.CheckResult
}
