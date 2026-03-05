# Repository / Postgres / Size

Tamanhos associados a produtos/variações.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Upsert(ctx, size)` | Insere ou atualiza mantendo SKU suffix único. |
| `ListByProduct(ctx)` | Retorna tamanhos ativos e preços. |
| `Deactivate(ctx)` | Marca tamanho como inativo mantendo histórico. |

## 2. Transações e locking
- Operações caminham junto do productRepository; use tx para manter sincronismo.

## 3. Exemplo de SQL
```sql
SELECT id, name, price
FROM product_sizes
WHERE product_id=@product AND status='active';
```

## 4. Notas operacionais
- Preço salvo como numeric(12,2).
- Nunca delete físico; apenas altere `status`.
