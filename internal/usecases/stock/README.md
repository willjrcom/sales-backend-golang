# Usecase / Stock

Responsável por manter o estoque consistente nas operações de pedido, cancelamento e ajustes manuais. Atua sobre `stock`, `stock_batch`, `stock_movement`, `order` e `order_item`, garantindo FIFO, alertas automáticos e reconciliação diária.

---

## 1. Pontos de entrada (handlers/métodos)

| Método | Handler / Usecase | Descrição |
|--------|-------------------|-----------|
| `ReserveFromItem` | Invocado via `POST /order/{orderID}/items` | Reserva quantidade quando um item é criado. |
| `DebitFromOrder` | `POST /order/{orderID}/finish` | Consome lotes FIFO ao finalizar o pedido. |
| `RestoreFromItem` | `DELETE /order/{orderID}/items/{itemID}` | Devolve reservas quando item é removido/cancelado. |
| `AddMovement` | `POST /stock/{id}/movement/add` | Entrada manual (compra, devolução). |
| `RemoveMovement` | `POST /stock/{id}/movement/remove` | Saída manual (quebra, desperdício). |
| `AdjustMovement` | `POST /stock/{id}/movement/adjust` | Ajuste inventário (define estoque alvo). |
| `ListMovements` | `GET /stock/{id}/movement` | Histórico paginado de movimentos. |
| `ListAlerts` | `GET /stock/alerts` | Alertas ativos (baixo, sem estoque, over, expiração). |
| `RunReconciliation` | job em `scheduler/daily_scheduler.go` | Recalcula estoque a partir dos movimentos. |

---

## 2. Dependências injetadas

- Repositórios: `stock`, `stock_movement`, `product`, `order`, `order_item`, `employee`.
- Serviços: `scheduler.CheckAlerts`, `report.StockReport`, `rabbitmq` (eventos `stock.alert.created`).
- Configurações: preferências de estoque da empresa (`allow_negative_stock`, `alert_threshold_days`).

---

## 3. Fluxos detalhados e exemplos

### 3.1 Reserva automática (`ReserveFromItem`)
1. Handler de pedido cria item → chama `ReserveFromItem(ctx, orderID, itemID)`.
2. `SELECT ... FOR UPDATE` no registro de estoque (produto/variação).
3. Se `allow_negative_stock=false` e `current_stock < qty`, retorna `ErrInsufficientStock`.
4. Atualiza `current_stock -= qty`, `reserved_stock += qty`, cria movimento `RESERVE`.
5. Executa `CheckAlerts`.

**Exemplo**
```
POST /order/ord-123/items
{
  "product_id": "prod-01",
  "variation_id": "var-200",
  "quantity": 2
}
```
Resposta (trecho):
```
{
  "item_id": "item-777",
  "stock_reservation": {
    "movement_id": "mov-abc",
    "type": "RESERVE",
    "current_stock_after": 18,
    "reserved_stock_after": 4
  }
}
```

### 3.2 Débito FIFO (`DebitFromOrder`)
1. Handler `finish order` chama `DebitFromOrder(ctx, orderID)`.
2. Lista itens/grupos em aberto; para cada um busca lotes ordenados por `created_at`.
3. Consome lote a lote, criando movimentos `OUT` com referência ao `batch_id`.
4. Ajusta `reserved_stock -= qty`, `current_stock` já reduzido na reserva.
5. Se lotes insuficientes e `allow_negative_stock=true`, gera `OUT_NEGATIVE`.

**Exemplo**
```
POST /order/ord-123/finish
```
Resposta:
```
{
  "order_id": "ord-123",
  "stock_debit": [
    {
      "stock_id": "stk-01",
      "movements": [
        {"type": "OUT", "batch_id": "batch-1", "quantity": 1},
        {"type": "OUT", "batch_id": "batch-2", "quantity": 1}
      ]
    }
  ]
}
```

### 3.3 Restauração (`RestoreFromItem` / cancelamento)
1. Item cancelado antes da finalização → `RestoreFromItem` reverte `RESERVE`.
2. Pedido finalizado cancelado → consulta movimentos `OUT` por `order_id` e repõe lotes.

**Exemplo**
```
DELETE /order/ord-123/items/item-777
```
Resposta:
```
{
  "restored": true,
  "movement_id": "mov-rest-01",
  "type": "RESTORE_RESERVE",
  "current_stock_after": 20,
  "reserved_stock_after": 0
}
```

### 3.4 Movimentações manuais (`Add/Remove/Adjust`)
**Entrada (`AddMovement`)**
```
POST /stock/stk-01/movement/add
{
  "quantity": 50,
  "reason": "Compra fornecedor ABC",
  "cost_price": 8.5,
  "expires_at": "2026-05-01"
}
```
Resposta:
```
{
  "movement_id": "mov-in-01",
  "type": "IN",
  "batch_id": "batch-99",
  "current_stock_after": 150
}
```

**Saída (`RemoveMovement`)**
```
POST /stock/stk-01/movement/remove
{
  "quantity": 5,
  "reason": "Produto vencido"
}
```
Resposta:
```
{
  "movement_id": "mov-out-05",
  "type": "OUT",
  "batches": [
    {"batch_id": "batch-90", "quantity": 3},
    {"batch_id": "batch-91", "quantity": 2}
  ]
}
```

**Ajuste (`AdjustMovement`)**
```
POST /stock/stk-01/movement/adjust
{
  "new_stock": 80,
  "reason": "Inventário trimestral"
}
```
Resposta:
```
{
  "movement_id": "mov-adjust-02",
  "type": "ADJUST_OUT",
  "difference": -10
}
```

### 3.5 Histórico e alertas
**Listar movimentos**
```
GET /stock/stk-01/movement?page=1&limit=20
```

**Listar alertas**
```
GET /stock/alerts?status=open
```
Resposta:
```
[
  {
    "alert_id": "al-001",
    "stock_id": "stk-01",
    "type": "low_stock",
    "current_stock": 5,
    "threshold": 10
  }
]
```

### 3.6 Reconciliação (`RunReconciliation`)
- Job diário percorre schemas, soma movimentos (`IN - OUT + RESERVE`), compara com `stock.current_stock`.
- Divergências geram alertas `adjustment_needed` e movimentos `ADJUST_*`.

---

## 4. Falhas tratadas
- `ErrStockInactive`: tentativa de mover estoque desativado → retorna 409.
- `ErrBatchLockTimeout`: retry com backoff antes de falhar 503.
- `ErrInsufficientStock`: bloqueia saídas manuais; reservas obedecem `allow_negative_stock`.

---

## 5. Observações operacionais
- Execute `RunReconciliation` após migrações ou ajustes diretos no banco.
- Alertas críticos podem disparar webhook configurado em `company.preferences.stock_alert_webhook`.
- Logs devem incluir `stock_id`, `product_id`, `order_id` e `movement_id`.
