# DTO / S3

Estruturas para assinar upload/download de arquivos.

---

## 1. Onde é usado
- handler/s3.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| PresignRequest | folder, filename, content_type | request |
| PresignResponse | upload_url, key, expires_in | response |

## 3. Regras de validação
- `content_type` obrigatório.

## 4. Exemplo de request
```json
{
  "folder": "products",
  "filename": "combo.png",
  "content_type": "image/png"
}
```

## 5. Exemplo de response
```json
{
  "upload_url": "https://s3?...",
  "key": "companies/cmp-100/products/combo.png",
  "expires_in": 900
}
```

## 6. Notas e compatibilidade
- Use HTTPS sempre.
