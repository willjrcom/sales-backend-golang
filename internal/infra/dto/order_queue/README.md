# DTO / Order Queue

DTOs para monitorar filas e posições.

---

## 1. Onde é usado
- handler/order_queue.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| QueueEntryResponse | order_id, process_id, position, status | response |

## 3. Regras de validação
- `position` inicia em 1 e não aceita valores negativos.

## 4. Exemplo de request
```json
{}
```

## 5. Exemplo de response
```json
{
  "order_id": "ord-500",
  "process_id": "proc-kitch",
  "position": 2,
  "status": "in_progress"
}
```

## 6. Notas e compatibilidade
- Atualizações devem ser transmitidas via SSE/WebSocket.
