# Repository / Postgres / Address

Gerencia endereços multi-tenant incluindo geocode e vínculos com empresas/clientes/pedidos.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `GetByID(ctx, id)` | Busca endereço por ID respeitando schema. |
| `ListByEntity(ctx, entityID, entityType)` | Lista todos os endereços associados a uma pessoa/empresa. |
| `Upsert(ctx, Address)` | Atualiza campos e geocode em operação única. |

## 2. Transações e locking
- Atualizações acontecem na mesma transação do cadastro de empresa/cliente para evitar órfãos.

## 3. Exemplo de SQL
```sql
SELECT id, zip_code, street, city, state, geocode_lat, geocode_lng
FROM addresses
WHERE company_id = @company AND id = @id;
```

## 4. Notas operacionais
- Aplicar `SetSearchPath` antes das queries.
- Geocode é opcional: mantenha colunas nulas quando provider não retornar.
