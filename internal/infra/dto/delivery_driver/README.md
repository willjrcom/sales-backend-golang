# DTO / Delivery Driver

Payloads para cadastro e atualização de entregadores.

---

## 1. Onde é usado
- handler/delivery_driver.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| DriverRequest | employee_id, vehicle, zones[], documents | request |
| DriverResponse | id, employee_id, status, vehicle, zones | response |

## 3. Regras de validação
- `vehicle` ∈ {bike,moto,car}.
- `zones` referenciam place/region válidos.

## 4. Exemplo de request
```json
{
  "employee_id": "emp-10",
  "vehicle": "bike",
  "zones": [
    "north"
  ]
}
```

## 5. Exemplo de response
```json
{
  "id": "drv-1",
  "employee_id": "emp-10",
  "status": "available",
  "vehicle": "bike",
  "zones": [
    "north"
  ]
}
```

## 6. Notas e compatibilidade
- Inclua `tracking_enabled` quando disponível.
