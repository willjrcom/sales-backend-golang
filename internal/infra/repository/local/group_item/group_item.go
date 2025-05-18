package groupitemrepositorylocal

import (
   "context"
   "sync"

   "github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

// GroupItemRepositoryLocal is an in-memory group item repository
type GroupItemRepositoryLocal struct {
   mu    sync.RWMutex
   items map[string]*model.GroupItem
}

func NewGroupItemRepositoryLocal() model.GroupItemRepository {
   return &GroupItemRepositoryLocal{items: make(map[string]*model.GroupItem)}
}

func (r *GroupItemRepositoryLocal) CreateGroupItem(ctx context.Context, groupitem *model.GroupItem) (err error) {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.items[groupitem.ID.String()] = groupitem
   return nil
}

func (r *GroupItemRepositoryLocal) UpdateGroupItem(ctx context.Context, groupitem *model.GroupItem) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.items[groupitem.ID.String()] = groupitem
   return nil
}

func (r *GroupItemRepositoryLocal) GetGroupByID(ctx context.Context, id string, withRelation bool) (*model.GroupItem, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   if gi, ok := r.items[id]; ok {
       return gi, nil
   }
   return nil, nil
}

func (r *GroupItemRepositoryLocal) DeleteGroupItem(ctx context.Context, id string, complementItemID *string) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   delete(r.items, id)
   return nil
}

func (r *GroupItemRepositoryLocal) GetGroupItemsByOrderIDAndStatus(ctx context.Context, id string, status string) ([]model.GroupItem, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := []model.GroupItem{}
   for _, gi := range r.items {
       if gi.OrderID.String() == id && gi.Status == status {
           out = append(out, *gi)
       }
   }
   return out, nil
}

func (r *GroupItemRepositoryLocal) GetGroupItemsByStatus(ctx context.Context, status string) ([]model.GroupItem, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   out := []model.GroupItem{}
   for _, gi := range r.items {
       if gi.Status == status {
           out = append(out, *gi)
       }
   }
   return out, nil
}
