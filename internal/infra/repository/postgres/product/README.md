# Repository / Postgres / Product

Produtos, variações e tamanhos com controle de estoque.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `CreateFull(ctx, Product)` | Insere produto + variations + sizes em tx. |
| `List(ctx, filters)` | JOIN categorias, estoque e patrocinadores. |
| `UpdateAvailability(ctx)` | Ativa/desativa produto por canal. |

## 2. Transações e locking
- Quando `track_stock=true`, criar registros em stock na mesma tx.
- Atualizações críticas usam versionamento.

## 3. Exemplo de SQL
```sql
SELECT p.id, p.name, p.track_stock, json_agg(v.*) variations
FROM products p
LEFT JOIN product_variations v ON v.product_id=p.id
WHERE p.company_id=@company
GROUP BY p.id;
```

## 4. Notas operacionais
- Indexar (company_id, lower(name)) para busca rápida.
- Sincronizar mudanças com caches do PDV.
