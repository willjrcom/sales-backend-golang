# Repository / Postgres / Product Category

Categorias hierárquicas e vínculos com process rules/patrocínios.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `ListTree(ctx)` | Retorna árvore via CTE recursiva. |
| `Create(ctx)` | Insere categoria com parent opcional. |
| `UpdateProcessRule(ctx)` | Atualiza relação categoria → process rule. |

## 2. Transações e locking
- Alterar hierarquia roda com lock/validação para evitar loops.

## 3. Exemplo de SQL
```sql
WITH RECURSIVE tree AS (
  SELECT id, name, parent_id, 1 depth FROM product_categories WHERE parent_id IS NULL
  UNION ALL
  SELECT pc.id, pc.name, pc.parent_id, t.depth+1
  FROM product_categories pc JOIN tree t ON pc.parent_id=t.id
)
SELECT * FROM tree;
```

## 4. Notas operacionais
- Garantir slug único para integrar com frontend.
- Atualizar caches após mudanças.
