# Usecase / Delivery Driver

Gerencia entregadores internos/terceirizados, suas taxas e disponibilidade para pedidos delivery.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/delivery-drivers` | handler/delivery_driver.go | Lista entregadores e status em tempo real. |
| POST | `/delivery-drivers` | handler/delivery_driver.go | Cadastra entregador vinculado a um employee. |
| PATCH | `/delivery-drivers/{id}/status` | handler/delivery_driver.go | Atualiza disponibilidade, veículo e zona. |

## 2. Dependências
- Repositories: delivery_driver, employee, shift.
- Services: rabbitmq (eventos de tracking).

## 3. Fluxos e exemplos
### Cadastro
Passos:
- Valida se employee existe e possui documentos obrigatórios.
- Registra taxas fixas/percentuais.
- Habilita tracking opcional.

Exemplo de request:
```json
{
  "employee_id": "emp-10",
  "vehicle": "bike",
  "zones": [
    "zona-1"
  ]
}
```
Resposta:
```json
{
  "id": "drv-1",
  "status": "available"
}
```

### Atualizar status
Passos:
- Recebe status `available`, `busy`, `off`.
- Grava timestamp e motivo.
- Notifica fila de pedidos para redistribuir.

Exemplo de request:
```json
{
  "status": "busy",
  "order_id": "ord-50"
}
```
Resposta:
```json
{
  "driver_id": "drv-1",
  "status": "busy",
  "updated_at": "2026-03-05T12:00:00Z"
}
```

## 4. Falhas conhecidas
- ErrDriverWithoutEmployee
- ErrInvalidZone

## 5. Notas operacionais
- Status é centralizado; qualquer alteração dispara broadcast para dashboards.
