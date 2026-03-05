# Service / JWT

Emissão e validação de tokens (id_token e access_token) considerando multi-tenant e refresh flow.

---

## 1. Responsabilidades
- Gerar id_token (login) e access_token (seleção empresa).
- Validar assinatura e expiração.
- Extrair claims customizadas (company_id, schema, roles).

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `GenerateIDToken(user User) (string, error)` | Expira em 30 min. |
| `GenerateAccessToken(user User, schema string) (string, error)` | Expira em 2h. |
| `Parse(token string) (Claims, error)` | Valida assinatura + expiração. |

## 3. Fluxo típico
- User login → `GenerateIDToken`.
- Seleção de empresa → `GenerateAccessToken` com schema.
- Middleware chama `Parse` para inserir claims em contexto.

## 4. Configuração / Env Vars
- `JWT_SECRET`
- `JWT_ISSUER`
- `JWT_AUDIENCE`

## 5. Exemplo de uso
```go
go
token, err := jwt.GenerateAccessToken(user, schema)
claims, err := jwt.Parse(token)
```

## 6. Falhas comuns
- ErrTokenExpired
- ErrInvalidSignature

## 7. Notas operacionais
- Alterar segredo exige derrubar tokens antigos (revogação).
