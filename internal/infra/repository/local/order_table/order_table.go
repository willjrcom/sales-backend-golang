package ordertablerepositorylocal

import (
   "context"
   "sync"

   "github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderTableRepositoryLocal struct {
   mu     sync.RWMutex
   tables map[string]*model.OrderTable
}

func NewOrderTableRepositoryLocal() model.OrderTableRepository {
   return &OrderTableRepositoryLocal{tables: make(map[string]*model.OrderTable)}
}

func (r *OrderTableRepositoryLocal) CreateOrderTable(ctx context.Context, table *model.OrderTable) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.tables[table.ID.String()] = table
   return nil
}

func (r *OrderTableRepositoryLocal) UpdateOrderTable(ctx context.Context, table *model.OrderTable) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.tables[table.ID.String()] = table
   return nil
}

func (r *OrderTableRepositoryLocal) DeleteOrderTable(ctx context.Context, id string) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   delete(r.tables, id)
   return nil
}

func (r *OrderTableRepositoryLocal) GetOrderTableById(ctx context.Context, id string) (*model.OrderTable, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   if t, ok := r.tables[id]; ok {
       return t, nil
   }
   return nil, nil
}

func (r *OrderTableRepositoryLocal) GetPendingOrderTablesByTableId(ctx context.Context, id string) ([]model.OrderTable, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := []model.OrderTable{}
   for _, t := range r.tables {
       if t.TableID.String() == id && t.Status == "Pending" {
           out = append(out, *t)
       }
   }
   return out, nil
}

func (r *OrderTableRepositoryLocal) GetAllOrderTables(ctx context.Context) ([]model.OrderTable, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := make([]model.OrderTable, 0, len(r.tables))
   for _, t := range r.tables {
       out = append(out, *t)
   }
   return out, nil
}
