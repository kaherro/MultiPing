package checker

import (
	"context"
	"net/http"
	"time"

	"github.com/kaherro/MultiPing/internal/domain"
)

type HTTPChecker struct {
	client *http.Client
}

func (c *HTTPChecker) Check(ctx context.Context, site domain.Site) domain.CheckResult {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(site.TimeoutMS) * time.Millisecond)
	defer cancel()

	res := domain.CheckResult{SiteID: site.ID, CheckedAT: time.Now()}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, site.URL, nil)
	if err != nil {
		res.Error = err
		return res
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	res.LatencyMS = int(time.Since(start).Milliseconds())
	if err != nil {
		res.Error = err
		return res
	}
	defer resp.Body.Close()

	res.HTTPCode = resp.StatusCode
	res.OK = (resp.StatusCode >= 200 && resp.StatusCode < 400)
	return res
}