# Service / S3

Cliente AWS S3 responsável por upload/download de arquivos (imagens de produto, PDFs fiscais).

---

## 1. Responsabilidades
- Gerar URLs pré-assinadas.
- Upload com ACL privada.
- Excluir assets órfãos.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `Upload(ctx, key string, body io.Reader, metadata map[string]string) (URL, error)` | Envia arquivo. |
| `GeneratePresignedURL(key string, ttl time.Duration) (string, error)` | URL temporária. |
| `Delete(key string) error` | Remove asset. |

## 3. Fluxo típico
- Handler recebe arquivo → chama `Upload` → salva URL retornada.
- Para downloads públicos controlados, usa `GeneratePresignedURL`.

## 4. Configuração / Env Vars
- `AWS_REGION`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `S3_BUCKET`

## 5. Exemplo de uso
```go
go
url, err := s3.Upload(ctx, fmt.Sprintf("companies/%s/products/%s.png", companyID, productID), file, nil)
```

## 6. Falhas comuns
- ErrFileTooLarge
- ErrS3Unavailable

## 7. Notas operacionais
- Sempre versionar chaves com IDs para facilitar limpeza.
