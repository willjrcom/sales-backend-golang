# DTO / Order Pickup

DTOs para pedidos de retirada.

---

## 1. Onde é usado
- handler/order_pickup.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| OrderPickupRequest | order_id, scheduled_for, contact_phone | request |
| OrderPickupResponse | order_id, status, pickup_code | response |

## 3. Regras de validação
- `scheduled_for` timezone ISO8601.

## 4. Exemplo de request
```json
{
  "order_id": "ord-500",
  "scheduled_for": "2026-03-05T19:30:00-03:00",
  "contact_phone": "+5511988887777"
}
```

## 5. Exemplo de response
```json
{
  "order_id": "ord-500",
  "status": "ready",
  "pickup_code": "AB12"
}
```

## 6. Notas e compatibilidade
- `pickup_code` deve ser único por dia.
