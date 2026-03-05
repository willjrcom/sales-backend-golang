# Repository / Postgres / Place

Locais físicos/dark kitchens e seus canais.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Create(ctx, place)` | Insere local com endereço e canais. |
| `List(ctx)` | JOIN addresses/tables para dashboards. |
| `UpdateChannels(ctx)` | Atualiza jsonb de canais e horários. |

## 2. Transações e locking
- Atualização de place e tables relacionadas roda em tx quando mover mesas.

## 3. Exemplo de SQL
```sql
SELECT p.id, p.name, p.channels, a.city
FROM places p
JOIN addresses a ON a.id=p.address_id
WHERE p.company_id=@company;
```

## 4. Notas operacionais
- Horários armazenados em JSON; respeitar timezone da empresa.
