# DTO / Keys

DTOs para distribuição de chaves públicas (ex.: PIX, integrações).

---

## 1. Onde é usado
- handler/public.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| PublicKeyResponse | service, key, updated_at | response |

## 3. Regras de validação
- Não expor segredos.

## 4. Exemplo de request
```json
{}
```

## 5. Exemplo de response
```json
{
  "service": "pix",
  "key": "222e...",
  "updated_at": "2026-03-05T12:00:00Z"
}
```

## 6. Notas e compatibilidade
- Cache recomendado para evitar hits diretos no banco.
