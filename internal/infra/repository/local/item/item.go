package itemrepositorylocal

import (
   "context"
   "sync"
   "github.com/google/uuid"
   "github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

// ItemRepositoryLocal is an in-memory implementation of ItemRepository
type ItemRepositoryLocal struct {
   mu    sync.RWMutex
   items map[string]*model.Item
}

func NewItemRepositoryLocal() model.ItemRepository {
   return &ItemRepositoryLocal{items: make(map[string]*model.Item)}
}

func (r *ItemRepositoryLocal) AddItem(ctx context.Context, item *model.Item) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.items[item.ID.String()] = item
   return nil
}

func (r *ItemRepositoryLocal) AddAdditionalItem(ctx context.Context, id uuid.UUID, productID uuid.UUID, additionalItem *model.Item) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   pid := id.String()
   if parent, ok := r.items[pid]; ok {
       parent.AdditionalItems = append(parent.AdditionalItems, *additionalItem)
       return nil
   }
   return nil
}

func (r *ItemRepositoryLocal) DeleteItem(ctx context.Context, id string) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   delete(r.items, id)
   return nil
}

func (r *ItemRepositoryLocal) DeleteAdditionalItem(ctx context.Context, idAdditional uuid.UUID) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   aid := idAdditional.String()
   for _, item := range r.items {
       newList := item.AdditionalItems[:0]
       for _, ai := range item.AdditionalItems {
           if ai.ID.String() != aid {
               newList = append(newList, ai)
           }
       }
       item.AdditionalItems = newList
   }
   return nil
}

func (r *ItemRepositoryLocal) UpdateItem(ctx context.Context, item *model.Item) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.items[item.ID.String()] = item
   return nil
}

func (r *ItemRepositoryLocal) GetItemById(ctx context.Context, id string) (*model.Item, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   if it, ok := r.items[id]; ok {
       return it, nil
   }
   return nil, nil
}

func (r *ItemRepositoryLocal) GetItemByAdditionalItemID(ctx context.Context, idAdditional uuid.UUID) (*model.Item, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   aid := idAdditional.String()
   for _, item := range r.items {
       for _, ai := range item.AdditionalItems {
           if ai.ID.String() == aid {
               return &ai, nil
           }
       }
   }
   return nil, nil
}
