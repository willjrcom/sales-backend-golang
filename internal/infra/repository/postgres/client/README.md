# Repository / Postgres / Client

CRUD de clientes finais, histórico e buscas por documento/telefone.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Search(ctx, filter)` | Suporta LIKE em nome, documento e telefone com paginação. |
| `GetWithStats(ctx, id)` | Retorna cliente + métricas (último pedido, ticket médio). |
| `Create(ctx, Client)` | Insere client + preferences usando RETURNING. |

## 2. Transações e locking
- Create/Update roda com person/contact/address dentro da mesma tx.

## 3. Exemplo de SQL
```sql
SELECT c.id, c.loyalty_score, COUNT(o.id) total_orders
FROM clients c
LEFT JOIN orders o ON o.client_id = c.id
WHERE c.id = @id
GROUP BY c.id;
```

## 4. Notas operacionais
- Use índices em documento/telefone para acelerar busca.
- Sempre normalizar telefone/documento antes de persistir.
