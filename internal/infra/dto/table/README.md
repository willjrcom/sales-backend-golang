# DTO / Table

DTOs para o CRUD de mesas.

---

## 1. Onde é usado
- handler/table.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| TableRequest | name, place_id, capacity, virtual | request |
| TableResponse | id, name, capacity, status, virtual | response |

## 3. Regras de validação
- `capacity` >=1.

## 4. Exemplo de request
```json
{
  "name": "Mesa 10",
  "place_id": "plc-2",
  "capacity": 4,
  "virtual": false
}
```

## 5. Exemplo de response
```json
{
  "id": "tbl-10",
  "name": "Mesa 10",
  "capacity": 4,
  "status": "available",
  "virtual": false
}
```

## 6. Notas e compatibilidade
- Virtual tables servem para hubs de entrega; manter flag true.
