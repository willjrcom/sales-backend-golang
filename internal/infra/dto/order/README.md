# DTO / Order

DTOs principais de pedidos (criação, resumo, status).

---

## 1. Onde é usado
- handler/order.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| OrderCreateRequest | client_id, type, items[], place_id, notes | request |
| OrderResponse | id, status, queue_number, totals, items | response |
| OrderStatusRequest | next_status, reason | request |

## 3. Regras de validação
- `type` ∈ {delivery,pickup,dine_in}.
- `items` não vazio.
- Status transitions validadas no usecase.

## 4. Exemplo de request
```json
{
  "client_id": "cli-200",
  "type": "delivery",
  "items": [
    {
      "product_id": "prod-1",
      "quantity": 1
    }
  ]
}
```

## 5. Exemplo de response
```json
{
  "id": "ord-500",
  "status": "pending",
  "queue_number": 32,
  "totals": {
    "gross": 35,
    "discount": 0
  },
  "items": [
    {
      "id": "item-1"
    }
  ]
}
```

## 6. Notas e compatibilidade
- Sempre incluir `source_channel` quando existir.
