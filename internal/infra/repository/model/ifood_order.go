package model

import (
"encoding/json"
"github.com/google/uuid"
"github.com/uptrace/bun"
entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type IfoodOrder struct {
entitymodel.Entity
bun.BaseModel `bun:"table:ifood_orders"`

IfoodOrderID string    `bun:"ifood_order_id,notnull,unique"`
Status       string    `bun:"status,notnull"`
RawPayload   json.RawMessage `bun:"raw_payload,type:jsonb,nullzero"`

OrderID        uuid.UUID  `bun:"order_id,type:uuid,notnull"`
OrderDeliveryID *uuid.UUID `bun:"order_delivery_id,type:uuid,nullzero"`
}
