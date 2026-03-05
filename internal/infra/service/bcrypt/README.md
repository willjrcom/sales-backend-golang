# Service / Bcrypt

Fornece funções utilitárias para gerar e validar hashes de senha usando bcrypt com custo configurável.

---

## 1. Responsabilidades
- Gerar hash seguro para senhas de usuários internos.
- Comparar senha informada com hash armazenado.
- Expor fator de custo vindo de configuração para balancear segurança x performance.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `Hash(password string) (string, error)` | Gera hash usando custo definido em `SECURITY_BCRYPT_COST` (padrão 12). |
| `Compare(hash, password string) error` | Retorna erro quando a senha não corresponde ou hash inválido. |

## 3. Fluxo típico
- Receive plaintext password from user onboarding/reset.
- Validate minimum length before hashing.
- Call bcrypt.GenerateFromPassword and persist result via user repository.
- During login, load hash and call Compare to authorize.

## 4. Configuração / Env Vars
- `SECURITY_BCRYPT_COST (opcional, default 12).`

## 5. Exemplo de uso
```go
go
hash, err := bcryptservice.Hash("senha-super-segura")
if err != nil { /* tratar */ }
if err := bcryptservice.Compare(hash, input); err != nil {
    return errors.New("invalid credentials")
}
```

## 6. Falhas comuns
- bcrypt.ErrPasswordTooLong se senha exceder limite.
- Erro genérico quando custo é inválido (<4 ou >31).

## 7. Notas operacionais
- Sempre sanitize entradas antes de logar para evitar vazamento de senha.
- Alterar o custo exige rehash progressivo dos usuários antigos.
