# Repository / Postgres / Contact

Tabela de telefones/emails compartilhados por pessoas.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `GetPrimary(ctx, personID, type)` | Retorna contato principal. |
| `Upsert(ctx, Contact)` | Cria/atualiza contato garantindo unicidade por tipo. |

## 2. Transações e locking
- Executado junto com person/client/employee; utilize a mesma tx.

## 3. Exemplo de SQL
```sql
SELECT id, value FROM contacts
WHERE person_id=@person AND type=@type
ORDER BY is_primary DESC, created_at ASC
LIMIT 1;
```

## 4. Notas operacionais
- Normalizar valor (lowercase / E.164) antes da query.
- Garanta apenas um `is_primary` por tipo.
