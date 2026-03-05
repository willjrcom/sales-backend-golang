# Repository / Postgres / Process Rule

Regras que ligam categorias/produtos às etapas de produção.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Create(ctx, rule)` | Insere regra e steps associados. |
| `List(ctx)` | JOIN categorias e estações para montar pipeline. |
| `Deactivate(ctx)` | Marca regra como inativa preservando histórico. |

## 2. Transações e locking
- Criar/atualizar steps de produção deleta e reinsere dentro de tx para garantir ordem.

## 3. Exemplo de SQL
```sql
SELECT r.id, r.name, json_agg(s.* ORDER BY s.position) steps
FROM process_rules r
JOIN process_rule_steps s ON s.rule_id=r.id
WHERE r.company_id=@company
GROUP BY r.id;
```

## 4. Notas operacionais
- Evitar cascade delete; mantenha histórico para auditoria.
