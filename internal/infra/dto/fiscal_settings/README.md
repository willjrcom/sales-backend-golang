# DTO / Fiscal Settings

DTOs para configurar credenciais fiscais.

---

## 1. Onde é usado
- handler/fiscal_settings.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| FiscalSettingsRequest | company_id, crt, cnae, certificate_base64, password | request |
| FiscalSettingsResponse | company_id, status, last_tested_at | response |

## 3. Regras de validação
- Certificado deve estar em base64 A1.
- `password` guardada de forma segura; não retornar.

## 4. Exemplo de request
```json
{
  "company_id": "cmp-100",
  "crt": 1,
  "cnae": "5611201",
  "certificate_base64": "...",
  "password": "secret"
}
```

## 5. Exemplo de response
```json
{
  "company_id": "cmp-100",
  "status": "ready",
  "last_tested_at": "2026-03-05T12:00:00Z"
}
```

## 6. Notas e compatibilidade
- Nunca logar payload completo.
