package model

import (
"time"

"github.com/uptrace/bun"
entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

// IfoodConnection stores iFood credentials/tokens per tenant schema.
// NOTE: this table is expected to be created in the PUBLIC schema.
type IfoodConnection struct {
entitymodel.Entity
bun.BaseModel `bun:"table:ifood_connections"`

Schema         string     `bun:"schema,notnull,unique"`
MerchantID     string     `bun:"merchant_id,notnull"`
ClientID       string     `bun:"client_id,notnull"`
ClientSecret   string     `bun:"client_secret,notnull"`
AccessToken    *string    `bun:"access_token,nullzero"`
RefreshToken   *string    `bun:"refresh_token,nullzero"`
TokenExpiresAt *time.Time `bun:"token_expires_at,nullzero"`
Sandbox        bool       `bun:"sandbox,notnull,default:false"`
}
