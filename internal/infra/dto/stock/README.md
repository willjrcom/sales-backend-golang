# DTO / Stock

DTOs para criação de estoque, movimentos e alertas.

---

## 1. Onde é usado
- handler/stock.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| StockMovementRequest | quantity, reason, cost_price, expires_at | request |
| StockMovementResponse | movement_id, type, batch_id, quantity, current_stock_after | response |
| StockAlertResponse | id, stock_id, type, current_stock, threshold | response |

## 3. Regras de validação
- `quantity` > 0 para add/remove.
- `expires_at` ISO8601 opcional.

## 4. Exemplo de request
```json
{
  "quantity": 10,
  "reason": "Compra fornecedor",
  "cost_price": 8.5,
  "expires_at": "2026-05-01"
}
```

## 5. Exemplo de response
```json
{
  "movement_id": "mov-1",
  "type": "IN",
  "batch_id": "batch-1",
  "quantity": 10,
  "current_stock_after": 50
}
```

## 6. Notas e compatibilidade
- Todos os valores retornam em unidades exatas, sem arredondar.
