# Domain / Table

Mesas físicas ou virtuais usadas em pedidos de salão.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Table | Identificador, capacidade, local. |
| TableStatus | Enum available, occupied, blocked. |

## 2. Regras de negócio
- Nome único por `place`.
- Mesas ocupadas não podem ser reatribuídas sem transferir pedido.

## 3. Interações e consumidores
- Usecases: table, order_table, place.

## 4. Exemplo de estrutura
```json
{
  "id": "tbl-10",
  "place_id": "plc-2",
  "name": "Mesa 10",
  "capacity": 4,
  "status": "available"
}
```
