# Repository / Postgres / User

Usuários internos, papéis (RBAC) e credenciais.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Create(ctx, user)` | Insere usuário e hash de senha. |
| `List(ctx, filters)` | Filtra por role/status com paginação. |
| `UpdatePassword(ctx)` | Atualiza hash e invalida tokens ativos. |
| `AssignRoles(ctx)` | Manipula tabela pivot user_roles. |

## 2. Transações e locking
- Alterar roles e employee vinculado acontece na mesma tx.
- Password update deve limpar refresh tokens.

## 3. Exemplo de SQL
```sql
SELECT u.id, u.email, array_agg(r.role) roles
FROM users u
LEFT JOIN user_roles r ON r.user_id=u.id
WHERE u.company_id=@company
GROUP BY u.id;
```

## 4. Notas operacionais
- Campos sensíveis (hash, tokens) nunca devem ser retornados.
- Logar operações com user_id/company_id para auditoria.
