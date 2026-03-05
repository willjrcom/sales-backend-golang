# DTO / Place

DTOs de locais físicos/dark kitchens.

---

## 1. Onde é usado
- handler/place.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| PlaceRequest | name, address_id, channels[], opening_hours | request |
| PlaceResponse | id, name, status, channels, opening_hours | response |

## 3. Regras de validação
- `channels` subset de {delivery,pickup,dine_in}.

## 4. Exemplo de request
```json
{
  "name": "Loja Norte",
  "address_id": "addr-9",
  "channels": [
    "delivery",
    "pickup"
  ],
  "opening_hours": [
    {
      "day": "monday",
      "from": "09:00",
      "to": "22:00"
    }
  ]
}
```

## 5. Exemplo de response
```json
{
  "id": "plc-2",
  "name": "Loja Norte",
  "status": "active",
  "channels": [
    "delivery",
    "pickup"
  ]
}
```

## 6. Notas e compatibilidade
- Horários devem considerar timezone da empresa.
