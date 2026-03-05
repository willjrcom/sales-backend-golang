# Usecase / User

Administra usuários do sistema, papéis (RBAC) e credenciais.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/users` | handler/user.go | Cria usuário vinculado à empresa. |
| POST | `/users/{id}/reset-password` | handler/user.go | Dispara reset. |
| POST | `/users/{id}/roles` | handler/user.go | Atualiza papéis. |

## 2. Dependências
- Repositories: user, employee, company.
- Services: bcrypt, jwt, email.

## 3. Fluxos e exemplos
### Criar usuário
Passos:
- Valida email único no schema.
- Hash de senha com bcrypt.
- Associa papéis padrão e envia convite.

Exemplo de request:
```json
{
  "email": "admin@loja.com",
  "role": "admin",
  "employee_id": "emp-10"
}
```
Resposta:
```json
{
  "user_id": "usr-10",
  "status": "invited"
}
```

### Reset de senha
Passos:
- Gera token JWT de reset.
- Envia email com link.
- Audita solicitação.

Exemplo de request:
```json
{}
```
Resposta:
```json
{
  "sent": true
}
```

## 4. Falhas conhecidas
- ErrUserDuplicate
- ErrRoleNotAllowed

## 5. Notas operacionais
- Bloquear usuário deve revogar todos os tokens ativos imediatamente.
