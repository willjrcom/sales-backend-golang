package orderprocessrepositorylocal

import (
   "context"
   "sync"

   "github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderProcessRepositoryLocal struct {
   mu        sync.RWMutex
   processes map[string]*model.OrderProcess
}

func NewOrderProcessRepositoryLocal() model.OrderProcessRepository {
   return &OrderProcessRepositoryLocal{processes: make(map[string]*model.OrderProcess)}
}

func (r *OrderProcessRepositoryLocal) CreateProcess(ctx context.Context, p *model.OrderProcess) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.processes[p.ID.String()] = p
   return nil
}

func (r *OrderProcessRepositoryLocal) UpdateProcess(ctx context.Context, p *model.OrderProcess) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.processes[p.ID.String()] = p
   return nil
}

func (r *OrderProcessRepositoryLocal) DeleteProcess(ctx context.Context, id string) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   delete(r.processes, id)
   return nil
}

func (r *OrderProcessRepositoryLocal) GetProcessById(ctx context.Context, id string) (*model.OrderProcess, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   if p, ok := r.processes[id]; ok {
       return p, nil
   }
   return nil, nil
}

func (r *OrderProcessRepositoryLocal) GetAllProcesses(ctx context.Context) ([]model.OrderProcess, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := make([]model.OrderProcess, 0, len(r.processes))
   for _, p := range r.processes {
       out = append(out, *p)
   }
   return out, nil
}

func (r *OrderProcessRepositoryLocal) GetProcessesByProcessRuleID(ctx context.Context, id string) ([]model.OrderProcess, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := []model.OrderProcess{}
   for _, p := range r.processes {
       if p.ProcessRuleID.String() == id {
           out = append(out, *p)
       }
   }
   return out, nil
}

func (r *OrderProcessRepositoryLocal) GetProcessesByProductID(ctx context.Context, id string) ([]model.OrderProcess, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := []model.OrderProcess{}
   for _, p := range r.processes {
       for _, prod := range p.Products {
           if prod.ID.String() == id {
               out = append(out, *p)
               break
           }
       }
   }
   return out, nil
}

func (r *OrderProcessRepositoryLocal) GetProcessesByGroupItemID(ctx context.Context, id string) ([]model.OrderProcess, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := []model.OrderProcess{}
   for _, p := range r.processes {
       if p.GroupItemID.String() == id {
           out = append(out, *p)
       }
   }
   return out, nil
}
