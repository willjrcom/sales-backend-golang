# DTO / Order Process

DTOs para filas/processos da cozinha.

---

## 1. Onde é usado
- handler/order_process.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| ProcessStatusRequest | order_id, process_id, status | request |
| ProcessStatusResponse | order_id, process_id, status, started_at, finished_at | response |

## 3. Regras de validação
- `status` ∈ {pending,in_progress,paused,done}.

## 4. Exemplo de request
```json
{
  "order_id": "ord-500",
  "process_id": "proc-10",
  "status": "in_progress"
}
```

## 5. Exemplo de response
```json
{
  "order_id": "ord-500",
  "process_id": "proc-10",
  "status": "in_progress",
  "started_at": "2026-03-05T12:00:00Z"
}
```

## 6. Notas e compatibilidade
- Sempre registrar `employee_id` responsável.
