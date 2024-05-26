package workers

import (
	"context"

	"github.com/romanyakovlev/gophermart/internal/service"
)

type WorkerPool struct {
	workers             []*AccrualWorker
	accrualRequestsChan chan AccrualRequest
}

func NewWorkerPool(numWorkers int, o service.OrdersService, u service.UserService, a service.AccrualService) *WorkerPool {
	pool := &WorkerPool{
		accrualRequestsChan: make(chan AccrualRequest, 100),
	}
	for i := 0; i < numWorkers; i++ {
		worker := InitAccrualWorker(o, u, &a, pool.accrualRequestsChan)
		pool.workers = append(pool.workers, worker)
	}
	return pool
}

func (p *WorkerPool) StartAll(ctx context.Context) {
	for i, worker := range p.workers {
		go worker.StartAccrualWorker(ctx, i+1)
	}
}

func (p *WorkerPool) SendAccrualRequest(req AccrualRequest) error {
	p.accrualRequestsChan <- req
	return nil
}
