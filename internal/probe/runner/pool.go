package runner

import (
	"context"
	"sync"

	"github.com/kaherro/MultiPing/internal/checker"
	"github.com/kaherro/MultiPing/internal/domain"
)

func RunPool(ctx context.Context, chk checker.Checker, workers int, 
	jobs <-chan domain.Site, results chan<- domain.CheckResult) {
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <- ctx.Done():
					return
				case site, ok := <- jobs:
					if !ok {
						return
					}
					res := chk.Check(ctx, site)
					select {
					case <- ctx.Done():
						return
					case results <- res:
					}
				}
			}
		}()
	}
	wg.Wait()
	close(results)
}
