# Service / Header

Helpers para construir e validar cabeçalhos multi-tenant (Authorization, X-Company-Schema, etc.).

---

## 1. Responsabilidades
- Extrair company schema e user ID dos headers.
- Montar headers quando chamamos APIs externas internas.
- Aplicar validações de presença obrigatória.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `Extract(ctx context.Context, req *http.Request) (Headers, error)` | Lê Authorization + X-Company-Schema. |
| `WithSchema(ctx, schema string) context.Context` | Anexa schema ao contexto para repositórios. |
| `ForwardHeaders(headers Headers) http.Header` | Cria header pronto para chamada a outro serviço. |

## 3. Fluxo típico
- Middleware chama `Extract`, valida tokens via JWT e injeta schema no contexto.
- Handlers usam `header.FromContext` para recuperar dados sem ler HTTP novamente.

## 4. Configuração / Env Vars
- (sem configuração específica)

## 5. Exemplo de uso
```go
go
hdr, err := header.Extract(ctx, r)
if err != nil { return unauthorized }
ctx = header.WithSchema(ctx, hdr.CompanySchema)
```

## 6. Falhas comuns
- ErrMissingSchema
- ErrInvalidAuthorization

## 7. Notas operacionais
- Sempre sanitize headers antes de logar.
