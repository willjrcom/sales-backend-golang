package queuerepositorylocal

import (
   "context"
   "sync"

   "github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type QueueRepositoryLocal struct {
   mu     sync.RWMutex
   queues map[string]*model.OrderQueue
}

func NewQueueRepositoryLocal() model.QueueRepository {
   return &QueueRepositoryLocal{queues: make(map[string]*model.OrderQueue)}
}

func (r *QueueRepositoryLocal) CreateQueue(ctx context.Context, p *model.OrderQueue) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.queues[p.ID.String()] = p
   return nil
}

func (r *QueueRepositoryLocal) UpdateQueue(ctx context.Context, p *model.OrderQueue) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.queues[p.ID.String()] = p
   return nil
}

func (r *QueueRepositoryLocal) DeleteQueue(ctx context.Context, id string) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   delete(r.queues, id)
   return nil
}

func (r *QueueRepositoryLocal) GetQueueById(ctx context.Context, id string) (*model.OrderQueue, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   if q, ok := r.queues[id]; ok {
       return q, nil
   }
   return nil, nil
}

func (r *QueueRepositoryLocal) GetOpenedQueueByGroupItemId(ctx context.Context, id string) (*model.OrderQueue, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   for _, q := range r.queues {
       if q.GroupItemID.String() == id && q.LeftAt == nil {
           return q, nil
       }
   }
   return nil, nil
}

func (r *QueueRepositoryLocal) GetQueuesByGroupItemId(ctx context.Context, id string) ([]model.OrderQueue, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := []model.OrderQueue{}
   for _, q := range r.queues {
       if q.GroupItemID.String() == id {
           out = append(out, *q)
       }
   }
   return out, nil
}

func (r *QueueRepositoryLocal) GetAllQueues(ctx context.Context) ([]model.OrderQueue, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := make([]model.OrderQueue, 0, len(r.queues))
   for _, q := range r.queues {
       out = append(out, *q)
   }
   return out, nil
}
