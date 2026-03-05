# Usecase / Order Queue

Controla filas em tempo real (cozinha, expedição) e métricas de tempo de preparo.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/order/queue` | handler/order_queue.go | Lista filas e posição atual. |
| POST | `/order/{id}/queue/advance` | handler/order_queue.go | Avança item na fila. |
| POST | `/order/{id}/queue/reassign` | handler/order_queue.go | Reatribui responsável/processo. |

## 2. Dependências
- Repositories: order_queue, order_process, employee.
- Services: rabbitmq/websocket broadcaster.

## 3. Fluxos e exemplos
### Avançar fila
Passos:
- Localiza entrada (order_id, process).
- Incrementa etapa e registra carimbo de data/hora.
- Emite evento para painel.

Exemplo de request:
```json
{
  "process_id": "proc-kitchen",
  "employee_id": "emp-10"
}
```
Resposta:
```json
{
  "order_id": "ord-500",
  "process_id": "proc-kitchen",
  "status": "ready"
}
```

## 4. Falhas conhecidas
- ErrQueueEntryNotFound
- ErrQueueOutOfOrder

## 5. Notas operacionais
- Fila deve ser bloqueada quando processo atinge estado terminal para evitar regressão.
