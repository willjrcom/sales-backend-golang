# Repository / Postgres / Company Category

Relaciona empresas às categorias operacionais e patrocinadores exclusivos.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Assign(ctx, companyID, categories[])` | Remove vínculos antigos e insere novos em batch. |
| `ListWithSponsors(ctx, companyID)` | JOIN em sponsors/patrocínios para montar resposta. |

## 2. Transações e locking
- `Assign` roda em tx para garantir consistência entre delete/insert.

## 3. Exemplo de SQL
```sql
INSERT INTO company_categories (company_id, category_id)
SELECT @company, UNNEST(@categories::text[]);
```

## 4. Notas operacionais
- Chave UNIQUE(company_id, category_id).
- Atualize caches de catálogo após mudanças.
