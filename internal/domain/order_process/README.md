# Domain / Order Process

Define workflow de produção (processos, filas, status).

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Process | Pipeline principal. |
| Queue | Instâncias em andamento. |
| StatusProcess | Enum de estados. |

## 2. Regras de negócio
- Cada produto pode cadastrar múltiplas etapas.
- Fila mantém posição e timestamps para analytics.

## 3. Interações e consumidores
- Usecases: order_process, order_queue, print_manager.

## 4. Exemplo de estrutura
```json
{
  "process_id": "proc-10",
  "name": "Cozinha quente",
  "steps": [
    {
      "station": "grill",
      "duration_sec": 300
    }
  ]
}
```
